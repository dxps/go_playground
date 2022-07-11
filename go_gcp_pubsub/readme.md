# Some GCP Pub/Sub Samples & Tests

This is a basic sample to test the behaviour of GCP's Pub/Sub producing and consuming (receiving from subscription) messages.

<br/>

## Setup

To use it, you need to have a _GCP Service Account_ with the appropriate Pub/Sub roles, such as `Pub/Sub Admin`.<br/>
Create a service account key with type JSON, download the generated file and set `GOOGLE_APPLICATION_CREDENTIALS` as environment variable with the value of (full or relative) path to that JSON file.

To export it in the current shell (to have it visibile in any child shell) use:
```shell
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/gcp-service-account.json
```

<br/>

### Usage

Remember that both producer and subscriber need to get access to a service account credentials (key) through the aforementioned environment variable. Otherwise, at startup it will fail with:
```
Failed to create pubsub client: pubsub(publisher): google: could not find default credentials. See https://developers.google.com/accounts/docs/application-default-credentials for more information.
```


Start the subscriber:
```shell
go run cmd/standard/subscriber/main.go -p {project-id} -t {topic-id} -s {subscription-id}
```

Start the producer:
```shell
go run cmd/standard/producer/main.go -p {project-id} -t {topic-id}
```

You should see both sides (production and consumption) showing off their parts. The order is not that important though.

<br/>
