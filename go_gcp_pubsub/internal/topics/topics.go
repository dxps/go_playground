package topics

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
)

func CreateTopic(c *pubsub.Client, topicID string) (*pubsub.Topic, error) {

	ctx := context.Background()
	return c.CreateTopic(ctx, topicID)
}

func InitTopic(c *pubsub.Client, topicID string) (*pubsub.Topic, error) {

	ctx, cc := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cc()
	t := c.Topic(topicID)
	ok, err := t.Exists(ctx)
	if !ok {
		return nil, fmt.Errorf("Topic '%s' does not exist or it's not accessible.", topicID)
	} else if err != nil {
		return nil, err
	}
	return t, nil
}
