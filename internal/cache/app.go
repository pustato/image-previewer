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
	"github.com/pustato/image-previewer/internal/cache/filesystem"
	"github.com/pustato/image-previewer/internal/cache/lru"
)

var _ app.App = (*AppCacheDecorator)(nil)

type AppCacheDecorator struct {
	app   app.App
	cache lru.Cache
	fs    filesystem.Filesystem
}

func NewCacheAppDecorator(app app.App, limit uint64, cachePath string) (*AppCacheDecorator, error) {
	fs, err := filesystem.NewDiskFilesystem(cachePath)
	if err != nil {
		return nil, fmt.Errorf("new cached app: %w", err)
	}

	return &AppCacheDecorator{
		app: app,
		cache: lru.NewCache(limit, func(item *lru.Item) {
			_ = fs.RemoveFile(item.FileName)
		}),
		fs: fs,
	}, nil
}

func (a *AppCacheDecorator) GetAndResize(
	ctx context.Context,
	url string,
	w, h int,
	headers http.Header,
) ([]byte, error) {
	key := a.generateKey(url, w, h)

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

	item = &lru.Item{
		FileName: key + ".jpg",
		Size:     uint64(len(content)),
	}

	if err := a.fs.WriteFile(item.FileName, content); err != nil {
		return nil, fmt.Errorf("cached app save content: %w", err)
	}
	a.cache.Set(key, item)

	return content, nil
}

func (a *AppCacheDecorator) generateKey(url string, w, h int) string {
	hash := sha256.New()

	io.WriteString(hash, url)
	io.WriteString(hash, strconv.Itoa(w))
	io.WriteString(hash, strconv.Itoa(h))

	return hex.EncodeToString(hash.Sum(nil))
}
