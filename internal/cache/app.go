package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pustato/image-previewer/internal/app"
)

var _ app.App = (*CachedApp)(nil)

type Item struct {
	key      string
	FileName string
	Size     uint64
}

type RemoveItemCallback func(item *Item)

type Cache interface {
	Get(key string) (*Item, bool)
	Set(key string, item *Item) bool
	Remove(key string)
}

type Filesystem interface {
	WriteFile(name string, content []byte) error
	ReadFile(name string) ([]byte, error)
	RemoveFile(name string) error
}

type CachedApp struct {
	app   app.App
	cache Cache
	fs    Filesystem
}

func New(app app.App, limit uint64, cachePath string) (app.App, error) {
	fs, err := NewFilesystem(cachePath)
	if err != nil {
		return nil, fmt.Errorf("new cached app: %w", err)
	}

	return &CachedApp{
		app: app,
		cache: NewCache(limit, func(item *Item) {
			_ = fs.RemoveFile(item.FileName) // log ? todo
		}),
		fs: fs,
	}, nil
}

func (a *CachedApp) GetAndResize(ctx context.Context, url string, w, h int, headers http.Header) ([]byte, error) {
	key := a.cacheKey(url, w, h)

	item, found := a.cache.Get(key)
	if found {
		content, err := a.fs.ReadFile(item.FileName)
		if err != nil {
			return nil, fmt.Errorf("cached app hit: %w", err)
		}

		return content, nil
	}

	content, err := a.app.GetAndResize(ctx, url, w, h, headers)
	if err != nil {
		return nil, fmt.Errorf("cached app proxy call: %w", err)
	}

	item = &Item{
		FileName: key + ".jpg",
		Size:     uint64(len(content)),
	}

	if err := a.fs.WriteFile(item.FileName, content); err != nil {
		return nil, fmt.Errorf("cached app save content: %w", err)
	}
	a.cache.Set(key, item)

	return content, nil
}

func (a *CachedApp) cacheKey(url string, w, h int) string {
	hash := sha256.New()

	io.WriteString(hash, url)
	io.WriteString(hash, strconv.Itoa(w))
	io.WriteString(hash, strconv.Itoa(h))

	return hex.EncodeToString(hash.Sum(nil))
}
