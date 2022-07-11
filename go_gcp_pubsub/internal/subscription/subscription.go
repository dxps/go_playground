package subscription

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

func CreateSubscription(c *pubsub.Client, topic *pubsub.Topic, subID string) (*pubsub.Subscription, error) {

	ctx := context.Background()
	sub, err := c.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func InitSubscription(c *pubsub.Client, topic *pubsub.Topic, subID string) *pubsub.Subscription {

	return c.Subscription(subID)
}

func ReceiveMessages(sub *pubsub.Subscription) error {

	var mu sync.Mutex
	received := 0
	ctx := context.Background()
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		log.Printf("Got '%q' with id %s.\n", string(msg.Data), msg.ID)
		msg.Ack()
		received++
	})
	if err != nil {
		return fmt.Errorf("Failed to receive: %v", err)
	}
	return nil
}
