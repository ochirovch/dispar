package LinkCollector

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/gocolly/colly/v2"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

const (
	begin = 1
	end   = 6
)

func LinkCollector(w http.ResponseWriter, req *http.Request) {
	projectID := os.Getenv("DISPARPROJECT")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
	}
	t := client.Topic("links")

	defer client.Close()

	c := colly.NewCollector()

	// add data to pub sub channel page project;url
	c.OnHTML(".news-item__body > a", func(e *colly.HTMLElement) {
		// e.Request.Visit(e.Attr("href"))
		sendUrlToPubSub(client, t, e.Attr("href"))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := begin; i < end; i++ {
		fmt.Println(i)
		c.Visit("https://www.chita.ru/news/?pg=" + strconv.Itoa(i))
	}
}

func sendUrlToPubSub(client *pubsub.Client, topic *pubsub.Topic, url string) {

	var totalErrors uint64
	ctx := context.Background()

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(url),
	})

	// The Get method blocks until a server-generated ID or
	// an error is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		// Error handling code can be added here.
		log.Printf("Failed to publish: %v", err)
		atomic.AddUint64(&totalErrors, 1)
		return
	}
	log.Printf("Published message; msg ID: %v\n", id)
}
