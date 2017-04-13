package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

func main() {

	// s := "postgres://user:pass@host.com:5432/path.sf?k=v#f"

	// s := "http://www.google.com/hamzaa/nis.png:8008"

	s := "http://www.google.com/hamzaanis.js"
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	fmt.Println("Scheme is " + u.Scheme)

	fmt.Println("Host is " + u.Host)
	host, port, _ := net.SplitHostPort(u.Host)
	fmt.Println("Host is " + host)
	fmt.Println(port)

	typee := u.Path[1:]
	filet := strings.Split(typee, ".")
	if filet[len(filet)-1] != "html" {
		fmt.Println("The image is found named " + typee)
	}
	fmt.Printf("Path is %v\n", u.Path)

}

// package main

// import (
//     "fmt"
//     "strings"
// )

// func get
// func main() {
// 	url := "www.i.imgur.com/m1UIjW1.jpg"
//     s := strings.Split(url, "/")
// 	for _, f := range s {
// 		fmt.Printf("%v\n",f)
// 	}
// 	a:=len(s)
// 	fmt.Println(s[a-1])

// }
// package main
// import (
//     "fmt"
//     "io"
//     "log"
//     "net/http"
//     "os"
// )

// func main() {
//     url := "http://i.imgur.com/m1UIjW1.jpg"
//     // don't worry about errors
//     response, e := http.Get(url)
//     if e != nil {
//         log.Fatal(e)
//     }

//     defer response.Body.Close()

//     //open a file for writing
//     file, err := os.Create("asdf.jpg")
//     if err != nil {
//         log.Fatal(err)
//     }
//     // Use io.Copy to just dump the response body to the file. This supports huge files
//     _, err = io.Copy(file, response.Body)
//     if err != nil {
//         log.Fatal(err)
//     }
//     file.Close()
//     fmt.Println("Success!")
// }
