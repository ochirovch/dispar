package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/pubsub"
)

var (
	readRange = "List1!A:H"
)

func main() {
	subID := "saveResults"
	SpreadsheetID := os.Getenv("DisparSpreadsheetID")
	projectID := os.Getenv("DisparProject")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("pubsub.NewClient: %v", err)
		return
	}
	defer client.Close()

	srvSpreadsheet := getExcelClient()

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
		data := make(map[string]string)
		var values [][]interface{}
		for msg := range cm {
			if overall == 0 {
				// add here rows
				//log.Println(tliker.Username)

			}
			fmt.Println(string(msg.Data))

			err := json.Unmarshal(msg.Data, &data)
			if err != nil {
				log.Println(err)
				continue
			}
			var value []interface{}
			for _, v := range data {
				value = append(value, v)
			}
			values = append(values, value)

			msg.Ack()
			if count == 1000 {
				// write bucket here
				addRowsToSpreadsheets(srvSpreadsheet, SpreadsheetID, values)
				// spreadsheets delay
				time.Sleep(1 * time.Second)
				count = 0
			}
			count++
			overall++
		}
		if count > 1 {
			addRowsToSpreadsheets(srvSpreadsheet, SpreadsheetID, values)
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

func getExcelClient() (srv *sheets.Service) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	srv, err = sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return srv
}

func addRowsToSpreadsheets(srv *sheets.Service, SpreadsheetID string, values [][]interface{}) {
	valueRange := sheets.ValueRange{
		MajorDimension: "ROWS",
		//		Range:          "Sales!",
		Values: values,
	}

	appendCall := srv.Spreadsheets.Values.Append(SpreadsheetID, readRange, &valueRange)
	appendCall.ValueInputOption("USER_ENTERED")
	_, err := appendCall.Do()
	if err != nil {
		log.Fatal(err)
	}

}
