package task

import (
	"github.com/indes/flowerss-bot/bot"
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/model"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"sync"
	"time"
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
	return &RssUpdateTask{
		observerList: []RssUpdateObserver{},
	}
}

// RssUpdateTask rss更新任务
type RssUpdateTask struct {
	observerList []RssUpdateObserver
	isStop       atomic.Bool
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

// Stop stop task
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

				newContents, err := source.GetNewContents()
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
