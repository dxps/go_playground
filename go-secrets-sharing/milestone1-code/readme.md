## A Secrets Sharing App

This is a simple secrets sharing web app written in Go, using the standard library and no external dependencies.

### Run

Use `./run.sh` to start it.

### Usage

The following HTTP API operations are implemented:

- A health check<br/>
  Usage example: `curl -v localhost:9001/healthcheck`
- Adding a secret<br/>
  Usage example: `curl -v -X POST localhost:9001/secrets -H "Content-Type: application/json" -d '{ "plain_text": "shoop" }'`
- Getting a secret<br/>
  The _id_ used to retrieve the secret can be used only once.<br/>
  Any additional retries using it would get a "not found" (http status code 404) response.<br/>
  Usage example: `curl -v localhost:9001/secrets/991abd0371a608a298b01cba186f7c5c`

#### Notes

The added secrets are stored in memory and persisted on disk, thus have them available across restarts: at startup, it loads into memory the entries from the file.

The path to the secrets file can be specified using the `DATA_FILE_PATH` environment variable. See an example in the `run.sh` script. If this variable is not defined, then the default value of `./secrets` is used.
