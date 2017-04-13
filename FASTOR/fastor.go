// package main

// import (
// 	"fmt"
// 	// "gorilla/http"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"net/url"
// 	"strings"
// )

// //A function to check the request and do the corrosponding action
// //E.g For image save the image to the file
// //If stylesheet save the stylesheet
// func handleStyle(link string){
// 	u, err := url.Parse(link)
// 	if err != nil {
// 		panic(err)
// 	}
// 	path:=u.Path[1:] //to avoid the / at the beginning
// 	filetype:=strings.Split(path,".")
// 	fmt.Printf("Type is %v\n",filetype[len(filetype)-1])
// 	if(filetype!="html")
// 	{

// 	}
// }

// //on every / request the handler will call this
// func torhandler(w http.ResponseWriter, r *http.Request) {
// 	domain := r.URL.Path[1:7]
// 	link :="http://" + r.URL.Path[8:]

// 	handleStyle(link)

// 	// if domain == "fastor" {
// 		fmt.Printf("Domain: %v\n", domain)
// 		fmt.Printf("Website:  %v\n", link)
// 		res, err := http.Get(link)
// 		if err != nil {
// 			// log.Fatal(err)
// 					http.Error(w,err.Error(),500)
// 			}
// 		buf, err := ioutil.ReadAll(res.Body)
// 		fmt.Fprintf(w, string(buf))
// 		defer res.Body.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	// } else {
// 	// 	fmt.Fprintf(w, "There is no services available for which you requested")
// 	// 	html, er := ioutil.ReadFile("test.html")
// 	// 	if er == nil {
// 	// 		fmt.Println("The file is  read")
// 	// 	} else {
// 	// 		fmt.Println("The file is not read")
// 	// 	}
// 	// 	fmt.Fprint(w, string(html))
// 	// }

// }

// func main() {
// 	port := ""
// 	if len(os.Args) > 1 {
// 		port = os.Args[1]
// 	} else {
// 		port = "9999"
// 	}
// 	//default will go here
// 	http.HandleFunc("/fastor/", torhandler)
//     http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("."))))
// 	http.ListenAndServe(":"+port, nil)
// }
