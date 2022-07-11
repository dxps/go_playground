package topic

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func InitTopic(c *pubsub.Client, topicID string) (*pubsub.Topic, error) {

	ctx := context.Background()
	t := c.Topic(topicID)
	ok, err := t.Exists(ctx)
	if !ok {
		return nil, fmt.Errorf("Topic '%s' does not exist or it's not accessible.", topicID)
	} else if err != nil {
		return nil, err
	}
	return t, nil
}
