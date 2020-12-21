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

	"github.com/devisions/go-playground/go-directio/config"
	"github.com/devisions/go-playground/go-directio/consumer/internal"
	"github.com/devisions/go-playground/go-directio/internal/data"
	"github.com/ncw/directio"
	"github.com/pkg/errors"
)

var block []byte

// Global state of the current file to read from.
var in *os.File

// Global state of the consumer.
var state *internal.ConsumerState

func main() {
	// Preparing the graceful shutdown elements.
	stopCtx, cancelFn := context.WithCancel(context.Background())
	stopWg := &sync.WaitGroup{}
	stopWg.Add(2)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}
	log.Printf("Using blocksize %d and reading files in path %s\n", cfg.BlockSize, cfg.Path)
	// directio.AlignSize = cfg.BlockSize
	block = directio.AlignedBlock(cfg.BlockSize)

	dataCh := make(chan *data.SomeData, 1_000_000)

	state, err = internal.InitConsumerState(cfg.Path, cfg.BlockSize)
	if err != nil {
		log.Fatalln("Failed to init state. Reason:", err)
	}
	log.Printf("Initial state: ReadFilepath:%s Readblocks:%d \n", state.ReadFilepath, state.ReadBlocks)

	go consumer(dataCh, stopCtx, stopWg)
	go reader(cfg.Path, cfg.FileMaxSize, dataCh, stopCtx, stopWg)

	waitingForGracefulShutdown(cancelFn, stopWg)
}

func reader(filepathPrefix string, fileMaxsize int64, dataCh chan *data.SomeData, stopCtx context.Context, stopWg *sync.WaitGroup) {
	running := true
	for running {
		select {
		case <-stopCtx.Done():
			log.Println("Stopping the reader ...")
			err := in.Close()
			if err != nil {
				log.Printf("Failed to close the file. Reason: %s", err)
			}
			running = false
			break
		default:
			f, err := data.GetFileForReading(in, filepathPrefix, state.ReadFilepath, fileMaxsize)
			if err != nil {
				if !os.IsNotExist(errors.Cause(err)) {
					log.Fatalln("Failed to get the file to read. Reason:", err)
				}
				// There is no file to read from. Let's wait ...
				fmt.Printf(".")
				time.Sleep(1 * time.Second)
				state.ReadBlocks = 0
				continue
			}
			if in == nil {
				log.Println("Reading from file", f.Name(), "and skipping", state.ReadBlocks, "blocks")
				if state.ReadBlocks > 0 {
					if _, err := f.Seek(int64(state.ReadBlocks*state.SaveBlocksize), 0); err != nil {
						log.Fatalln("Failed to skip already read blocks. Reason", err)
					}
				}
				state.ReadFilepath = f.Name()
				in = f
			}
			if f.Name() != in.Name() {
				log.Println("Reading from file", f.Name())
				state.ReadFilepath = f.Name()
				state.ReadBlocks = 0
				if err = data.DeleteFile(in.Name()); err != nil {
					log.Fatalln("Failed to delete completely read file. Reason", err)
				}
				log.Println("Deleted completely read file", in.Name())
				in = f
			}
			_, err = in.Read(block)
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
	running := true
	for running {
		select {
		case <-stopCtx.Done():
			log.Println("Stopping the consumer ...")
			running = false
			break
		case data := <-dataCh:
			log.Printf("Consumed %+v\n", *data)
			state.ReadBlocks += 1
			err := state.SaveToFile()
			if err != nil {
				log.Fatalln("Failed to save state to file. Reason:", err)
			}
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
