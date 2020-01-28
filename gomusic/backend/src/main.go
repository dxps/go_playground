package main

import (
	"log"

	"devisions.org/gomusic-be/src/rest"
)

func main() {
	log.Println("Main log....")
	rest.RunAPI(":9090")
}
