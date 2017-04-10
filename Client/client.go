// client.go
package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6000")
	if err != nil {
		// handle error
	}
	fmt.Println("Connection successful!!", conn.RemoteAddr())
	recvdSlice := make([]byte, 11)
	conn.Read(recvdSlice)
	fmt.Println(string(recvdSlice))
}
