package client

import (
	"context"

	"cloud.google.com/go/pubsub"
)

func InitClient(projectID string) (*pubsub.Client, error) {

	ctx := context.Background()
	return pubsub.NewClient(ctx, projectID)
}
