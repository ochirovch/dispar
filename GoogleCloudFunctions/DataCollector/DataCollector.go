package DataCollector

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
func DataCollector(ctx context.Context, m PubSubMessage) error {

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".masha_index", func(e *colly.HTMLElement) {
		log.Println(e.Text)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		log.Println(e.Text)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(string(m.Data))
	return nil
}
