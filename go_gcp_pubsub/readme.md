## GCP Pub/Sub

This is a basic sample to test the behaviour of GCP's Pub/Sub producing and consuming (receiving from subscription) messages.

<br/>

### The Service Account

To use it, you need to have a _GCP Service Account_ with the appropriate Pub/Sub roles.
And pass to the sample a reference to the JSON file that contains the private key.

To export it in the current shell (and have it visibile in any child shell) use:
```shell
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/gcp-service-account.json
```

### Usage

First of all, as mentioned before, both producer and subscriber need to get the service account credentials.
Just use that `export GOOGLE_APPLICATION_CREDENTIALS` or make sure `GOOGLE_APPLICATION_CREDENTIALS` exists as an environment variable.

Otherwise, the startup will fail with:
```
Failed to create pubsub client: pubsub(publisher): google: could not find default credentials. See https://developers.google.com/accounts/docs/application-default-credentials for more information.
```


Then start the subscriber:
```shell
go run cmd/subcriber/main.go -p your-existing-project-id -t your-existing-topic -s your-existing-subscription
```

And start the producer:
```shell
go run cmd/subcriber/main.go -p your-existing-project-id -t your-existing-topic
```

You should see both sides (producer's production and subscriber's consumption) showing off their parts.

<br/>