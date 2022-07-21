package consume

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/data"
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
	var currID uint
	ctx := context.Background()
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		obj, err := data.NewMyDataFromBytes(msg.Data)
		if err != nil {
			log.Printf("Error on receive: %v", err)
		} else {
			log.Printf("Received objID: %d", obj.ID)
			if currID == 0 {
				currID = obj.ID
			} else {
				if currID != obj.ID-1 {
					log.Printf("Got an unordered message: currID=%v objID=%v", currID, obj.ID)
				}
				currID = obj.ID
			}
		}
		msg.Ack()
		received++
	})
	if err != nil {
		return fmt.Errorf("Failed to receive: %v", err)
	}
	return nil
}
