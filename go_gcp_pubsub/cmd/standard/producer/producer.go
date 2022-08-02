package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/dxps/go_playground_go_gcp_pubsub/internal/clients"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/data"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/endpoints"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/produce"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/topics"
)

func main() {

	projectID, topicID, endpoint := "", "", ""
	numberOfEvents := uint(0)
	delay := uint(0)

	flag.StringVarP(&projectID, "projectID", "p", "tbd-project-id", "The Project ID")
	flag.StringVarP(&topicID, "topicID", "t", "tbd-topic-id", "The Topic ID")
	flag.StringVarP(&endpoint, "endpoint", "e", endpoints.PUBSUB_GLOBAL_ENDPOINT, "The Pub/Sub service endpoint.")
	flag.UintVarP(&numberOfEvents, "numberOfEvents", "n", 10, "Number of events to publish")
	flag.UintVarP(&delay, "delay", "d", 0, "Delay (in milliseconds) between publishing an event")

	flag.Parse()

	log.Printf(`Starting up using:
	       endpoint: '%s'
	     project ID: '%s'
	       topic ID: '%s'
     numberOfEvents: %d
	          delay: %d
	`, endpoint, projectID, topicID, numberOfEvents, delay)

	log.Printf("Using projectID: '%s', topicID: '%s', publishing %d events at an interval of %d milliseconds.",
		projectID, topicID, numberOfEvents, delay)

	client, err := clients.InitClient(endpoint, projectID)
	if err != nil {
		log.Fatalf("Failed to create PubSub client: %v", err)
	}

	topic, err := topics.InitTopic(client, topicID)
	if err != nil {
		log.Printf("Cannot use topic due to: %v", err)
		log.Printf("Trying to create it...")
		topic, err = topics.CreateTopic(client, topicID)
		if err != nil {
			log.Fatalf("Failed to use topic due to: %v. Existing now.", err)
		}
	}

	var wg sync.WaitGroup
	idChan := make(chan string, numberOfEvents)
	errChan := make(chan error, numberOfEvents)
	msgs := make(map[uint][]byte, numberOfEvents)

	ctx := context.Background()
	start := time.Now()
	sid := uint(start.Nanosecond())

	log.Println("Preparing the messages ...")
	obj := data.NewMyData()
	for n := uint(0); n < numberOfEvents; n++ {
		obj.ID = sid + n

		data, err := json.Marshal(obj)
		if err != nil {
			log.Fatalf("Failed to marshal object: %v", err)
		}
		msgs[obj.ID] = data
	}

	log.Println("Starting the publishing ...")

	for n := uint(0); n < numberOfEvents; n++ {

		if n > 0 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
		produce.PublishBytesAsyncRes(ctx, topic, msgs[n], &wg, idChan, errChan)
		if err != nil {
			log.Fatalf("Failed to publish msg: %v due to: %v", msgs[n], err)
		}
	}

	wg.Wait()
	duration := time.Since(start)
	log.Printf("Publishing %d events took %v\n", numberOfEvents, duration)

}
