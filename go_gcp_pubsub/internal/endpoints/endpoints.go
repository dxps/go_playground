package endpoints

// The list of all the endpoints can be found at:
// https://cloud.google.com/pubsub/docs/reference/service_apis_overview#service_endpoints
const (
	// GCP Pub/Sub global endpoint
	PUBSUB_GLOBAL_ENDPOINT = "pubsub.googleapis.com:443"

	// Netherlands regional endpoint.
	PUBSUB_EUROPE_WEST_4_ENDPOINT = "europe-west4-pubsub.googleapis.com:443"
)

// The list of all the endpoints, used for validating the user's input (the -e {endpoint} flag).
var PubSubEndpoints = []string{
	PUBSUB_GLOBAL_ENDPOINT,
	PUBSUB_EUROPE_WEST_4_ENDPOINT,
}
