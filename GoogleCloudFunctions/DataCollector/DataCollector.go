package DataCollector

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gocolly/colly"
	"log"
	"os"
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

	projectID := os.Getenv("DISPARPROJECT")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	iter := client.Collection("disparSettings").Where("ProjectName", "==", attrs["project"]).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(doc.Data())
	dataselector := doc.Data()["DataSelector"].(string)

	// Find and visit all links
	c.OnHTML(dataselector, func(e *colly.HTMLElement) {
		log.Println(e.Text)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(string(m.Data))
	return nil
}
