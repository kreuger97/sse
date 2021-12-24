package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type (
	Client struct {
		addres string
	}
)

var listeners map[string]*Client
var counter int

func main() {
	listeners = make(map[string]*Client)
	go loop()
	go func() {
		for {
			fmt.Println("listeners:")
			for key := range listeners {
				fmt.Printf("\tlistener: %s\n", key)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("")))
	http.HandleFunc("/sse", handler)
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		panic("could not start the server")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	connected := true

	fmt.Printf("%s connected\n", r.RemoteAddr)
	if listeners[r.RemoteAddr] == nil {
		listeners[r.RemoteAddr] = &Client{r.RemoteAddr}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for ;connected; {
		messageID := uint(rand.Uint32())
		fmt.Fprintf(w, "data: %d\nid:%d\n\n", counter, messageID)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		
		select {
		case <-r.Context().Done():
			connected = false
			delete(listeners, r.RemoteAddr)
			fmt.Printf("%s close connection\n", r.RemoteAddr)
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func loop() {
	for counter = 0; ; counter++ {
		time.Sleep(10 * time.Second)
	}
}
