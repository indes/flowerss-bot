package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/SlyMarbo/rss"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/fetch"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/pkg/client"
)

// RssUpdateObserver Rss Update observer
type RssUpdateObserver interface {
	SourceUpdate(*model.Source, []*model.Content, []*model.Subscribe)
	SourceUpdateError(*model.Source)
}

// NewRssTask new RssUpdateTask
func NewRssTask() *RssUpdateTask {
	clientOpts := []client.HttpClientOption{
		client.WithTimeout(10 * time.Second),
	}
	if config.Socks5 != "" {
		clientOpts = append(clientOpts, client.WithProxyURL(fmt.Sprintf("socks5://%s", config.Socks5)))
	}

	if config.UserAgent != "" {
		clientOpts = append(clientOpts, client.WithUserAgent(config.UserAgent))
	}
	httpClient := client.NewHttpClient(clientOpts...)

	return &RssUpdateTask{
		observerList: []RssUpdateObserver{},
		httpClient:   httpClient,
	}
}

// RssUpdateTask rss更新任务
type RssUpdateTask struct {
	observerList []RssUpdateObserver
	isStop       atomic.Bool
	httpClient   *client.HttpClient
}

// Register 注册rss更新订阅者
func (t *RssUpdateTask) Register(observer RssUpdateObserver) {
	t.observerList = append(t.observerList, observer)
}

// Stop task
func (t *RssUpdateTask) Stop() {
	t.isStop.Store(true)
}

// Start run task
func (t *RssUpdateTask) Start() {
	if config.RunMode == config.TestMode {
		return
	}

	t.isStop.Store(false)
	go func() {
		for {
			if t.isStop.Load() == true {
				zap.S().Info("RssUpdateTask stopped")
				return
			}

			sources := model.GetSubscribedNormalSources()
			for _, source := range sources {
				if !source.NeedUpdate() {
					continue
				}

				newContents, err := t.getSourceNewContents(source)
				if err != nil {
					if source.ErrorCount >= config.ErrorThreshold {
						t.notifyAllObserverErrorUpdate(source)
					}
					continue
				}

				if len(newContents) > 0 {
					subs := model.GetSubscriberBySource(source)
					t.notifyAllObserverUpdate(source, newContents, subs)
				}
			}

			time.Sleep(time.Duration(config.UpdateInterval) * time.Minute)
		}
	}()
}

// getSourceNewContents 获取rss新内容
func (t *RssUpdateTask) getSourceNewContents(source *model.Source) ([]*model.Content, error) {
	zap.S().Debugw("fetch source updates", "source", source)

	var newContents []*model.Content
	feed, err := rss.FetchByFunc(fetch.FetchFunc(t.httpClient), source.Link)
	if err != nil {
		zap.S().Errorw("unable to fetch SourceUpdate", "error", err, "source", source)
		source.AddErrorCount()
		return nil, err
	}

	source.EraseErrorCount()
	items := feed.Items
	for _, item := range items {
		c, isBroad, _ := model.GenContentAndCheckByFeedItem(source, item)
		if !isBroad {
			newContents = append(newContents, c)
		}
	}
	return newContents, nil
}

// notifyAllObserverUpdate notify all rss SourceUpdate observer
func (t *RssUpdateTask) notifyAllObserverUpdate(
	source *model.Source, newContents []*model.Content, subscribes []*model.Subscribe,
) {
	wg := sync.WaitGroup{}
	for _, observer := range t.observerList {
		wg.Add(1)
		go func(o RssUpdateObserver) {
			defer wg.Done()
			o.SourceUpdate(source, newContents, subscribes)
		}(observer)
	}
	wg.Wait()
}

// notifyAllObserverErrorUpdate notify all rss error SourceUpdate observer
func (t *RssUpdateTask) notifyAllObserverErrorUpdate(source *model.Source) {
	wg := sync.WaitGroup{}
	for _, observer := range t.observerList {
		wg.Add(1)
		go func(o RssUpdateObserver) {
			defer wg.Done()
			o.SourceUpdateError(source)
		}(observer)
	}
	wg.Wait()
}
