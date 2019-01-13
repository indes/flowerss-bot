package rss

import (
	"github.com/indes/rssflow/model"
	"log"
	"time"
)

func init() {

}

func Update() {
	for {
		sources := model.GetSources()
		for _, source := range sources {
			c, _ := source.GetNewContents()
			log.Println(c)
		}
		log.Println(len(sources))
		//log.Println(time.Now())
		time.Sleep(10 * time.Second)

	}
}
