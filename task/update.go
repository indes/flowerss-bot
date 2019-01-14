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
		sources := model.GetSubscribedSources()
		for _, source := range sources {
			c, _ := source.GetNewContents()
			subs := model.GetSubscriberBySource(&source)
			bot.BroadNews(&source, subs, c)
		}
		time.Sleep(10 * time.Minute)
	}
}
