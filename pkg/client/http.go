package client

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type HttpClientOptions struct {
	UserAgent string
	Timeout   time.Duration
	ProxyURL  string
}

func NewHttpClientOptions() *HttpClientOptions {
	return &HttpClientOptions{Timeout: time.Second}
}

type HttpClientOption func(opts *HttpClientOptions)

func WithUserAgent(userAgent string) HttpClientOption {
	return func(o *HttpClientOptions) {
		o.UserAgent = userAgent
	}
}

func WithTimeout(timeout time.Duration) HttpClientOption {
	return func(o *HttpClientOptions) {
		o.Timeout = timeout
	}
}

func WithProxyURL(url string) HttpClientOption {
	return func(opts *HttpClientOptions) {
		opts.ProxyURL = url
	}
}

type HttpClient struct {
	client    *http.Client
	userAgent string
}

func NewHttpClient(opts ...HttpClientOption) *HttpClient {
	o := NewHttpClientOptions()
	for _, opt := range opts {
		opt(o)
	}

	transport := &http.Transport{}
	if o.ProxyURL != "" {
		proxyURL, _ := url.Parse(o.ProxyURL)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Timeout:   o.Timeout,
		Transport: transport,
	}

	return &HttpClient{
		client:    client,
		userAgent: o.UserAgent,
	}
}

func (c *HttpClient) GetWithContext(ctx context.Context, url string, opts ...HttpClientOption) (*http.Response, error) {
	o := NewHttpClientOptions()
	for _, opt := range opts {
		opt(o)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if o.UserAgent != "" {
		req.Header.Set("User-Agent", o.UserAgent)
	}
	return c.client.Do(req)
}

func (c *HttpClient) Get(url string, opts ...HttpClientOption) (*http.Response, error) {
	return c.GetWithContext(context.Background(), url, opts...)
}

// Client method returns the current `http.Client` used by HttpClient.
func (c *HttpClient) Client() *http.Client {
	return c.client
}
