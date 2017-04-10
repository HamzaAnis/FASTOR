package main

import (
	"fmt"
	"net/http"
	"os"
)

func cdcHandler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Path[1:7]
	link := r.URL.Path[11:]
	fmt.Println(r.URL.Path[1:])

	fmt.Fprintf(w, "Domain: %v     ", domain)
	fmt.Fprintf(w, "Website:  %v   ", link)
}

func main() {
	port := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "9999"
	}
	//default will go here
	http.HandleFunc("/", cdcHandler)
	http.ListenAndServe(":"+port, nil)
}
