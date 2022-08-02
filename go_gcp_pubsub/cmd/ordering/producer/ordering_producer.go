package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
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

	var projectID, topicID, endpoint string
	var numberOfEvents uint

	flag.StringVarP(&projectID, "projectID", "p", "tbd-project-id", "The Project ID")
	flag.StringVarP(&topicID, "topicID", "t", "tbd-topic-id", "The Topic ID")
	flag.StringVarP(&endpoint, "endpoint", "e", endpoints.PUBSUB_GLOBAL_ENDPOINT, "The Pub/Sub service endpoint.")
	flag.UintVarP(&numberOfEvents, "numberOfEvents", "n", 10, "Number of events to publish")

	flag.Parse()

	log.Printf(`Starting up using:
	       endpoint: '%s'
	     project ID: '%s'
	       topic ID: '%s'
	 numberOfEvents: %d
	`, endpoint, projectID, topicID, numberOfEvents)

	client, err := clients.InitClient(endpoint, projectID)
	if err != nil {
		log.Fatalf("Failed to create PubSub client: %v", err)
	}

	topic, err := topics.InitTopic(client, topicID)
	if err != nil {
		log.Fatalf("Failed to use topic: %v", err)
	}

	// Enabling the message ordering at the topic level.
	topic.EnableMessageOrdering = true

	var wg sync.WaitGroup
	idChan := make(chan string, numberOfEvents)
	errChan := make(chan error, numberOfEvents)

	ctx := context.Background()
	start := time.Now()
	sid := uint(start.Nanosecond())
	msgs := make(map[string][]byte, 0)

	log.Println("Preparing the messages ...")
	obj := data.NewMyData()
	for n := uint(0); n < numberOfEvents; n++ {
		obj.ID = sid + n

		data, err := json.Marshal(obj)
		if err != nil {
			log.Fatalf("Failed to marshal object: %v", err)
		}
		msgs[fmt.Sprintf("%d", obj.ID)] = data
	}

	log.Println("Starting the publishing ...")
	keys := make([]string, 0)
	for k := range msgs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		_ = json.Unmarshal(msgs[k], &obj)
		log.Printf("Publishing msg with Data: %v and OrderingKey: %v", obj, k)
		produce.PublishBytesWithOrderingAsyncRes(ctx, topic, msgs[k], k, &wg, idChan, errChan)
		if err != nil {
			log.Fatalf("Failed to publish msg: %v due to: %v", string(msgs[k]), err)
		}
	}

	wg.Wait()
	duration := time.Since(start)
	log.Printf("Publishing %d events in order took %v\n", numberOfEvents, duration)

}
