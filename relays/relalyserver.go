package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apcera/termtables"
	"github.com/fatih/color"
)

//Relay is a struct to store the information of the relays
type Relay struct {
	conn        net.Conn
	name        string
	number      int
	ch          chan string
	participate bool
}

//Request have the information about the request that goes to directory server
type Request struct {
	url   string
	relay *Relay
}

func main() {
	//relay directory server
	relaysDatabse := make(map[int]Relay)
	numberOfRelays := 0
	totalRelays := 0
	port := ""
	relaynNumberPort := ""
	hearbeatport := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
		temp, _ := strconv.Atoi(port)
		temp++
		relaynNumberPort = string(temp)
		temp++
		hearbeatport = string(temp)
	} else {
		port = "9696"
		relaynNumberPort = "9697"
		hearbeatport = "9698"
	}
	connection, err := net.Listen("tcp", ":"+port)
	connection2, err2 := net.Listen("tcp", ":"+relaynNumberPort)
	connection3, err3 := net.Listen("tcp", ":"+hearbeatport)
	if err != nil || err2 != nil || err3 != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	requestchan := make(chan Request)

	addRelay := make(chan Relay)

	rmRelay := make(chan Relay)

	go handleRelays(relaysDatabse, requestchan, addRelay, rmRelay, &totalRelays)
	for {
		conn, err := connection.Accept()
		conn2, err2 := connection2.Accept()
		conn3, err3 := connection3.Accept()
		if err != nil || err2 != nil || err3 != nil {
			fmt.Println(err)
			continue
		} else {
			color.Red("\tA client has connected")
			conn.Write([]byte("Hello FASTOR user!"))
			go handleConnection(conn, conn2, conn3, requestchan, addRelay, rmRelay, &numberOfRelays, &totalRelays)
		}
	}
}

func promptName(c net.Conn) string {
	clr := color.New(color.FgGreen)
	io.WriteString(c, "What is your relay name? ")
	name := make([]byte, 20)
	n, _ := c.Read(name)
	// clr.Println("The length is ", n)
	clr.Printf("The name of the relay: %v", string(name[:n]))
	return string(name[:n-2])
}

func promptChoice(c net.Conn) bool {
	c.Write([]byte("Do you want to participate in the anonymous service?(Y/N)"))
	choice := make([]byte, 1)
	c.Read(choice)
	// fmt.Printf("The length of choice is %v", n)
	participate := false
	if string(choice) == "Y" {
		participate = true
		color.Green("Relay is participating")
	} else {
		participate = false
		color.Red("Relay is not participating")
	}
	return participate
}

//Core
func handleConnection(c net.Conn, num net.Conn, heatrbeat net.Conn, requestchan chan<- Request, addRelay chan<- Relay, rmRelay chan<- Relay, numberRelay *int, totalRelays *int) {
	//we first need to add current relay to the channel
	//filling in the relay structure
	relay := Relay{
		conn:        c,
		name:        promptName(c),
		ch:          make(chan string),
		participate: promptChoice(c),
		number:      *numberRelay,
	}
	*numberRelay++
	if strings.TrimSpace(relay.name) == "" {
		io.WriteString(c, "Invalid relay name\n")
		return
	}

	//it populates the map
	addRelay <- relay

	//ignore for the time being
	defer func() {
		// color.Yellow("Exiting 3")
		color.Yellow("Connection from %v closed.\n", c.RemoteAddr())
		rmRelay <- relay
	}()

	go relay.ServerRequest(requestchan)

	//to send the nnumber of relays
	go sendNumber(num, totalRelays)

	go relay.heartB(heatrbeat, rmRelay)

	relay.WriteLinesFrom(relay.ch)

}

func (c Relay) heartB(hB net.Conn, rmRelay chan<- Relay) {
	for {
		time.Sleep(1 * time.Second)

		_, err := hB.Write([]byte("Server is Up"))
		if err != nil {
			rmRelay <- c
			break
		}
		res := make([]byte, 20)
		_, err = hB.Read(res)
		// color.Yellow(string(res[:n]))
		if err != nil {
			rmRelay <- c
			break
		}
	}
}

// Initially check the number of the total relays available
func sendNumber(num net.Conn, totalRelays *int) {
	for {
		temp := make([]byte, 100)

		t := *totalRelays
		a := strconv.Itoa(t)
		// fmt.Println("Sending ", a)
		//TO check that a request is made
		_, err := num.Read(temp)
		if err != nil {
			break
		}
		// fmt.Println("Check ", string(temp))
		//writing number to relay
		_, err = num.Write([]byte(a))
		if err != nil {
			break
		}
	}
	// fmt.Println("\t\t\tExiting1\n")

}

//ServerRequest is a method which will take request from the client
func (c Relay) ServerRequest(requestchan chan<- Request) {
	for {
		for i := 0; i < 1000; i++ {

		}
		time.Sleep(1 * time.Second)
		// fmt.Println("\tStarting read")
		link := make([]byte, 100)
		n, err := c.conn.Read(link)
		// fmt.Printf("The number of values read is %v\n", n)
		// fmt.Printf("The read value is %vEND", string(link[:n]))

		if err != nil {
			color.Red("\t\tConnection from %v is closed", c.number)
			break
		} else {
			req := Request{
				url:   string(link[:n]),
				relay: &c,
			}
			// c.getRequest(string(link[:n]))
			requestchan <- req
			//need to add a delay
		}

	}
	// fmt.Println("\t\t\tExiting2\n")

}

func (c Relay) getRequest(link string) {
	fmt.Print("The request found is " + link + "||\n")

	_, err := url.ParseRequestURI(link)
	if err != nil {
		color.Yellow(link)
		color.Magenta("Url is not correct")
		// c.conn.Write([]byte("The url received is not correct"))
		return
	}
	res, err := http.Get(link)
	if err != nil {
		color.Magenta("The url is not correct")
		// log.Fatal(err)
	}
	responseData, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(responseData))

	c.conn.Write(responseData)
	defer res.Body.Close()
	if err != nil {
		color.Yellow("Can not send the data")
	}
}

//WriteLinesFrom is a method no use
func (c Relay) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.conn, msg)
		if err != nil {
			break
		}
	}
}
func handleRelays(relaysDatabse map[int]Relay, requestchan <-chan Request, addRelay <-chan Relay, rmRelay <-chan Relay, totalRelays *int) {
	// relaysDatabse := make(map[int]Relay)

	for {
		select {
		case relay := <-requestchan:
			color.Magenta("New request: %s", relay.url)
			relay.relay.getRequest(relay.url)
		case relay := <-addRelay:
			relaysDatabse[relay.number] = relay
			displayTable(relaysDatabse)
			*totalRelays++
		case relay := <-rmRelay:
			color.Yellow("Relay# %v which is %v is down\n", relay.number, relay.name)
			delete(relaysDatabse, relay.number)
			displayTable(relaysDatabse)
			*totalRelays--
		}
	}
}

func displayTable(relaysDatabse map[int]Relay) {
	table := termtables.CreateTable()
	table.AddHeaders("Number", "Relay name", "Particpating")
	for _, value := range relaysDatabse {
		table.AddRow(value.number, value.name, value.participate)
	}
	color.Yellow(table.Render())
}
