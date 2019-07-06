package main

import (
	"encoding/json"
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := rpio.Close(); err != nil {
			log.Println(err)
		}
	}()

	pin := rpio.Pin(15)
	pin.Output()

	state := NewState(pin)

	r := http.NewServeMux()

	r.HandleFunc("/api/state", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			jsonMessage(state, writer)
		case http.MethodPut:
			var payload State
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				log.Println("Failed to decode PUT", err)
				writer.WriteHeader(400)
				jsonMessage("Failed to decode", writer)
				return
			}
			state.Schedules = payload.Schedules
			// TODO handle schedule changes
			state.Turn(payload.On)
			jsonMessage(state, writer)
		}
	})

	fs := http.FileServer(http.Dir("/home/pi/web/"))
	r.Handle("/", fs)

	log.Println("serving on http://127.0.0.1:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
}

func jsonMessage(msg interface{}, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(msg); err != nil {
		log.Println(err)
	}
}

type State struct {
	sync.RWMutex
	On bool
	Schedules []Schedule
	pin rpio.Pin
}

type Schedule struct {
	At string // e.g. "15:00"
	On bool
}

func NewState(pin rpio.Pin) *State {
	s := &State{pin: pin}
	s.Lock()
	s.load()
	s.Unlock()
	return s
}

const fileName = "/home/pi/canigetup.json"

func (s *State) save() {
	// assumes we are already locked
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Failed to create file", err)
		return
	}
	if err = json.NewEncoder(file).Encode(s); err != nil {
		log.Println("Failed to serialize state", err)
	}
}

func (s *State) load() {
	// assumes we are already locked
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("failed to load state", err)
		return
	}
	if err = json.NewDecoder(file).Decode(s); err != nil {
		log.Println("failed to deserialize state", err)
	}
	if s.On	{
		s.pin.High()
	} else {
		s.pin.Low()
	}
	// TODO handle schedule
}

func (s *State) Turn(on bool) {
	s.Lock()
	defer s.Unlock()
	s.On = on
	if s.On {
		s.pin.High()
	} else {
		s.pin.Low()
	}
	s.save()
}
