//go:build intgrtest
// +build intgrtest

package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

func dirSize(t *testing.T, dir string) (size uint64) {
	t.Helper()

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Error(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			t.Error(err.Error())
		}

		if info.Mode().IsRegular() {
			size += uint64(info.Size())
		}
	}

	return
}

func buildUrl(c *Config, w, h int, path string) string {
	b := strings.Builder{}
	b.WriteString("http://")
	b.WriteString(c.serviceAddr)
	b.WriteString("/")
	b.WriteString(strconv.Itoa(w))
	b.WriteString("/")
	b.WriteString(strconv.Itoa(h))
	b.WriteString("/")
	b.WriteString(c.staticAddr)
	b.WriteString(path)

	return b.String()
}

func TestSimple(t *testing.T) {
	t.Parallel()
	config := NewConfig()
	client := http.Client{}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		testData := []struct {
			url          string
			expectedFile string
		}{
			{buildUrl(config, 100, 50, "/gopher.jpg"), "testdata/gopher_100_50.jpg"},
			{buildUrl(config, 100, 1000, "/gopher.jpg"), "testdata/gopher_100_1000.jpg"},
			{buildUrl(config, 1024, 504, "/gopher.jpg"), "testdata/gopher_1024_504.jpg"},
			{buildUrl(config, 1000, 50, "/gopher.jpg"), "testdata/gopher_1000_50.jpg"},
			{buildUrl(config, 2000, 1000, "/gopher.jpg"), "testdata/gopher_2000_1000.jpg"},

			{buildUrl(config, 100, 1000, url.PathEscape("/file?name=gopher.jpg")), "testdata/gopher_100_1000.jpg"},
			{buildUrl(config, 2000, 1000, url.PathEscape("/file?name=gopher.jpg")), "testdata/gopher_2000_1000.jpg"},
		}

		for i, td := range testData {
			td := td
			t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
				t.Parallel()

				rsp, err := client.Get(td.url)
				require.NoError(t, err)
				defer rsp.Body.Close()
				require.Equal(t, http.StatusOK, rsp.StatusCode)

				actual, err := io.ReadAll(rsp.Body)
				require.NoError(t, err)

				expected, err := os.ReadFile(td.expectedFile)
				require.NoError(t, err)

				require.EqualValues(t, expected, actual)
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		testData := []struct {
			url        string
			httpStatus int
			err        string
		}{
			{"http://" + config.serviceAddr + "/", http.StatusNotFound, "malformed request path"},
			{"http://" + config.serviceAddr + "/1", http.StatusNotFound, "malformed request path"},
			{"http://" + config.serviceAddr + "/1/1", http.StatusNotFound, "malformed request path"},
			{"http://" + config.serviceAddr + "/1/1/x", http.StatusBadGateway, "bad request"},
			{"http://" + config.serviceAddr + "/x/1/x", http.StatusNotFound, "width is not a number"},
			{"http://" + config.serviceAddr + "/1/x/x", http.StatusNotFound, "height is not a number"},
			{"http://" + config.serviceAddr + "/x/x/x", http.StatusNotFound, "width is not a number"},
			{buildUrl(config, 2000, 1000, "/file?name=gopher.jpg"), http.StatusBadGateway, "bad request"},
		}

		for i, td := range testData {
			td := td
			t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
				t.Parallel()

				rsp, err := client.Get(td.url)
				require.NoError(t, err)
				defer rsp.Body.Close()
				require.Equal(t, td.httpStatus, rsp.StatusCode)

				body, err := io.ReadAll(rsp.Body)
				require.NoError(t, err)

				require.Contains(t, string(body), td.err)
			})
		}
	})
}

func TestProxyHeaders(t *testing.T) {
	config := NewConfig()
	client := http.Client{}
	protectedUrl := buildUrl(config, 2000, 1000, "/protected/gopher.jpg")

	rsp, err := client.Get(protectedUrl)
	defer rsp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusBadGateway, rsp.StatusCode)

	rq, err := http.NewRequest(http.MethodGet, protectedUrl, nil)
	require.NoError(t, err)

	rq.Header.Add("X-Access", "secret")

	rsp, err = client.Do(rq)
	require.NoError(t, err)
	defer rsp.Body.Close()
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	actual, err := io.ReadAll(rsp.Body)
	require.NoError(t, err)

	expected, err := os.ReadFile("testdata/gopher_2000_1000.jpg")
	require.NoError(t, err)

	require.EqualValues(t, expected, actual)
}

func TestCacheRotation(t *testing.T) {
	config := NewConfig()
	client := http.Client{}

	for i := 10; i < 300; i += 1 {
		u := buildUrl(config, 100, i, "/gopher.jpg")

		rsp, err := client.Get(u)
		require.NoError(t, err)
		rsp.Body.Close()
		require.Equal(t, http.StatusOK, rsp.StatusCode)

		cacheDirSize := dirSize(t, config.cacheDir)
		require.LessOrEqual(t, cacheDirSize, config.cacheSize)
	}
}
