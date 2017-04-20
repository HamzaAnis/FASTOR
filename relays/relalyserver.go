package main

import(
	"fmt"
	"net"
	"os"
)

type Relays struct {
	conn     net.Conn
	name string
	number int
	ch       chan string
}
func main(){
	fmt.Println("Hamza Anis")
	port:=""
	if len(os.Args) > 1 {
		port =":"+ os.Args[1]
	} else {
		port = ":"+"9696"
	}
	connection, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for {
			_	, err := connection.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}else{
				fmt.Println("A client has connected")
			}
	}
}




}

