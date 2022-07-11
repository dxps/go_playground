package main

import (
	"encoding/json"
	"log"

	flag "github.com/spf13/pflag"

	"github.com/dxps/go_playground_go_gcp_pubsub/internal/client"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/publish"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/topic"
)

func main() {

	projectID, topicID := "", ""

	flag.StringVarP(&projectID, "projectID", "p", "some-project", "project ID")
	flag.StringVarP(&topicID, "topicID", "t", "test-topic", "topic ID")

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

	obj := struct {
		SomeTestID   int    `json:"someTestID"`
		SomeTestName string `json:"someTestName"`
	}{
		SomeTestName: "marile testing",
	}

	// Publishing two messages.
	for n := 0; n <= 1; n++ {
		obj.SomeTestID = n

		data, err := json.Marshal(obj)
		if err != nil {
			log.Fatalf("Failed to marshal object: %v", err)
		}
		msgID, err := publish.PublishBytes(topic, data)
		if err != nil {
			log.Fatalf("Failed to publish: %v", err)
		}
		log.Printf("Published %+v as msg with ID: %s", obj, msgID)
	}

}
