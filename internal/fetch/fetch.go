package fetch

import (
	"github.com/SlyMarbo/rss"
	"github.com/indes/flowerss-bot/pkg/client"
	"io"
	"net/http"
	"strings"
	"unicode"
)

// FetchFunc rss fetch func
func FetchFunc(client *client.HttpClient) rss.FetchFunc {
	return func(url string) (resp *http.Response, err error) {
		resp, err = client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var data []byte
		if data, err = io.ReadAll(resp.Body); err != nil {
			return nil, err
		}

		resp.Body = io.NopCloser(strings.NewReader(strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		}, string(data))))
		return
	}
}
