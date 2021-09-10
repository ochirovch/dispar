package LinkCollector

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	pubsub "cloud.google.com/go/pubsub"
	"github.com/gocolly/colly/v2"
	//pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
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
	c.OnHTML(".news-item__content > a", func(e *colly.HTMLElement) {
		log.Println(e.Request.URL.Host + e.Attr("href"))
		sendUrlToPubSub(client, t, e.Request.URL.Host+e.Attr("href"))
	})

	c.OnHTML("html", func(e *colly.HTMLElement) { // Title
		//log.Println(string(e.Response.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := begin; i < end; i++ {
		fmt.Println(i)
		c.Visit("https://www.chita.ru/news/?pg=" + strconv.Itoa(i))
	}
}

func sendUrlToPubSub(client *pubsub.Client, topic *pubsub.Topic, url string) {

	ctx := context.Background()

	// result := topic.Publish(ctx, &pubsub.Message{
	// 	Data: []byte(url),
	// })

	// // The Get method blocks until a server-generated ID or
	// // an error is returned for the published message.
	// id, err := result.Get(ctx)
	// if err != nil {
	// 	// Error handling code can be added here.
	// 	log.Printf("Failed to publish: %v", err)
	// 	return
	// }
	// log.Printf("Published message; msg ID: %v\n", id)
	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(url),
		Attributes: map[string]string{
			"origin":   "golang",
			"username": "gcp",
		},
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		log.Println(err)
	}
	log.Println(id)
}
