package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"go.uber.org/atomic"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/feed"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/pkg/client"
)

// RssUpdateObserver Rss Update observer
type RssUpdateObserver interface {
	SourceUpdate(*model.Source, []*model.Content, []*model.Subscribe)
	SourceUpdateError(*model.Source)
}

// NewRssTask new RssUpdateTask
func NewRssTask(appCore *core.Core) *RssUpdateTask {
	return &RssUpdateTask{
		observerList: []RssUpdateObserver{},
		core:         appCore,
		feedParser:   appCore.FeedParser(),
		httpClient:   appCore.HttpClient(),
	}
}

// RssUpdateTask rss更新任务
type RssUpdateTask struct {
	observerList []RssUpdateObserver
	isStop       atomic.Bool
	core         *core.Core
	feedParser   *feed.FeedParser
	httpClient   *client.HttpClient
}

// Register 注册rss更新订阅者
func (t *RssUpdateTask) Register(observer RssUpdateObserver) {
	t.observerList = append(t.observerList, observer)
}

// Stop scheduler
func (t *RssUpdateTask) Stop() {
	t.isStop.Store(true)
}

// Start run scheduler
func (t *RssUpdateTask) Start() {
	if config.RunMode == config.TestMode {
		return
	}

	t.isStop.Store(false)
	go func() {
		for {
			if t.isStop.Load() {
				log.Info("RssUpdateTask stopped")
				return
			}

			sources, err := t.core.GetSources(context.Background())
			if err != nil {
				log.Errorf("get sources failed, %v", err)
				time.Sleep(time.Duration(config.UpdateInterval) * time.Minute)
				continue
			}
			for _, source := range sources {
				if source.ErrorCount >= config.ErrorThreshold {
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
					subs, err := t.core.GetSourceAllSubscriptions(
						context.Background(), source.ID,
					)
					if err != nil {
						log.Errorf("get subscriptions failed, %v", err)
						continue
					}
					t.notifyAllObserverUpdate(source, newContents, subs)
				}
			}

			time.Sleep(time.Duration(config.UpdateInterval) * time.Minute)
		}
	}()
}

// getSourceNewContents 获取rss新内容
func (t *RssUpdateTask) getSourceNewContents(source *model.Source) ([]*model.Content, error) {
	log.Debugf("fetch source [%d]%s update", source.ID, source.Link)

	rssFeed, err := t.feedParser.ParseFromURL(context.Background(), source.Link)
	if err != nil {
		log.Errorf("unable to fetch feed, source %#v, err %v", source, err)
		t.core.SourceErrorCountIncr(context.Background(), source.ID)
		return nil, err
	}
	t.core.ClearSourceErrorCount(context.Background(), source.ID)

	newContents, err := t.saveNewContents(source, rssFeed.Items)
	if err != nil {
		return nil, err
	}
	return newContents, nil
}

// saveNewContents generate content by fetcher item
func (t *RssUpdateTask) saveNewContents(
	s *model.Source, items []*gofeed.Item,
) ([]*model.Content, error) {
	var newItems []*gofeed.Item
	for _, item := range items {
		hashID := model.GenHashID(s.Link, item.GUID)
		exist, err := t.core.ContentHashIDExist(context.Background(), hashID)
		if err != nil {
			log.Errorf("check item hash id failed, %v", err)
		}

		if exist {
			// 已存在，跳过
			continue
		}
		newItems = append(newItems, item)
	}
	return t.core.AddSourceContents(context.Background(), s, newItems)
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
