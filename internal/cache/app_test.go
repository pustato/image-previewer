package cache

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/pustato/image-previewer/internal/app"
	mockapp "github.com/pustato/image-previewer/internal/app/mocks"
	"github.com/pustato/image-previewer/internal/cache/filesystem"
	mockfilesystem "github.com/pustato/image-previewer/internal/cache/filesystem/mocks"
	"github.com/pustato/image-previewer/internal/cache/lru"
	mocklru "github.com/pustato/image-previewer/internal/cache/lru/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	ctx          = context.Background()
	anyCacheKey  = mock.MatchedBy(func(_ string) bool { return true })
	anyFileName  = mock.MatchedBy(func(_ string) bool { return true })
	anyCacheItem = mock.MatchedBy(func(_ *lru.Item) bool { return true })
	headers      = http.Header{}
	url          = "http://google.com/"
)

func createApp(app app.App, cache lru.Cache, fs filesystem.Filesystem) *AppCacheDecorator {
	return &AppCacheDecorator{
		app:   app,
		cache: cache,
		fs:    fs,
	}
}

func TestAppCacheDecorator_GetAndResize_Success(t *testing.T) {
	t.Run("hit cache", func(t *testing.T) {
		fileName := "some_file_name"
		item := &lru.Item{
			FileName: fileName,
		}
		result := []byte("success result")

		cache := &mocklru.Cache{}
		cache.
			On("Get", anyCacheKey).
			Once().
			Return(item, true)

		fs := &mockfilesystem.Filesystem{}
		fs.
			On("ReadFile", fileName).
			Once().
			Return(result, nil)

		unit := createApp(&mockapp.App{}, cache, fs)

		actual, err := unit.GetAndResize(ctx, url, 100, 100, headers)
		require.NoError(t, err)
		require.EqualValues(t, result, actual)
	})

	t.Run("miss cache", func(t *testing.T) {
		w, h := 100, 100
		result := []byte("success result")
		var item *lru.Item
		var fileName string

		cache := &mocklru.Cache{}
		cache.
			On("Get", anyCacheKey).
			Once().
			Return(nil, false)
		cache.
			On("Set", anyCacheKey, anyCacheItem).
			Once().
			Run(func(args mock.Arguments) {
				item = args[1].(*lru.Item)
			}).
			Return(false)

		appp := &mockapp.App{}
		appp.
			On("GetAndResize", ctx, url, w, h, headers).
			Once().
			Return(result, nil)

		fs := &mockfilesystem.Filesystem{}
		fs.
			On("WriteFile", anyFileName, result).
			Once().
			Run(func(args mock.Arguments) {
				fileName = args.String(0)
			}).
			Return(nil)

		unit := createApp(appp, cache, fs)

		actual, err := unit.GetAndResize(ctx, url, w, h, headers)
		require.NoError(t, err)
		require.EqualValues(t, result, actual)
		require.Equal(t, fileName, item.FileName)
		require.Equal(t, uint64(len(result)), item.Size)
	})
}

func TestAppCacheDecorator_GetAndResize_FS_Errors(t *testing.T) {
	t.Run("read file", func(t *testing.T) {
		fileName := "some_file_name"
		item := &lru.Item{
			FileName: fileName,
		}
		testError := errors.New("test error")

		cache := &mocklru.Cache{}
		cache.
			On("Get", anyCacheKey).
			Once().
			Return(item, true)

		fs := &mockfilesystem.Filesystem{}
		fs.
			On("ReadFile", fileName).
			Once().
			Return(nil, testError)

		unit := createApp(&mockapp.App{}, cache, fs)

		result, err := unit.GetAndResize(ctx, url, 100, 100, headers)
		require.Nil(t, result)
		require.Error(t, err)
		require.ErrorIs(t, err, testError)
	})

	t.Run("write file", func(t *testing.T) {
		w, h := 100, 100
		result := []byte("error result")
		testError := errors.New("test error")

		cache := &mocklru.Cache{}
		cache.
			On("Get", anyCacheKey).
			Once().
			Return(nil, false)

		appp := &mockapp.App{}
		appp.
			On("GetAndResize", ctx, url, w, h, headers).
			Once().
			Return(result, nil)

		fs := &mockfilesystem.Filesystem{}
		fs.
			On("WriteFile", anyFileName, result).
			Once().
			Return(testError)

		unit := createApp(appp, cache, fs)

		result, err := unit.GetAndResize(ctx, url, 100, 100, headers)
		require.Nil(t, result)
		require.Error(t, err)
		require.ErrorIs(t, err, testError)
	})
}

func TestAppCacheDecorator_GetAndResize_App_Error(t *testing.T) {
	w, h := 100, 100
	testError := errors.New("test error")

	cache := &mocklru.Cache{}
	cache.
		On("Get", anyCacheKey).
		Once().
		Return(nil, false)

	appp := &mockapp.App{}
	appp.
		On("GetAndResize", ctx, url, w, h, headers).
		Once().
		Return(nil, testError)

	unit := createApp(appp, cache, &mockfilesystem.Filesystem{})
	result, err := unit.GetAndResize(ctx, url, 100, 100, headers)
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorIs(t, err, testError)
}
