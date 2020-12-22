package main

import (
	"context"
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

	dataCh := make(chan *internal.ConsumerData, 1_000_000)

	state, err = internal.InitConsumerState(cfg.Path, cfg.BlockSize)
	if err != nil {
		log.Fatalln("Failed to init state. Reason:", err)
	}
	if !state.IsEmpty() {
		log.Printf("Starting with state: ReadFilepath:%s Readblocks:%d \n", state.ReadFilepath, state.ReadBlocks)
	} else {
		log.Printf("Starting with an empty state")
	}

	fileMaxSize := cfg.MaxBlocks * int64(cfg.BlockSize)
	go consumer(fileMaxSize, dataCh, stopCtx, stopWg)
	go reader(cfg.Path, fileMaxSize, dataCh, stopCtx, stopWg)

	waitingForGracefulShutdown(cancelFn, stopWg)
}

func reader(filepathPrefix string, fileMaxsize int64, dataCh chan *internal.ConsumerData, stopCtx context.Context, stopWg *sync.WaitGroup) {
	if !state.IsEmpty() {
		var f *os.File
		var err error
		f, err = data.OpenFileForReading(state.ReadFilepath)
		if err != nil {
			if os.IsNotExist(errors.Cause(err)) {
				log.Println("Last read file is missing. Let's look for any first file to read...")
			} else {
				log.Fatalln("Failed to open last read file (according to the state). Reason:", err)
			}
		}
		for f == nil {
			f, err = internal.GetFileForReading(nil, filepathPrefix, fileMaxsize)
			if err != nil {
				if !os.IsNotExist(errors.Cause(err)) {
					log.Fatalln("Failed to get a file to read. Reason:", err)
				}
				// No file exists. Let's wait for anything new.
				time.Sleep(1 * time.Second)
				continue
			}
		}
		in = f
		log.Println("Reading from file", f.Name())
	}

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
			f, err := internal.GetFileForReading(in, filepathPrefix, fileMaxsize)
			if err != nil {
				if !os.IsNotExist(errors.Cause(err)) {
					log.Fatalln("Failed to get a file to read. Reason:", err)
				}
				// There is no file to read from. Let's wait ...
				time.Sleep(1 * time.Second)
				continue
			}
			if in == nil {
				if state.ReadBlocks > 0 {
					log.Println("Reading from file", f.Name(), "and skipping", state.ReadBlocks, "blocks")
					if _, err := f.Seek(state.SeekOffset(), 0); err != nil {
						log.Fatalln("Failed to skip already read blocks. Reason", err)
					}
				} else {
					log.Println("Reading from file", f.Name())
				}
				in = f
			} else {
				// New file for reading provided, let's close the existing and start using it.
				if f != nil && f.Name() != in.Name() {
					if err := in.Close(); err != nil {
						log.Printf("[WARN] Failed to close existing file '%s'. Reason:%s\n", in.Name(), err)
					}
					log.Println("Reading from new file", f.Name())
					in = f
				}
			}
			_, err = in.Read(block)
			if err != nil && err != io.EOF {
				log.Fatalln("Failed to read from file. Reason:", err)
			}
			if err == io.EOF {
				time.Sleep(1 * time.Second)
				continue
			}
			d, err := data.Decode(block)
			if err != nil {
				log.Fatalln("Failed to decode data", err)
			}
			dataCh <- &internal.ConsumerData{Data: d, FromFilepath: in.Name()}
		}
	}
	log.Println("Reader has stopped.")
	stopWg.Done()
}

func consumer(fileMaxsize int64, dataCh chan *internal.ConsumerData, stopCtx context.Context, stopWg *sync.WaitGroup) {
	running := true
	for running {
		select {

		case <-stopCtx.Done():
			log.Println("Stopping the consumer ...")
			running = false
			break

		case cd := <-dataCh:
			log.Printf("Consumed %+v\n", *cd.Data)
			tryDelete(state.ReadFilepath, fileMaxsize)
			if cd.FromFilepath != state.ReadFilepath {
				state.UseNew(cd.FromFilepath)
			} else {
				state.ReadBlocks += 1
			}
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

func tryDelete(filepath string, maxSize int64) bool {
	deleted, err := data.DeleteFileIfReachedMaxSize(filepath, maxSize)
	if err != nil {
		log.Println("[WARN] Failed while trying to check and delete the consumed file", filepath, "Reason:", err)
	}
	if deleted {
		log.Println("Deleted the consumed file", state.ReadFilepath)
		return true
	}
	return false
}

func waitingForGracefulShutdown(cancelFn context.CancelFunc, stopWg *sync.WaitGroup) {
	osStopChan := make(chan os.Signal, 1)
	signal.Notify(osStopChan, syscall.SIGINT, syscall.SIGTERM)
	<-osStopChan
	log.Println("Shutting down ...")
	cancelFn()
	stopWg.Wait()
}
