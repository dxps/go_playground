package main

import (
	"log"
	"net/http"
	"os"
)

// DefaultPort is the default port used for listening for requests.
// Use SERVER_PORT environment variable if you need another value.
const DefaultPort = "8000"

func getServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port != "" {
		return port
	}

	return DefaultPort
}

// MirrorHandler returns back the request data.
func MirrorHandler(writer http.ResponseWriter, request *http.Request) {

	log.Println(">>> Mirroring the request '" + request.URL.Path + "'")

	writer.Header().Set("Access-Control-Allow-Origin", "*")

	// pre-flight headers are permitted
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")

	request.Write(writer)
}

func main() {

	log.Println(">>> Starting to listen on port " + getServerPort())

	http.HandleFunc("/", MirrorHandler)
	http.ListenAndServe(":"+getServerPort(), nil)
}
