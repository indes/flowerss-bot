package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestServer(t *testing.T) *httptest.Server {
	handle := func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Method: %v", r.Method)
		t.Logf("Path: %v", r.URL.Path)

		if r.Method == http.MethodGet {
			switch r.URL.Path {
			case "/":
				_, _ = w.Write([]byte("/"))

			case "/useragent":
				for i := range r.Header["User-Agent"] {
					_, _ = w.Write([]byte(r.Header["User-Agent"][i]))
				}
			case "/timeout":
				time.Sleep(time.Second)
			}
		}
	}
	return httptest.NewServer(http.HandlerFunc(handle))
}

func TestHttpClient_Get(t *testing.T) {
	ts := createTestServer(t)
	defer ts.Close()

	t.Run("get", func(t *testing.T) {
		client := NewHttpClient()

		resp, err := client.Get(ts.URL)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Nil(t, err)
	})

	t.Run("custom client user-agent", func(t *testing.T) {
		userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
		client := NewHttpClient(WithUserAgent(userAgent))
		url := fmt.Sprintf("%s/useragent", ts.URL)
		resp, err := client.Get(url)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		t.Logf("Got: %v", string(body))
		assert.Equal(t, userAgent, string(body))
	})

	t.Run("custom get user-agent", func(t *testing.T) {
		userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
		client := NewHttpClient()
		url := fmt.Sprintf("%s/useragent", ts.URL)
		resp, err := client.Get(url, WithUserAgent(userAgent))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		t.Logf("Got: %v", string(body))
		assert.Equal(t, userAgent, string(body))
	})

	t.Run("timeout", func(t *testing.T) {
		client := NewHttpClient(WithTimeout(time.Millisecond))
		url := fmt.Sprintf("%s/timeout", ts.URL)
		_, err := client.Get(url)
		assert.Error(t, err)

		client = NewHttpClient(WithTimeout(time.Minute))
		response, err := client.Get(url)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}
