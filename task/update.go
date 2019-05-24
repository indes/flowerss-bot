package task

import (
	"github.com/indes/flowerss-bot/bot"
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/model"
	"time"
)

func init() {

}

func Update() {
	for {
		sources := model.GetSubscribedNormalSources()
		for _, source := range sources {
			c, err := source.GetNewContents()
			if err == nil {
				subs := model.GetSubscriberBySource(&source)
				bot.BroadNews(&source, subs, c)
			}
			if source.ErrorCount >= config.ErrorThreshold {
				bot.BroadSourceError(&source)
			}
		}
		time.Sleep(time.Duration(config.UpdateInterval) * time.Minute)
	}
}
