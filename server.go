package main

import (
	"fmt"
	"log"
	"net/http"
)

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Cats Service. Version 0.1\n")
}

func server() {
	log.Println("Server started")
	http.HandleFunc("/ping", ping)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
