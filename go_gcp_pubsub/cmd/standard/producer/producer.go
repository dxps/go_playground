package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/dxps/go_playground_go_gcp_pubsub/internal/client"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/data"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/produce"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/topic"
)

func main() {

	projectID, topicID := "", ""
	eventsCount := uint(0)

	flag.StringVarP(&projectID, "projectID", "p", "tbd-project-id", "The Project ID")
	flag.StringVarP(&topicID, "topicID", "t", "tbd-topic-id", "The Topic ID")
	flag.UintVarP(&eventsCount, "eventsCount", "e", 10, "Number of events to publish")

	flag.Parse()

	log.Printf("Using projectID: '%s', topicID: '%s'.", projectID, topicID)

	client, err := client.InitClient(projectID)
	if err != nil {
		log.Fatalf("Failed to create PubSub client: %v", err)
	}

	topic, err := topic.InitTopic(client, topicID)
	if err != nil {
		log.Fatalf("Failed to use topic: %v", err)
	}

	obj := data.NewMyData()

	var wg sync.WaitGroup
	idChan := make(chan string, eventsCount)
	errChan := make(chan error, eventsCount)
	msgs := make([][]byte, 0)

	ctx := context.Background()
	start := time.Now()
	sid := uint(start.Nanosecond())

	log.Println("Preparing the messages ...")
	for n := uint(0); n < eventsCount; n++ {
		obj.ID = sid + n

		data, err := json.Marshal(obj)
		if err != nil {
			log.Fatalf("Failed to marshal object: %v", err)
		}
		msgs = append(msgs, data)
	}

	log.Println("Starting the publishing ...")

	for n := uint(0); n < eventsCount; n++ {

		produce.PublishBytesAsyncRes(ctx, topic, msgs[n], &wg, idChan, errChan)
		if err != nil {
			log.Fatalf("Failed to publish msg: %v due to: %v", msgs[n], err)
		}
	}

	wg.Wait()
	duration := time.Since(start)
	log.Printf("Publishing %d events took %v\n", eventsCount, duration)

}
