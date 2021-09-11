// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {
	//	log.Println(string(m.Data))

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".pane__button", func(e *colly.HTMLElement) {
		log.Println(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(string(m.Data))
	return nil
}
