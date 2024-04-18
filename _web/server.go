package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":1234", "http server address listen on")
)

// go run server.go
func main() {
	flag.Parse()

	log.Println("http listen:", *addr)
	log.Fatal(http.ListenAndServe(*addr, http.FileServer(http.Dir("./"))))
}
