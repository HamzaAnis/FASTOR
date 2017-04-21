package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
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
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	} else {
		port = ":" + "9696"
	}
	connection, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	requestchan := make(chan string)

	addRelay := make(chan Relay)

	rmRelay := make(chan Relay)

	go handleRelays(requestchan, addRelay, rmRelay)

	for {
		conn, err := connection.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			fmt.Println("A client has connected")
			go handleConnection(conn, requestchan, addRelay, rmRelay, &numberOfRelays)
		}
	}
}

func promptName(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "What is your relay name? ")
	name, _, _ := bufc.ReadLine()
	return string(name)
}

func promptChoice(c net.Conn, bufc *bufio.Reader) bool {
	io.WriteString(c, "Do you want to participate in the anonymous service?(Y/N)")
	choice, _, _ := bufc.ReadLine()
	participate := false
	if string(choice) == "Y" {
		participate = true
		fmt.Println("Relay is participating")
	} else if string(choice) == "N" {
		participate = false
		fmt.Println("Relay is not participating")
	}
	return participate
}

//Core
func handleConnection(c net.Conn, requestchan chan<- string, addRelay chan<- Relay, rmRelay chan<- Relay, numberRelay *int) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	//we first need to add current relay to the channel
	//filling in the relay structure
	relay := Relay{
		conn:        c,
		name:        promptName(c, bufc),
		ch:          make(chan string),
		participate: promptChoice(c, bufc),
		number:      *numberRelay,
	}
	*numberRelay++
	if strings.TrimSpace(relay.name) == "" {
		io.WriteString(c, "Invalid relay name\n")
		return
	}

	// Register user, our messageHandler is waiting on this channel
	//it populates the map
	addRelay <- relay

	//ignore for the time being
	defer func() {
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmRelay <- relay
	}()

	//just a welcome message
	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\n", relay.name))

	//We are now populating the other channel now
	//our message handler is waiting on this channel as well
	//it reads this message and copies to the individual channel of each Client in map
	// effectively the broadcast

	// another go routine whose purpose is to keep on waiting for user input
	//and write it with nick to the
	go relay.ReadLinesInto(requestchan)

	//given a channel, writelines prints lines from it
	//we are giving here client.ch and this routine is for each client
	//so effectively each client is printitng its channel
	//to which our messagehandler has added messages for boroadcast
	relay.WriteLinesFrom(relay.ch)
}

//ReadLinesInto is a method on Client type
//it keeps waiting for user to input a line, ch chan is the msgchannel
//it formats and writes the message to the channel
func (c Relay) ReadLinesInto(ch chan<- string) {
	bufc := bufio.NewReader(c.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		ch <- fmt.Sprintf("%s: %s", c.name, line)
	}
}

//WriteLinesFrom is a method
//each client routine is writing to channel
func (c Relay) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.conn, msg)
		if err != nil {
			return
		}
	}
}
func handleRelays(requestchan <-chan string, addRelay <-chan Relay, rmRelay <-chan Relay) {
	relays := make(map[net.Conn]chan<- string)

	for {
		select {
		case site := <-requestchan:
			log.Printf("New request: %s", site)
			for _, ch := range relays {
				go func(mch chan<- string) { mch <- "\033[1;33;40m" + site + "\033[m" }(ch)
			}
		case relay := <-addRelay:
			log.Printf("New relay: %v\n", relay.conn)
			relays[relay.conn] = relay.ch
		case relay := <-rmRelay:
			log.Printf("Relay disconnects: %v\n", relay.conn)
			delete(relays, relay.conn)
		}
	}
}
