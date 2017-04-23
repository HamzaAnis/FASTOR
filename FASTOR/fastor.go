package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"io"

	"github.com/fatih/color"
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
func torhandler(w http.ResponseWriter, r *http.Request, server net.Conn) {
	domain := r.URL.Path[1:7]
	link := "http://" + r.URL.Path[8:]

	// handleStyle(link)

	// if domain == "fastor" {
	fmt.Printf("Domain: %v\n", domain)
	fmt.Printf("Website:  %v\n", link)
	// server.Write([]byte(link))
	io.WriteString(server, link)
	_, err := server.Write([]byte(link))

	content := make([]byte, 1000000)

	n, err := server.Read(content)
	color.Blue(string(content[:n]))

	fmt.Fprintf(w, string(content[:n]))
	// defer res.Body.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func enterDetails(a net.Conn) {
	c := color.New(color.FgHiCyan)
	c.Add(color.Bold)

	//Read message from the server
	message := make([]byte, 100)
	a.Read(message)
	c.Println("\t\t\t", string(message))
	//Welcome message
	a.Read(message)
	c.Println(string(message))

	reader := bufio.NewReader(os.Stdin)
	//relay name
	name, _ := reader.ReadString('\n')
	a.Write([]byte(name))

	//Participate message
	a.Read(message)
	c.Println(string(message))

	//particiption choice
	part, _ := reader.ReadString('\n')
	// part := "Yes"
	a.Write([]byte(part))
}

func main() {
	// c := color.New(color.FgBlue)
	//1st port wil be for the webserver and second port will be for the torserver
	port := ""
	relaysserverport := ""
	relaysCountPort := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
		relaysserverport = os.Args[2]
		relaysCountPort = os.Args[3]
		// fmt.Printf("the sencond one is %vd", relaysserverport)
		// temp, _ := strconv.Atoi(relaysserverport)
		// temp++
		// relaysCountPort = string(temp)
	} else {
		port = "9825"
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
		number := make([]byte, 1)
		numRelays.Write([]byte("Number of relays"))
		numRelays.Read(number)
		// color.Blue("The number ofH the relays online is %v\n", string(number))
		if string(number) == "1" {
			color.Red("The minimum relays are available on the server. Starting HTTP Server")
			break
		}
	}

	// default will go here
	http.HandleFunc("/fastor/", func(w http.ResponseWriter, r *http.Request) {
		torhandler(w, r, a)
	})
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("."))))
	http.ListenAndServe(":"+port, nil)
}
