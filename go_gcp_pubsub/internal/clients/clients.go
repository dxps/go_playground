package clients

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/dxps/go_playground_go_gcp_pubsub/internal/endpoints"
	"google.golang.org/api/option"
)

func IsEndpointValid(endpoint string) bool {

	for _, e := range endpoints.PubSubEndpoints {
		if e == endpoint {
			return true
		}
	}
	return false
}

func InitClient(endpoint, projectID string) (*pubsub.Client, error) {

	if !IsEndpointValid(endpoint) {
		return nil, fmt.Errorf("Provided endpoint is not valid")
	}
	ctx := context.Background()
	opt := option.WithEndpoint(endpoint)
	return pubsub.NewClient(ctx, projectID, opt)
}
