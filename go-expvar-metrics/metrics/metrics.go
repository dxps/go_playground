package metrics

import (
	"expvar"
	"fmt"
	"runtime"
)

var metrics *expvar.Map

func init() {
	metrics = expvar.NewMap("metrics")
}

func IncreaseRequests() {
	metrics.Add("requests", 1)
}

// -----------------------------------------------------------------------------------
// Example:
//  NewStatCount("stats")
//  ..
//  Increment("active_ws_users")
//  defer Decrement("active_ws_users")
//  doWebSocket()

// Recommend using expvarmon (see https://github.com/divan/expvarmon )
//   expvarmon  -ports="http://localhost:9000" -i=1s \
//   -vars="goroutines,stats.auth_failures,stats.open_ws,stats.rabbit_messages_rx,stats.rabbit_reconnections, \
//   mem:memstats.Alloc,mem:memstats.Sys,mem:memstats.HeapAlloc,mem:memstats.HeapInuse"

// NewStatCount sets up a stat counter
// See http://go-talks.appspot.com/github.com/sajari/talks/201610/simplifying-storage/storage.slide#36 for more
// or http://www.mikeperham.com/2014/12/17/expvar-metrics-for-golang/ or http://blog.ralch.com/tutorial/golang-metrics-with-expvar/

var stats *expvar.Map

func NewStatCount(statName string) {
	stats = expvar.NewMap(statName)
	stats.Set("auth_failures", new(expvar.Int))
	stats.Set("open_ws", new(expvar.Int))
	stats.Set("rabbit_messages_rx", new(expvar.Int))
	stats.Set("rabbit_reconnections", new(expvar.Int))

	// Export goroutines
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return fmt.Sprintf("%d", runtime.NumGoroutine())
	}))
}

// Increment a certain stat
func Increment(stat string) {
	if stats == nil {
		return
	}
	stats.Add(stat, 1)
}

// Decrement a certain stat
func Decrement(stat string) {
	if stats == nil {
		return
	}
	stats.Add(stat, -1)
}

// SetInt sets a particular particular stat to a specific integer value
func SetInt(stat string, n int64) {
	if stats == nil {
		return
	}
	stats.Get(stat).(*expvar.Int).Set(n)
}
