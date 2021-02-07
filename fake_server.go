package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func startServer() {
	s := http.NewServeMux()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count := r.URL.Query().Get("count")
		time.Sleep(time.Duration(time.Duration(rand.Intn(5)) * time.Second))
		w.Write([]byte(fmt.Sprintf("Finished request %s", count)))
	})

	go http.ListenAndServe("localhost:3000", s)
	time.Sleep(100 * time.Millisecond)
}
