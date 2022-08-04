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
Failed to create pubsub client: pubsub(publisher): google: could not find default credentials.
See https://developers.google.com/accounts/docs/application-default-credentials for more information.
```


#### Start the Consumer, Run the Producer(s)

Start the consumer using:<br/>
`go run cmd/standard/consumer/consumer.go -p {project-id} -t {topic-id} -s {subscription-id}`

Whether messages are being received in order (considering they were published using the same ordering key) or not it depends on how the subcription settings. There is nothing else to be done at the consumer level about it.

On the production side, things are a little bit different:
- If you want to send the messages in order, then an ordering key value must be specified for all the messages.<br/>
  That's where _ordering producer_ comes in.
- Otherwise, a _standard producer_ can be used: it publishes messages without any ordering key.

Run the producers using:
1. For the standard producer use:<br/>
   `go run cmd/standard/producer/producer.go -p {project-id} -t {topic-id} -n {number-of-messages|default=10} -d {delay|default=0ms}`
2. For the producer that uses an ordering key use:<br/>
   `go run cmd/ordering/producer/ordering_producer.go -p {project-id} -t {topic-id} -n {number-of-messages|default=10}`

You should see both sides (production and consumption) showing off their parts.

##### Endpoints

For both producers and consumer, you can also specify an additional flag `-e {endpoint}`. Btw, `-h` can be used for getting the usage flags and their details.

Currently, the defined endpoints are:
| endpoint | meaning |
| -------- | ------- |
| `pubsub.googleapis.com:443`              | Global Pub/Sub service endpoint |
| `europe-west4-pubsub.googleapis.com:443` | Regional (Netherlands) Pub/Sub service endpoint |

Producers and consumer can use different endpoints.

<br/>

### Stats

The followings represents just some execution snapshots while publishing just 1000 messages.

| number of messages | operation & type   | duration     |
| ------------------ | ------------------ | ------------ |
| 1_000              | published in order | 252.701501ms |
| 1_000              | published standard | 225.168448ms |


<br/>

### Todos

The _async result_ design must be improved as trying to publish 1M messages would probably take longer than the synchronous way
due to spawning too many goroutines. To be evaluated ...

<br/>