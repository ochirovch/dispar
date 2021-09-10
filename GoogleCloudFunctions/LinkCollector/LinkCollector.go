package LinkCollector

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/functions/metadata"
	pubsub "cloud.google.com/go/pubsub"
	"github.com/gocolly/colly/v2"
	//pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

// FirestoreEvent is the payload of a Firestore event.
// Please refer to the docs for additional information
// regarding Firestore events.
type FirestoreEvent struct {
	OldValue FirestoreValue `json:"oldValue"`
	Value    FirestoreValue `json:"value"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	Fields interface{} `json:"fields"`
}

func (v FirestoreValue) getIntegerValue(name string) (int, error) {
	fields, ok := v.Fields.(map[string]interface{})
	mapped, ok := fields[name].(map[string]interface{})
	if !ok {
		return 0, errors.New(fmt.Errorf("Error extracting value %s from %+v", name, fields).Error())
	}
	strValue, ok := mapped["integerValue"].(string)
	if !ok {
		return 0, errors.New(fmt.Errorf("Error extracting value %s from %+v", name, fields).Error())
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// getStringValue extracts a string value from a Firestore value
func (v FirestoreValue) getStringValue(name string) (string, error) {
	fields, ok := v.Fields.(map[string]interface{})
	mapped, ok := fields[name].(map[string]interface{})
	if !ok {
		return "", errors.New(fmt.Errorf("Error extracting value %s from %+v", name, fields).Error())
	}
	value, ok := mapped["stringValue"].(string)
	if !ok {
		return "", errors.New(fmt.Errorf("Error extracting value %s from %+v", name, fields).Error())
	}
	return value, nil
}

// HelloFirestore is triggered by a change to a Firestore document.
func LinkCollector(ctx context.Context, e FirestoreEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)

	PaginationURL, err := e.Value.getStringValue("PaginationURL")
	if err != nil {
		log.Println(err)
	}

	LinkSelector, err := e.Value.getStringValue("LinkSelector")
	if err != nil {
		log.Println(err)
	}

	PaginationStartPage, err := e.Value.getIntegerValue("PaginationStartPage")
	if err != nil {
		log.Println(err)
	}
	PaginationEndPage, err := e.Value.getIntegerValue("PaginationEndPage")
	if err != nil {
		log.Println(err)
	}

	projectID := os.Getenv("DISPARPROJECT")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
	}
	t := client.Topic("links")

	defer client.Close()

	c := colly.NewCollector()

	// add data to pub sub channel page project;url  LinkSelector
	// c.OnHTML(".news-item__content > a", func(e *colly.HTMLElement) {

	c.OnHTML(LinkSelector, func(e *colly.HTMLElement) {
		log.Println(e.Request.URL.Host + e.Attr("href"))
		sendUrlToPubSub(client, t, e.Request.URL.Host+e.Attr("href"))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := PaginationStartPage; i <= PaginationEndPage; i++ {
		fmt.Println(i)
		//c.Visit("https://www.chita.ru/news/?pg=" + strconv.Itoa(i))
		c.Visit(PaginationURL + strconv.Itoa(i))
	}
	return nil
}

func sendUrlToPubSub(client *pubsub.Client, topic *pubsub.Topic, url string) {

	ctx := context.Background()

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(url),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		log.Println(err)
	}
	log.Println(id)
}
