package publish

import (
	"context"
	"sync"

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

// Publish the `data` bytes as message data, and gets the result in an async fashion.
// This means, publishing does not block waiting for the server-generated ID of the message.
// The message ID or any potential error are being sent back through the respective channel.
func PublishBytesAsyncRes(ctx context.Context, topic *pubsub.Topic, data []byte,
	wg *sync.WaitGroup, idChan chan string, errChan chan error) {

	wg.Add(1)
	res := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	go func(res *pubsub.PublishResult) {
		defer wg.Done()
		// Block until the result is returned and a server-generated ID of the message is returned.
		id, err := res.Get(ctx)
		idChan <- id
		errChan <- err
		// fmt.Printf("Published, got msg id %v\n", id)
	}(res)
}

func PublishBytesWithOrdering(topic *pubsub.Topic, data []byte, orderingKey string) (msgID string, err error) {

	ctx := context.Background()
	res := topic.Publish(ctx, &pubsub.Message{
		Data:        data,
		OrderingKey: orderingKey,
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	return res.Get(ctx)
}

// Publish the `data` bytes as message data, with ordering key and gets the result in an async fashion.
// This means, publishing does not block waiting for the server-generated ID of the message.
// The message ID or any potential error are being sent back through the respective channel.
func PublishBytesWithOrderingAsyncRes(ctx context.Context, topic *pubsub.Topic, data []byte, orderingKey string,
	wg *sync.WaitGroup, idChan chan string, errChan chan error) {

	wg.Add(1)
	res := topic.Publish(ctx, &pubsub.Message{
		Data:        data,
		OrderingKey: orderingKey,
	})
	go func(res *pubsub.PublishResult) {
		defer wg.Done()
		// Block until the result is returned and a server-generated ID of the message is returned.
		id, err := res.Get(ctx)
		idChan <- id
		errChan <- err
		// fmt.Printf("Published, got msg id %v\n", id)
	}(res)
}
