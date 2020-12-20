package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/devisions/go-playground/go-directio/data"
	"github.com/ncw/directio"
)

func main() {
	// Preparing the graceful shutdown elements.
	stopCtx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	stopWg := &sync.WaitGroup{}
	stopWg.Add(2)

	dataCh := make(chan *data.SomeData, 1_000_000)

	go consumer(dataCh, stopCtx, stopWg)
	go reader(dataCh, stopCtx, stopWg)

	waitingForGracefulShutdown(cancelFn, stopWg)
}

func reader(dataCh chan *data.SomeData, stopCtx context.Context, stopWg *sync.WaitGroup) {
	filePath := "/Users/devisions/tmp/test_dio"
	block := directio.AlignedBlock(data.BlockSize)
	file, err := directio.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open file for reading. Reason: %s", err)
	}
	log.Println("Ready to read.")
	running := true
	for running {
		select {
		case <-stopCtx.Done():
			log.Println("Stopping the reader ...")
			err := file.Close()
			if err != nil {
				log.Printf("Failed to close the file. Reason: %s", err)
			}
			running = false
			break
		default:
			_, err := file.Read(block)
			if err != nil && err != io.EOF {
				log.Fatalln("Failed to read from file. Reason:", err)
			}
			if err == io.EOF {
				fmt.Print(".")
				time.Sleep(1 * time.Second)
				continue
			}
			d, err := data.Decode(block)
			if err != nil {
				log.Fatalln("Failed to decode data", err)
			}
			dataCh <- d
		}
	}
	log.Println("Reader has stopped.")
	stopWg.Done()
}

func consumer(dataCh chan *data.SomeData, stopCtx context.Context, stopWg *sync.WaitGroup) {
	log.Println("Starting to consume ...")
	running := true
	for running {
		select {
		case <-stopCtx.Done():
			log.Println("Stopping the consumer ...")
			running = false
			break
		case data := <-dataCh:
			log.Printf("Got %+v\n", *data)
		default:
			time.Sleep(1 * time.Second)
		}
	}
	log.Println("Consumer has stopped.")
	stopWg.Done()
}

func waitingForGracefulShutdown(cancelFn context.CancelFunc, stopWg *sync.WaitGroup) {
	osStopChan := make(chan os.Signal, 1)
	signal.Notify(osStopChan, syscall.SIGINT, syscall.SIGTERM)
	<-osStopChan
	log.Println("Shutting down ...")
	cancelFn()
	stopWg.Wait()
}
