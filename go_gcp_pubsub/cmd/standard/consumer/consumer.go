package main

import (
	"context"
	"log"

	flag "github.com/spf13/pflag"

	"github.com/dxps/go_playground_go_gcp_pubsub/internal/clients"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/consume"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/endpoints"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/topics"
)

func main() {

	var projectID, topicID, endpoint, subscriptionID string

	flag.StringVarP(&projectID, "projectID", "p", "some-project", "ID of the GCP project")
	flag.StringVarP(&topicID, "topicID", "t", "test-topic", "name of the topic")
	flag.StringVarP(&endpoint, "endpoint", "e", endpoints.PUBSUB_GLOBAL_ENDPOINT, "The Pub/Sub service endpoint.")
	flag.StringVarP(&subscriptionID, "subscriptionID", "s", "test-subscription", "name of the subscription")

	flag.Parse()

	log.Printf(`Starting up using:
	       endpoint: '%s'
	     project ID: '%s'
	       topic ID: '%s'
	subscription ID: '%s'
	`, endpoint, projectID, topicID, subscriptionID)

	client, err := clients.InitClient(endpoint, projectID)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}

	topic, err := topics.InitTopic(client, topicID)
	if err != nil {
		log.Printf("Cannot use topic due to: %v", err)
		log.Printf("Trying to create it ...")
		topic, err = topics.CreateTopic(client, topicID)
		if err != nil {
			log.Fatalf("Failed to use topic due to: %v. Existing now.", err)
		}
	}

	sub, err := consume.CreateSubscriptionWithOrdering(client, topic, subscriptionID)
	if err != nil {
		_ = client.Close()
		log.Fatalf("Failed to create subscription with ordering due to: %v", err)
	}
	sc, _ := sub.Config(context.Background())
	log.Printf("Created the subscription. Config: %+v", sc)
	log.Println("Ready to receive messages.")

	if err := consume.ReceiveMessages(sub); err != nil {
		log.Println(err)
	}

}
