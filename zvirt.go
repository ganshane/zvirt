package zvirt
import (
	libvirt "github.com/libvirt/libvirt-go"
	"io"
	"github.com/facebookgo/rpool"
	"time"
	"log"
	"os"
	"syscall"
	"os/signal"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/ganshane/zvirt/protocol"
)

//libvirt connection wrapper
type libvirtConnWrapper struct {
	conn *libvirt.Connect
}
func (conn *libvirtConnWrapper)Close() error  {
	_, e :=conn.conn.Close()
	return e;
}

//ConnectionMaker for ResourcePool
type connectionMaker struct{
	uri string  //global uri
}
//new connection
func (r *connectionMaker) new() (io.Closer, error) {
	conn, err := libvirt.NewConnect(r.uri)
	if err != nil {
		return nil,err
	}
	return &libvirtConnWrapper{conn:conn},nil
}

//global ZvirtAgent server
type ZvirtAgent struct {
	pool *rpool.Pool            //resource pool
	listener net.Listener       //network binding
	rpc *grpc.Server            //rpc server instance
}
//fatal method
func(s *ZvirtAgent) Fatal(args ...interface{})()  {
	log.Fatal(args)
}
//start zvirt agent
func (s*ZvirtAgent) Start() {
	//register grpc service
	pb.RegisterZvirtDomainServiceServer(s.rpc,s)
	reflection.Register(s.rpc)
	//handle system singal
	go s.handleSignal()
	log.Println("starting rpc server....")
	if err := s.rpc.Serve(s.listener); err != nil {
		log.Println("failed to serve: ", err)
	} else {
		log.Printf("success")
	}
}
func (s*ZvirtAgent) handleSignal(){
	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	for{
		sig := <-ch
		log.Println("receive signal ",sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// this ensures a subsequent INT/TERM will trigger standard go behaviour of
			// terminating.
			signal.Stop(ch)
			s.close()
			return
		case syscall.SIGUSR2:
		// we only return here if there's an error, otherwise the new process
		// will send us a TERM when it's ready to trigger the actual shutdown.
		}
	}
}
//close zvirt agent
func (s*ZvirtAgent) close() error {
	log.Println("closing zvirt agent")
	log.Println("closing rpc server")
	s.rpc.GracefulStop()
	log.Println("closing net listener")
	s.listener.Close()
	log.Println("closing pool")
	log.Println("zvirt agent closed")
	return nil
}
//create new server for zvirt
func NewServer(uri string,bind string) *ZvirtAgent {
	//resource pool
	maker := connectionMaker{uri:uri}
	p := rpool.Pool{
		New:           maker.new,
		Max:           5,
		MinIdle:       1,
		IdleTimeout:   time.Hour,
		ClosePoolSize: 2,
	}
	//bind
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//new rpc server
	s := grpc.NewServer()
	return &ZvirtAgent{pool:&p,listener:lis,rpc:s}
}
