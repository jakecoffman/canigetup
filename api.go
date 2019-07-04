package main

import (
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"net/http"
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

	r := http.NewServeMux()

	r.HandleFunc("/api/on", func(writer http.ResponseWriter, request *http.Request) {
		pin.High()
		jsonMessage(`{"msg":"on"}`, writer)
	})

	r.HandleFunc("/api/off", func(writer http.ResponseWriter, request *http.Request) {
		pin.Low()
		jsonMessage(`{"msg":"off"}`, writer)
	})

	r.HandleFunc("/api/toggle", func(writer http.ResponseWriter, request *http.Request) {
		pin.Toggle()
		jsonMessage(`{"msg":"toggle"}`, writer)
	})

	fs := http.FileServer(http.Dir("/home/pi/web/"))
	r.Handle("/", fs)

	log.Println("serving on http://127.0.0.1:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
}

func jsonMessage(msg string, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")
	if _, err := writer.Write([]byte(msg)); err != nil {
		log.Println(err)
	}
}
