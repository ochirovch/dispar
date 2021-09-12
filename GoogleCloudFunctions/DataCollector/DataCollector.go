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
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

// HelloPubSub consumes a Pub/Sub message.
func DataCollector(ctx context.Context, m PubSubMessage) error {

	c := colly.NewCollector()
	attrs := m.Attributes
	log.Println(attrs["project"])

	iter := client.Collection("disparSettings").Where("ProjectName", "==", attrs["project"]).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(doc.Data())
	dataselector := doc.Data()["DataSelector"]

	// Find and visit all links
	c.OnHTML(dataselector, func(e *colly.HTMLElement) {
		log.Println(e.Text)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(string(m.Data))
	m.Ack()
	return nil
}
