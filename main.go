package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"./handlers"
)

var (
    HOST string
    PORT int
)

func init() {
    flag.StringVar(&HOST, "host", "127.0.0.1", "host of http serve")
    flag.IntVar(&PORT, "port", 8000, "port of http serve")
    flag.Parse()
    fmt.Printf("Serving HTTP http://%s:%d/\n", HOST, PORT)
}

func main() {
    handlers.Setup()
    log.Fatal(http.ListenAndServe(
        fmt.Sprintf("%s:%d", HOST, PORT), nil))
}
