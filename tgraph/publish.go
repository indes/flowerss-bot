package tgraph

import (
	"log"
)

func PublishItem(title string, html string) string {

	url, _ := PublisHtml(title, html)

	return url
}

func PublisHtml(title string, html string) (string, error) {

	// CreatePage
	if page, err := client.CreatePageWithHTML(title, authorName, authorUrl, html, true); err == nil {
		log.Printf("Created telegraph page url: %s", page.URL)
		return page.URL, err
	} else {
		log.Printf("Create telegraph page error: %s", err)
		return "", nil
	}
}
