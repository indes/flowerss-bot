package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/SlyMarbo/rss"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/indes/flowerss-bot/internal/bot"
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/fetch"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/pkg/client"
)

func init() {
	task := NewRssTask()
	task.Register(&telegramBotRssUpdateObserver{})
	registerTask(task)
}

// RssUpdateObserver Rss update observer
type RssUpdateObserver interface {
	update(*model.Source, []*model.Content, []*model.Subscribe)
	errorUpdate(*model.Source)
	id() string
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

// Name 任务名称
func (t *RssUpdateTask) Name() string {
	return "RssUpdateTask"
}

// Register 注册rss更新订阅者
func (t *RssUpdateTask) Register(observer RssUpdateObserver) {
	t.observerList = append(t.observerList, observer)
}

// Register 注销rss更新订阅者
func (t *RssUpdateTask) Deregister(removeObserver RssUpdateObserver) {
	for i, observer := range t.observerList {
		if observer.id() == removeObserver.id() {
			t.observerList = append(t.observerList[:i], t.observerList[i+1:]...)
			return
		}
	}
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
		zap.S().Errorw("unable to fetch update", "error", err, "source", source)
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

// notifyAllObserverUpdate notify all rss update observer
func (t *RssUpdateTask) notifyAllObserverUpdate(
	source *model.Source, newContents []*model.Content, subscribes []*model.Subscribe) {

	wg := sync.WaitGroup{}
	for _, observer := range t.observerList {
		wg.Add(1)
		go func(o RssUpdateObserver) {
			defer wg.Done()
			o.update(source, newContents, subscribes)
		}(observer)
	}
	wg.Wait()
}

// notifyAllObserverErrorUpdate notify all rss error update observer
func (t *RssUpdateTask) notifyAllObserverErrorUpdate(source *model.Source) {
	wg := sync.WaitGroup{}
	for _, observer := range t.observerList {
		wg.Add(1)
		go func(o RssUpdateObserver) {
			defer wg.Done()
			o.errorUpdate(source)
		}(observer)
	}
	wg.Wait()
}

type telegramBotRssUpdateObserver struct {
}

func (o *telegramBotRssUpdateObserver) update(
	source *model.Source, newContents []*model.Content, subscribes []*model.Subscribe) {
	zap.S().Debugf("%v receiving [%d]%v update", o.id(), source.ID, source.Title)
	bot.BroadcastNews(source, subscribes, newContents)
}

func (o *telegramBotRssUpdateObserver) errorUpdate(source *model.Source) {
	zap.S().Debugf("%v receiving [%d]%v error update", o.id(), source.ID, source.Title)
	bot.BroadcastSourceError(source)
}

func (o *telegramBotRssUpdateObserver) id() string {
	return "telegramBotRssUpdateObserver"
}
