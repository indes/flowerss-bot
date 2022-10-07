package feed

import (
	"context"
	"errors"
	"net/http"

	"github.com/indes/flowerss-bot/pkg/client"

	"github.com/mmcdole/gofeed"
)

type FeedParser struct {
	client *client.HttpClient
	parser *gofeed.Parser
}

func NewFeedParser(httpClient *client.HttpClient) *FeedParser {
	return &FeedParser{
		client: httpClient,
		parser: gofeed.NewParser(),
	}
}

func (p *FeedParser) ParseFromURL(ctx context.Context, URL string) (*gofeed.Feed, error) {
	resp, err := p.client.GetWithContext(ctx, URL)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.New(resp.Status)
	}
	return p.parser.Parse(resp.Body)
}
