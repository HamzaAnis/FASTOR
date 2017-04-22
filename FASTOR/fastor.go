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
	"strconv"
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

func enterDetails(a net.Conn) {
	//Read message from the server
	message := make([]byte, 100)
	a.Read(message)
	fmt.Println(string(message))
	//Welcome message
	a.Read(message)
	fmt.Println(string(message))

	reader := bufio.NewReader(os.Stdin)
	//relay name
	name, _ := reader.ReadString('\n')
	a.Write([]byte(name))

	//Participate message
	a.Read(message)
	fmt.Println(string(message))

	//particiption choice
	part, _ := reader.ReadString('\n')
	// part := "Yes"
	a.Write([]byte(part))
}

func main() {
	//1st port wil be for the webserver and second port will be for the torserver
	port := ""
	relaysserverport := ""
	relaysCountPort := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
		relaysserverport = os.Args[2]
		temp, _ := strconv.Atoi(relaysserverport)
		temp++
		relaysCountPort = string(temp)
	} else {
		port = "8081"
		relaysserverport = "9696"
		relaysCountPort = "9697"
	}

	a, err := net.Dial("tcp", "localhost:"+relaysserverport)
	numRelays, err := net.Dial("tcp", "localhost:"+relaysCountPort)
	if err != nil {
		fmt.Println("Dial error:", err)
	}
	defer a.Close()
	defer numRelays.Close()
	go enterDetails(a)

	for {
		number := make([]byte, 20)
		numRelays.Write([]byte("Number of relays"))
		numRelays.Read(number)
		// fmt.Println("The number of the relays online is ", string(number))
		if string(number) == "1" {
			fmt.Println("The relays are available on the server. Starting HTTp Server")
			break
		}
	}
	// default will go here
	http.HandleFunc("/fastor/", torhandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("."))))
	http.ListenAndServe(":"+port, nil)
}
