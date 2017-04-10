package main

import (
	"fmt"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func cdcHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Testing1")
}

func main() {
	http.HandleFunc("/fastor", cdcHandler)
	http.ListenAndServe(":"+os.Args[1], nil)
}
