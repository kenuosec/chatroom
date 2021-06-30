package main

import (
	"chatroom/config"
	"chatroom/ctrl"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", "127.0.0.1:9000", "http service address")

func main() {
	config.NewLoggerWithRotate()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", ctrl.HandleConnections)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
