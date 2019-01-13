package tgraph

import (
	"fmt"
	"github.com/SlyMarbo/rss"
	"log"
)

func PublishItems(items []*rss.Item) error {
	for _, item := range items {
		url, _ := PublisHtml(item.Title, item.Content)
		fmt.Println(url)
	}
	return nil
}

func PublisHtml(title string, html string) (string, error) {

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

		fmt.Println(page)
		return page.URL, err
	} else {
		log.Printf("* CreatePage error: %s", err)
		return "", nil
	}
}
