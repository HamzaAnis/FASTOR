// package main

// import (
// 	"fmt"
// 	"net"
// 	"net/url"
// 	"strings"
// )

// func main() {

// 	// s := "postgres://user:pass@host.com:5432/path.sf?k=v#f"

// 	// s := "http://www.google.com/hamzaa/nis.png:8008"

// 	s := "http://www.google.com/hamzaanis.js"
// 	u, err := url.Parse(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Scheme is " + u.Scheme)

// 	fmt.Println("Host is " + u.Host)
// 	host, port, _ := net.SplitHostPort(u.Host)
// 	fmt.Println("Host is " + host)
// 	fmt.Println(port)

// 	typee := u.Path[1:]
// 	filet := strings.Split(typee, ".")
// 	if filet[len(filet)-1] != "html" {
// 		fmt.Println("The image is found named " + typee)
// 	}
// 	fmt.Printf("Path is %v\n", u.Path)

// }

// // package main

// // import (
// //     "fmt"
// //     "strings"
// // )

// // func get
// // func main() {
// // 	url := "www.i.imgur.com/m1UIjW1.jpg"
// //     s := strings.Split(url, "/")
// // 	for _, f := range s {
// // 		fmt.Printf("%v\n",f)
// // 	}
// // 	a:=len(s)
// // 	fmt.Println(s[a-1])

// // }
// // package main
// // import (
// //     "fmt"
// //     "io"
// //     "log"
// //     "net/http"
// //     "os"
// // )

// // func main() {
// //     url := "http://i.imgur.com/m1UIjW1.jpg"
// //     // don't worry about errors
// //     response, e := http.Get(url)
// //     if e != nil {
// //         log.Fatal(e)
// //     }

// //     defer response.Body.Close()

// //     //open a file for writing
// //     file, err := os.Create("asdf.jpg")
// //     if err != nil {
// //         log.Fatal(err)
// //     }
// //     // Use io.Copy to just dump the response body to the file. This supports huge files
// //     _, err = io.Copy(file, response.Body)
// //     if err != nil {
// //         log.Fatal(err)
// //     }
// //     file.Close()
// //     fmt.Println("Success!")
// // }

// package main

// import "fmt"

// // define Dog object type
// type Dog struct {
// 	Name  string
// 	Color string
// }

// func main() {

// 	// create instance of object and set properties
// 	Spot := Dog{Name: "Spot", Color: "brown"}

// 	// get pointer of object
// 	SpotPointer := &Spot

// 	// modify field through pointer
// 	SpotPointer.Color = "black"

// 	fmt.Println(Spot.Color)

// }

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func printSlowly(s string, n int) {
// 	for i := 0; i < n; i++ {
// 		fmt.Println(i, s)
// 		time.Sleep(300 * time.Millisecond)
// 	}
// }

// func main() {

// 	// This is a normal function call.
// 	// Main() will finish this off before continuing.
// 	printSlowly("directly functioning", 3)

// 	// The go functions below will each spin off to happen each in their own thread,
// 	// Meaning they'll be called _concurrently_.

// 	// Calling the named function as a go routine.
// 	go printSlowly("red fish goroutine", 3)
// 	go printSlowly("blue fish goroutine", 3)

// 	// Call an anonymous function as a go routine.
// 	go func(ss string, nn int) {
// 		for i := 0; i < nn; i++ {
// 			fmt.Println(i, ss)
// 			time.Sleep(150 * time.Millisecond)
// 		}
// 	}("anony fish goroutine", 3)

// 	// Waits for a button to be pushed.
// 	// Try commenting this!
// 	var input string
// 	fmt.Scanln(&input) // Just push RETURN to finish the program.
// 	fmt.Println("DONE.")
// }

// package main

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"fmt"
// )

// type MyString struct {
// 	Length  int32
// 	Message [10]byte
// }

// type MyMessage struct {
// 	First   uint64
// 	Second  byte
// 	_       byte // padding
// 	Third   uint32
// 	Message MyString
// }

// func main() {
// 	buf := new(bytes.Buffer)
// 	a := MyMessage{
// 		First:   10,
// 		Second:  10,
// 		Third:   10,
// 		Message: MyString{0, [10]byte{'H', 'e', 'l', 'l', 'o', '\n'}},
// 	}
// 	b := MyMessage{
// 		First:   100,
// 		Second:  0,
// 		Third:   100,
// 		Message: MyString{0, [10]byte{'H', 'e', '\n'}},
// 	}
// 	test := []MyMessage{a, b}
// 	fmt.Println(test)
// 	err := binary.Write(buf, binary.LittleEndian, &test)
// 	if err != nil {
// 		fmt.Printf("binary.Read failed:", err)
// 		return
// 	}

// 	// <<--- CONN -->>
// 	// msg := []MyMessage{}
// 	msg2 := new(MyMessage)

// 	err2 := binary.Read(buf, binary.LittleEndian, msg2)
// 	if err2 != nil {
// 		fmt.Printf("binary.Read failed:", err2)
// 		return
// 	}
// 	fmt.Println(msg2)

// }

// package main

// import (
// 	"net/url"

// 	"github.com/fatih/color"
// )

// func main() {

// 	_, err := url.ParseRequestURI("http://www.google.com")
// 	if err != nil {
// 		color.Red("Invalid URI")
// 	}
// 	// link := "http://google.com"
// 	// u, err := url.ParseRequestURI(link)
// 	// fmt.Println(u.Scheme + "://" + u.Host)
// 	// if err != nil {
// 	// 	fmt.Println("The rl not correct")
// 	// }
// 	// response, err := http.Get(link)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// } else {
// 	// 	defer response.Body.Close()
// 	// 	io.Copy(os.Stdout, response.Body)

// 	// }
// }

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

var m = map[int]string{1: "one", 2: "two", 3: "three"}
var a string = "Hamza"

func main() {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(m)
	if err != nil {
		panic(err)
	}

	// your encoded stuff
	// fmt.Println(buf.Bytes())

	var decodedMap map[int]string
	decoder := gob.NewDecoder(buf)

	err = decoder.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", decodedMap)

	call()
}

func call() {
	fmt.Println(a)
}
