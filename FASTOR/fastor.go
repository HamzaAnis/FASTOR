package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//A function to check the request and do the corrosponding action
//E.g For image save the image to the file
//If stylesheet save the stylesheet
func handleStyle(link string) {
	u, err := url.Parse(link)
	if err != nil {
		panic(err)
	}
	path := u.Path[1:] //to avoid the / at the beginning
	filetype := strings.Split(path, ".")
	fmt.Printf("Type is %v\n", filetype[len(filetype)-1])
	if filetype[len(filetype)-1] != "html" {
		fmt.Println("The image is found named " + path)
	}
}

//on every / request the handler will call this
func torhandler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Path[1:7]
	link := "http://" + r.URL.Path[8:]

	// handleStyle(link)

	// if domain == "fastor" {
	fmt.Printf("Domain: %v\n", domain)
	fmt.Printf("Website:  %v\n", link)
	res, err := http.Get(link)
	if err != nil {
		// log.Fatal(err)
		http.Error(w, err.Error(), 500)
	}
	buf, err := ioutil.ReadAll(res.Body)
	fmt.Fprintf(w, string(buf))
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func ReadLinesInto(a net.Conn) {
	bufc := bufio.NewReader(a)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Println("The message get is " + line)
	}
}

func main() {
	//1st port wil be for the webserver and second port will be for the torserver
	port := ""
	relaysserverport := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
		relaysserverport = os.Args[2]
	} else {
		port = "9999"
		relaysserverport = "9696"
	}

	a, err := net.Dial("tcp", "localhost:"+relaysserverport)
	if err != nil {
		fmt.Println("Dial error:", err)
	}
	defer a.Close()
	go ReadLinesInto(a)
	fmt.Println("Usma anis")
	//default will go here
	http.HandleFunc("/fastor/", torhandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("."))))
	http.ListenAndServe(":"+port, nil)
}
