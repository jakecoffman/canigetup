package main

import (
	"encoding/json"
	"github.com/jakecoffman/canigetup/pi"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	rpi := pi.New()
	defer rpi.Close()

	state := NewState(rpi)
	lock := sync.RWMutex{}

	go func() {
		before, err := timeToMinutes(time.Now().Format("15:04"))
		if err != nil {
			log.Fatal(err)
		}
		for {
			time.Sleep(1 * time.Second)
			now, err := timeToMinutes(time.Now().Format("15:04"))
			if err != nil {
				log.Fatal(err)
			}
			if now == before {
				continue
			}
			lock.Lock()
			for _, schedule := range state.Schedules {
				t, err := timeToMinutes(schedule.At)
				if err != nil {
					log.Fatal(err)
				}
				if before < t && now == t {
					state.On = schedule.On
				}
			}
			lock.Unlock()
			before = now
		}
	}()

	r := http.NewServeMux()
	r.HandleFunc("/api/state", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			lock.RLock()
			defer lock.RUnlock()

			jsonMessage(state, writer)
		case http.MethodPut:
			var payload State
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				log.Println("Failed to decode PUT", err)
				writer.WriteHeader(400)
				jsonMessage("Failed to decode", writer)
				return
			}

			for _, schedule := range payload.Schedules {
				_, err := timeToMinutes(schedule.At)
				if err != nil {
					writer.WriteHeader(400)
					jsonMessage("Invalid time", writer)
					return
				}
			}

			lock.Lock()
			defer lock.Unlock()

			state.Schedules = payload.Schedules
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
	On bool
	Schedules []Schedule
	pi *pi.Pi
}

type Schedule struct {
	At string // e.g. "15:00"
	On bool
}

func NewState(rpi *pi.Pi) *State {
	s := &State{pi: rpi}
	s.load()
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
		s.pi.High()
	} else {
		s.pi.Low()
	}
}

func (s *State) Turn(on bool) {
	s.On = on
	if s.On {
		s.pi.High()
	} else {
		s.pi.Low()
	}
	s.save()
}

func timeToMinutes(time string) (int, error) {
	splitTime := strings.Split(time, ":")
	hours, err := strconv.Atoi(splitTime[0])
	if err != nil {
		log.Println(err)
		return 0, err
	}
	minutes, err := strconv.Atoi(splitTime[1])
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return hours * 60 + minutes, nil
}
