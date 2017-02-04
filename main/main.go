
package main

import (
	"log"
	"net"
	libvirt "github.com/libvirt/libvirt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/ganshane/zvirt"
	pb "github.com/ganshane/zvirt/protocol"
)

const (
	port = ":50051"
)

func BuildTestConnection() *libvirt.Connect {
	conn, err := libvirt.NewConnect("test:///default")
	if err != nil {
		panic(err)
	}
	return conn
}
func main() {
	//build connection
	conn:= BuildTestConnection()
	defer conn.Close()

	server:=zvirt.NewServer(conn)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterZvirtDomainServiceServer(s,server)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}else{
		log.Printf("success")
	}
}
