package main

import (
	"log"

	flag "github.com/spf13/pflag"

	"github.com/dxps/go_playground_go_gcp_pubsub/internal/client"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/subscription"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/topic"
)

func main() {

	projectID, topicID, subscriptionID := "", "", ""

	flag.StringVarP(&projectID, "projectID", "p", "some-project", "ID of the GCP project")
	flag.StringVarP(&topicID, "topicID", "t", "test-topic", "name of the topic")
	flag.StringVarP(&subscriptionID, "subscriptionID", "s", "example-subscription", "name of subcriber")

	flag.Parse()

	log.Printf("Using projectID: '%s', topicID: '%s', subscriptionID: '%s'.", projectID, topicID, subscriptionID)

	client, err := client.InitClient(projectID)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}

	topic, err := topic.InitTopic(client, topicID)
	if err != nil {
		log.Fatalf("Failed to use topic: %v", err)
	}

	sub := subscription.InitSubscription(client, topic, subscriptionID)

	if err := subscription.ReceiveMessages(sub); err != nil {
		log.Println(err)
	}

}