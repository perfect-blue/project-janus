package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := http.Server{
		Addr: "localhost:8000",
		Handler: setupEndpoints(),
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func setupEndpoints() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	}))
	mux.HandleFunc("/health", healthCheck)
	return mux
}