## go-expvar-metrics

A minimal example of using `expvar` for exposing metrics and showing them nicely using `expvarmon`.

### Run

Use provided scripts for starting up both parts: the app and the monitor.

### Get something out

Since the _app_ is just a bare minimal HTTP server that increases a `metrics.requests` variable, hit http://localhost:8080 for multiple times to see how this metric increases and get the free histogram of these changes in the `expvarmon` output (execution of `mon.sh`).

So, you can use for example `curl localhost:8080` in a third shell terminal to increase the value of that metric.

