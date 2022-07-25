package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {

	s := &Server{
		ch: make(chan []byte, 2048),
	}
	srv := http.Server{
		Handler:      s.routes(),
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// Server is the main Http Server of the app.
type Server struct {
	m        sync.Map
	lock     sync.Mutex
	upgrader websocket.Upgrader
	ch       chan []byte
}

func (s *Server) routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/crud/{key}", s.HandleCreate).Methods(http.MethodPost)
	router.HandleFunc("/crud/{key}", s.HandleRetrieve).Methods(http.MethodGet)
	router.HandleFunc("/crud/{key}", s.HandleUpdate).Methods(http.MethodPut)
	router.HandleFunc("/crud/{key}", s.HandleDelete).Methods(http.MethodDelete)
	router.HandleFunc("/ws", s.HandleWS)
	return router
}

func (s *Server) merge(dst interface{}, src map[string]interface{}) map[string]interface{} {
	switch dst := dst.(type) {
	case map[string]interface{}:
		// range over keys
		for key, value := range src {
			// populate if not exists
			if _, ok := dst[key]; !ok {
				dst[key] = value
				continue
			}

			// merge if exists
			switch t := value.(type) {
			case map[string]interface{}:
				dst[key] = s.merge(dst[key], t)
			default:
				dst[key] = value
			}
		}
		return dst
	default:
		return src
	}
}

// HandleWS is the WebSockets handler.
func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("HandleWS > Upgrade error:", err)
		return
	}
	defer c.Close()
	for update := range s.ch {
		err = c.WriteMessage(websocket.TextMessage, update)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// HandleCreate is the handler of create request.
func (s *Server) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// get key
	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// read body
	body := r.Body
	defer body.Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println("unable to read body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate it as json
	throwAway := struct{}{}
	err = json.Unmarshal(data, &throwAway)
	if err != nil {
		log.Printf("Unable to read POST body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// store data
	_, exists := s.m.LoadOrStore(key, data)
	if exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// send to subscriber
	s.ch <- []byte(fmt.Sprintf(`{"key":%q,"method":"created"}`, key))

	// write created header
	w.WriteHeader(http.StatusCreated)
}

// HandleRetrieve is the handler of retrieving an item request.
func (s *Server) HandleRetrieve(w http.ResponseWriter, r *http.Request) {
	// read key
	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// load key from store
	v, ok := s.m.Load(key)
	if !ok || v == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// convert to []byte
	data, ok := v.([]byte)
	if !ok || v == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write value
	w.Header().Add("content-type", "application/json")
	if _, err := w.Write(data); err != nil {
		log.Printf("unable to write response: %s\n", err.Error())
	}
}

// HandleUpdate is the handler of updating an item request.
func (s *Server) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// get key
	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// read body
	body := r.Body
	defer body.Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println("unable to read body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate it as json
	update := make(map[string]interface{})
	err = json.Unmarshal(data, &update)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// prevent update races
	s.lock.Lock()
	defer s.lock.Unlock()

	// get existing value
	v, ok := s.m.Load(key)
	if !ok || v == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// convert to []byte
	existingData, ok := v.([]byte)
	if !ok || v == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// convert to json
	existing := make(map[string]interface{})
	err = json.Unmarshal(existingData, &existing)
	if err != nil {
		log.Println("Unable to unmarshal from storage", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// merge the two
	result := s.merge(existing, update)

	// marshal json
	output, err := json.Marshal(result)
	if err != nil {
		log.Println("Unable to marshal json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// store data
	s.m.Store(key, output)

	// send to subscriber
	s.ch <- []byte(fmt.Sprintf(`{"key":%q,"method":"updated"}`, key))
}

// HandleDelete is the handler of deleting an item request.
func (s *Server) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// load key
	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// delete if exists
	s.m.Delete(key)

	// send to subscriber
	s.ch <- []byte(fmt.Sprintf(`{"key":%q,"method":"deleted"}`, key))
}
