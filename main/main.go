
package main

import (
	"github.com/ganshane/zvirt"
	"flag"
)


func main() {
	flag.Parse()
	server:=zvirt.NewServer()
	server.Start()
}
