package main

import (
	"flag"
	"github.com/ganshane/zvirt"
)

func main() {
	flag.Parse()
	server := zvirt.NewServer()
	server.Serve()
}
