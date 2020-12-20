package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/devisions/go-playground/go-directio/data"
	"github.com/ncw/directio"
	"github.com/pkg/errors"
)

var block = directio.AlignedBlock(data.BlockSize)

func main() {
	// Preparing the graceful shutdown elements.
	stopCtx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	stopWg := &sync.WaitGroup{}
	stopWg.Add(2)

	dataCh := make(chan data.SomeData, 1_000_000)
	defer close(dataCh)

	go writer(dataCh, stopCtx, stopWg)
	go producer(dataCh, stopCtx, stopWg)

	waitingForGracefulShutdown(cancelFn, stopWg)
}

func writer(dataCh chan data.SomeData, stopCtx context.Context, stopWg *sync.WaitGroup) {
	file := "/Users/devisions/tmp/test_dio"

	out, err := directio.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open file for writing. Reason: %s", err)
	}
	log.Println("Ready to write.")
	running := true
	for running {
		select {
		case <-stopCtx.Done():
			log.Println("Stopping the writer ...")
			l := len(dataCh)
			if l > 0 {
				log.Printf("Writing to file the remaining %d data items ...", l)
				for len(dataCh) > 0 {
					d := <-dataCh
					if err := write(out, &d); err != nil {
						log.Println("Failed writing to file. Reason:", err)
						break
					}
					log.Println("Wrote", d)
				}
			}
			err := out.Close()
			if err != nil {
				log.Printf("Failed closing the file. Reason: %s", err)
			}
			running = false
			break
		case data := <-dataCh:
			if err := write(out, &data); err != nil {
				log.Println("Failed writing to file. Reason:", err)
				running = false
				break
			}
		default:
			fmt.Print(".")
			time.Sleep(1 * time.Second)
		}
	}
	log.Println("Writer stopped.")
	stopWg.Done()
}

func producer(dataCh chan data.SomeData, stopCtx context.Context, stopWg *sync.WaitGroup) {
	var i uint32
	log.Println("Starting to produce ...")
	running := true
	for running {
		select {
		case <-stopCtx.Done():
			log.Println("Stopping the producer ...")
			running = false
			break
		default:
			i++
			d := data.SomeData{Value: i}
			dataCh <- d
			log.Printf("Produced %+v\n", d)
			time.Sleep(1 * time.Second)
		}
	}
	log.Println("Producer stopped.")
	stopWg.Done()
}

func write(out *os.File, d *data.SomeData) error {
	err := d.Encode(block)
	if err != nil {
		return errors.Wrap(err, "encoding data")
	}
	_, err = out.Write(block)
	if err != nil {
		return errors.Wrap(err, "writing to file")
	}
	return nil
}

func waitingForGracefulShutdown(cancelFn context.CancelFunc, stopWg *sync.WaitGroup) {
	osStopChan := make(chan os.Signal, 1)
	signal.Notify(osStopChan, syscall.SIGINT, syscall.SIGTERM)
	<-osStopChan
	log.Println("Shutting down ...")
	cancelFn()
	stopWg.Wait()
}
