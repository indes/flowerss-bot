package task

import (
	"github.com/indes/rssflow/bot"
	"github.com/indes/rssflow/config"
	"github.com/indes/rssflow/model"
	"time"
)

func init() {

}

func Update() {
	for {
		sources := model.GetSubscribedSources()
		for _, source := range sources {
			c, err := source.GetNewContents()
			if err == nil {
				subs := model.GetSubscriberBySource(&source)
				bot.BroadNews(&source, subs, c)
			}
		}
		time.Sleep(time.Duration(config.UpdateInterval) * time.Minute)
	}
}
