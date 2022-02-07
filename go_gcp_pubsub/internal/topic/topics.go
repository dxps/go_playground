package topic

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func InitTopic(c *pubsub.Client, name string) (*pubsub.Topic, error) {

	ctx := context.Background()
	t := c.Topic(name)
	ok, err := t.Exists(ctx)
	if !ok {
		return nil, fmt.Errorf("Topic '%s' does not exist.", name)
	} else if err != nil {
		return nil, err
	}
	return t, nil
}
