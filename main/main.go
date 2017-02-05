
package main

import (
	"log"
	"github.com/ganshane/zvirt"
)

const (
	port = ":50051"
)

func main() {
	uri := "test:///default"
	log.Println("starting zvirt agent for ",uri)
	server:=zvirt.NewServer(uri,port)
	server.Start()
}
