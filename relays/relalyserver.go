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

func main() {
	numberOfRelays := 0
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
	requestchan := make(chan string)

	addRelay := make(chan Relay)

	rmRelay := make(chan Relay)

	go handleRelays(requestchan, addRelay, rmRelay)

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
			go handleConnection(conn, conn2, conn3, requestchan, addRelay, rmRelay, &numberOfRelays)
		}
	}
}

func promptName(c net.Conn) string {
	clr := color.New(color.FgGreen)
	io.WriteString(c, "What is your relay name? ")
	name := make([]byte, 20)
	n, _ := c.Read(name)
	clr.Println("The length is ", n)
	clr.Printf("The name of the relay: %v\n", string(name[:n]))
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
func handleConnection(c net.Conn, num net.Conn, heatrbeat net.Conn, requestchan chan<- string, addRelay chan<- Relay, rmRelay chan<- Relay, numberRelay *int) {
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
	go sendNumber(num, numberRelay)

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
func sendNumber(num net.Conn, numberRelay *int) {
	for {
		temp := make([]byte, 100)

		t := *numberRelay
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

//ServerRequest is a method on Client type
//it keeps waiting for user to input a line, ch chan is the msgchannel
//it formats and writes the message to the channel
func (c Relay) ServerRequest(requestchan chan<- string) {
	for {
		for i := 0; i < 1000; i++ {

		}
		// fmt.Println("\tStarting read")
		link := make([]byte, 100)
		n, err := c.conn.Read(link)
		// fmt.Printf("The number of values read is %v\n", n)
		// fmt.Printf("The read value is %vEND", string(link[:n]))

		if err != nil {
			color.Red("\t\tConnection from %v is closed", c.number)
			break
		} else {
			time.Sleep(1 * time.Second)
			fmt.Print("The request found is " + string(link[:n]) + "||\n")

			_, err := url.ParseRequestURI(string(link[:n]))
			if err != nil {
				color.Yellow(string(link[:n]) + "C")
				color.Magenta("Url is not correct")
				// c.conn.Write([]byte("The url received is not correct"))
				continue
			}
			res, err := http.Get(string(link[:n]))
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
			requestchan <- string(link)

			//need to add a delay
		}

	}
	// fmt.Println("\t\t\tExiting2\n")

}

//WriteLinesFrom is a method
//each client routine is writing to channel
func (c Relay) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.conn, msg)
		if err != nil {
			break
		}
	}
}
func handleRelays(requestchan <-chan string, addRelay <-chan Relay, rmRelay <-chan Relay) {
	relaysDatabse := make(map[int]Relay)

	for {
		select {
		case site := <-requestchan:
			color.Magenta("New request: %s", site)
		case relay := <-addRelay:
			relaysDatabse[relay.number] = relay
			table := termtables.CreateTable()
			table.AddHeaders("Number", "Relay name", "Particpating")
			for _, value := range relaysDatabse {
				table.AddRow(value.number, value.name, value.participate)
			}
			color.Yellow(table.Render())
		case relay := <-rmRelay:
			color.Yellow("Relay# %v which is %v is down\n", relay.number, relay.name)
			delete(relaysDatabse, relay.number)
		}
	}
}
