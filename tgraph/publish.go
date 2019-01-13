package tgraph

import (
	"fmt"
	"log"
)

func PublishItem(title string, html string) string {

	url, _ := PublisHtml(title, html)

	return url
}

func PublisHtml(title string, html string) (string, error) {
	//log.Println(html)
	// CreatePage
	if page, err := client.CreatePageWithHTML(title, authorName, authorUrl, html, true); err == nil {
		//log.Printf("> CreatePage result: %#+v", page)
		log.Printf("> Created page url: %s", page.URL)

		fmt.Println(page)
		return page.URL, err
	} else {
		log.Printf("* CreatePage error: %s", err)
		return "", nil
	}
}
