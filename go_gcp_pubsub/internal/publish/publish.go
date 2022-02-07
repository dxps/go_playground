package publish

import (
	"context"

	"cloud.google.com/go/pubsub"
)

func PublishBytes(topic *pubsub.Topic, data []byte) (msgID string, err error) {

	ctx := context.Background()
	res := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	return res.Get(ctx)
}
