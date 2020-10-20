package main

import (
	"log"
	"mime"
	"mime/multipart"
	"net/http"
)

type myResponse struct {
	Values []string `json:"values"`
}

func main() {

	srv := http.Server{
		Addr:    "localhost:8001",
		Handler: http.HandlerFunc(respondMultipart),
	}
	log.Fatal(srv.ListenAndServe())
}

func respondMultipart(w http.ResponseWriter, r *http.Request) {

	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Accept"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	if mediatype != "multipart/form-data" {
		http.Error(w, "set Accept: multipart/form-data", http.StatusMultipleChoices)
		return
	}

	resp := myResponse{
		Values: []string{
			"abc", "def", "ghi",
		},
	}

	mw := multipart.NewWriter(w)
	w.Header().Set("Content-Type", mw.FormDataContentType())
	for _, value := range resp.Values {
		fw, err := mw.CreateFormField("value")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fw.Write([]byte(value)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := mw.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
