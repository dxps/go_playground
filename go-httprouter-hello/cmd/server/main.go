package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dxps/go_playground/go-httprouter-hello/internal/app"
)

func main() {

	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("App init failed: %s", err)
	}

	gracefulShutdown(app)

	log.Println("Shutdown complete.")
}

func gracefulShutdown(app *app.App) {

	wg := &sync.WaitGroup{}
	wg.Add(1)
	if err := app.Start(wg); err != nil {
		log.Fatalf("Startup error: %v\n", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	// Waiting for shutdown signal.
	<-signalChan
	log.Println("\nShutting down ...")

	stopCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFn()

	app.Stop(stopCtx)

	wg.Wait()
}
