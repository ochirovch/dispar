package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/pubsub"
)

func main() {
	subID := "saveResults"
	projectID := os.Getenv("DISPARPROJECT")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("pubsub.NewClient: %v", err)
		return
	}
	defer client.Close()

	sub := client.Subscription(subID)

	// Turn on synchronous mode. This makes the subscriber use the Pull RPC rather
	// than the StreamingPull RPC, which is useful for guaranteeing MaxOutstandingMessages,
	// the max number of messages the client will hold in memory at a time.
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 1

	// Receive messages for 5 seconds.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message, 1)
	defer close(cm)
	// Handle individual messages in a goroutine.
	go func() {
		count, overall := 0, 0
		for msg := range cm {
			if overall == 0 {
				// add here rows
			}
			fmt.Println(string(msg.Data))
			data := make(map[string]string)
			err := json.Unmarshal(msg.Data, data)
			if err != nil {
				log.Println(err)
				continue
			}
			msg.Ack()
			if count == 1000 {
				// write bucket here

				// spreadsheets delay
				time.Sleep(1 * time.Second)
				count = 0
			}
			count++
			overall++
		}
		if count > 1 {
			// write rest of bucket here
		}

	}()

	// Receive blocks until the passed in context is done.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil && status.Code(err) != codes.Canceled {
		log.Println(err)
		return
	}

}
