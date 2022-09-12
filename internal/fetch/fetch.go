package fetch

import (
	"github.com/SlyMarbo/rss"
	"github.com/indes/flowerss-bot/pkg/client"
	"net/http"
)

// FetchFunc rss fetch func
func FetchFunc(client *client.HttpClient) rss.FetchFunc {
	return func(url string) (resp *http.Response, err error) {
		resp, err = client.Get(url)
		if err != nil {
			return nil, err
		}
		return
	}
}
