package tgraph

import (
	"fmt"
	"github.com/SlyMarbo/rss"
	"github.com/meinside/telegraph-go"
	"log"
)

func PublishItems(items []*rss.Item) error {
	for _, item := range items {
		PublisHtml(item.Title, item.Content)
	}
	return nil
}

func PublisHtml(title string, html string) (error) {

	// CreatePage
	if page, err := client.CreatePageWithHTML(title, authorName, authorUrl, html, true); err == nil {
		log.Printf("> CreatePage result: %#+v", page)
		log.Printf("> Created page url: %s", page.URL)

		// GetPage
		if page, err := client.GetPage(page.Path, true); err == nil {
			log.Printf("> GetPage result: %#+v", page)
		} else {
			log.Printf("* GetPage error: %s", err)
		}

		page, err := telegraph.NewNodesWithHTML(html)

		fmt.Println(page)
		return err
	} else {
		log.Printf("* CreatePage error: %s", err)
		return nil
	}
}
