package DataCollector

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
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

	projectID := os.Getenv("DISPARPROJECT")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	iter := client.Collection("DisparSettings").Where("ProjectName", "==", attrs["project"]).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(doc.Data())
	dataselector := doc.Data()["DataSelector"].(string)
	dataselectors := doc.Data()["DataSelectors"].(map[string]interface{})

	// Find and visit all links
	c.OnHTML(dataselector, func(e *colly.HTMLElement) {
		results := make(map[string]string)
		for k, v := range dataselectors {
			results[k] = e.ChildText(v.(string))
		}
		log.Printf("%+v\n", results)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(string(m.Data))
	return nil
}
