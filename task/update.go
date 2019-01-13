package task

import (
	"github.com/indes/rssflow/bot"
	"github.com/indes/rssflow/model"
	"time"
)

func init() {

}

func Update() {
	for {
		sources := model.GetSources()
		for _, source := range sources {
			c, _ := source.GetNewContents()
			subs := model.GetSubscriberBySource(&source)
			bot.BroadNews(subs, c)
		}
		time.Sleep(10 * time.Minute)
	}
}
