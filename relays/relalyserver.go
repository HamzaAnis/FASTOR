package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

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
	if len(os.Args) > 1 {
		port = os.Args[1]
		temp, _ := strconv.Atoi(port)
		temp++
		relaynNumberPort = string(temp)
	} else {
		port = "9696"
		relaynNumberPort = "9697"
	}
	connection, err := net.Listen("tcp", ":"+port)
	connection2, err2 := net.Listen("tcp", ":"+relaynNumberPort)
	if err != nil || err2 != nil {
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
		if err != nil || err2 != nil {
			fmt.Println(err)
			continue
		} else {
			color.Red("\tA client has connected")
			conn.Write([]byte("Hello FASTOR user!"))
			go handleConnection(conn, conn2, requestchan, addRelay, rmRelay, &numberOfRelays)
		}
	}
}

func promptName(c net.Conn) string {
	clr := color.New(color.FgGreen)
	io.WriteString(c, "What is your relay name? ")
	name := make([]byte, 8)
	c.Read(name)
	clr.Printf("The name of the relay: %v\n", string(name))
	return string(name)
}

func promptChoice(c net.Conn) bool {
	c.Write([]byte("Do you want to participate in the anonymous service?(Y/N)"))
	choice := make([]byte, 1)
	c.Read(choice)
	// fmt.Printf("The choice is %vJ\n", choice)
	// fmt.Printf("The length of choice is %vJ\n", len(choice))
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
func handleConnection(c net.Conn, num net.Conn, requestchan chan<- string, addRelay chan<- Relay, rmRelay chan<- Relay, numberRelay *int) {
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

	// Register user, our messageHandler is waiting on this channel
	//it populates the map
	addRelay <- relay

	//ignore for the time being
	defer func() {
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmRelay <- relay
	}()

	//We are now populating the other channel now
	//our message handler is waiting on this channel as well
	//it reads this message and copies to the individual channel of each Client in map
	// effectively the broadcast

	// another go routine whose purpose is to keep on waiting for user input
	//and write it with nick to the
	go relay.ReadLinesInto(requestchan)

	//to send the nnumber of relays
	go sendNumber(num, numberRelay)

	//given a channel, writelines prints lines from it
	//we are giving here client.ch and this routine is for each client
	//so effectively each client is printitng its channel
	//to which our messagehandler has added messages for boroadcast
	relay.WriteLinesFrom(relay.ch)

}

// Initially check the number of the total relays available
func sendNumber(num net.Conn, numberRelay *int) {
	for {
		temp := make([]byte, 100)
		t := *numberRelay
		a := strconv.Itoa(t)
		// fmt.Println("Sending ", a)
		//TO check that a request is made
		num.Read(temp)
		// fmt.Println("Check ", string(temp))
		//writing number to relay
		num.Write([]byte(a))
	}

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
			color.Magenta("New request: %s", site)
			for _, ch := range relays {
				go func(mch chan<- string) { mch <- "\033[1;33;40m" + site + "\033[m" }(ch)
			}
		case relay := <-addRelay:
			color.Blue("New relay: %v\n\tNumber= %v\n\tParticipating= %v", relay.name, relay.number, relay.participate)
			relays[relay.conn] = relay.ch
		case relay := <-rmRelay:
			log.Printf("Relay disconnects: %v\n", relay.conn)
			delete(relays, relay.conn)
		}
	}
}
