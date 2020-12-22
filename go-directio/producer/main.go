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

	"github.com/devisions/go-playground/go-directio/config"
	"github.com/devisions/go-playground/go-directio/internal/data"
	"github.com/devisions/go-playground/go-directio/producer/internal"
	"github.com/ncw/directio"
	"github.com/pkg/errors"
)

var block []byte

// Global state of the current file to write into.
var out *os.File

func main() {
	// Preparing the graceful shutdown elements.
	stopCtx, cancelFn := context.WithCancel(context.Background())
	stopWg := &sync.WaitGroup{}
	stopWg.Add(2)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}
	log.Printf("Using blocksize %d and writing files in path %s\n", cfg.BlockSize, cfg.Path)
	// directio.AlignSize = cfg.BlockSize
	block = directio.AlignedBlock(cfg.BlockSize)

	dataCh := make(chan data.SomeData, 1_000_000)
	defer close(dataCh)

	fileMaxSize := cfg.MaxBlocks * int64(cfg.BlockSize)
	go writer(cfg.Path, fileMaxSize, dataCh, stopCtx, stopWg)
	go producer(dataCh, stopCtx, stopWg)

	waitingForGracefulShutdown(cancelFn, stopWg)
}

func writer(filepathPrefix string, fileMaxsize int64, dataCh chan data.SomeData, stopCtx context.Context, stopWg *sync.WaitGroup) {

	f, err := internal.GetFileForWriting(nil, filepathPrefix, fileMaxsize)
	if err != nil {
		log.Fatalln("Failed to look for the next file to write into. Reason:", err)
	}
	if f == nil { // This should never happen; used just for safety.
		log.Fatalln("No file to write could be used.")
	}
	out = f
	log.Println("Ready to write on file", out.Name())

	running := true
	for running {
		select {
		case <-stopCtx.Done():
			log.Println("Stopping the writer ...")
			l := len(dataCh)
			if l > 0 {
				log.Printf("Draining the channel: writing to file the remaining %d data items ...", l)
				for len(dataCh) > 0 {
					d := <-dataCh
					if err := write(filepathPrefix, fileMaxsize, &d); err != nil {
						log.Println("Failed writing to file. Reason:", err)
						break
					}
					log.Println("Wrote", d)
				}
			}
			if out != nil { // Just for safety reasons.
				err := out.Close()
				if err != nil {
					log.Printf("Failed closing the file. Reason: %s", err)
				}
			}

			running = false
			break
		case d := <-dataCh:
			if err := write(filepathPrefix, fileMaxsize, &d); err != nil {
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

func write(filepathPrefix string, fileMaxsize int64, d *data.SomeData) error {
	err := d.Encode(block)
	if err != nil {
		return errors.Wrap(err, "encoding data")
	}
	f, err := internal.GetFileForWriting(out, filepathPrefix, fileMaxsize)
	if err != nil {
		return err
	}
	// A new file has been provided, so close existing and start using it.
	if f != nil {
		if err := out.Close(); err != nil {
			log.Printf("[WARN] Failed to close existing file '%s'. Reason:%s\n", out.Name(), err)
		}
		log.Println("Writing to new file", f.Name())
		out = f
	}
	_, err = out.Write(block)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("writing to file (the item %+v)", *d))
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
