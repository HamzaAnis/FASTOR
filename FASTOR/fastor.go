package main

import (
	"fmt"
	"net/http"
	"os"
)

//on every / request the handler will call this
func torhandler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Path[1:7]
	link := r.URL.Path[11:]
	// fmt.Println(r.URL.Path[1:])
	// fmt.Println(r.URL.Path[11:])
	if domain == "fastor" {
		fmt.Fprintf(w, "Domain: %v     ", domain)
		fmt.Fprintf(w, "Website:  %v   ", link)
	} else {
		fmt.Fprintf(w, "There is no services available for which you requested")
		// html, er := ioutil.ReadFile("test.html")
		// if er == nil {
		// 	fmt.Println("The file is  read")
		// } else {
		// 	fmt.Println("The file is not read")
		// }
		// fmt.Fprint(w, string(html))
	}

}

func main() {
	port := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "9999"
	}
	//default will go here
	http.HandleFunc("/", torhandler)
	http.ListenAndServe(":"+port, nil)
}
