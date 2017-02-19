package zvirt

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/facebookgo/ensure"
	"github.com/facebookgo/rpool"
	pb "github.com/ganshane/zvirt/protocol"
	libvirt "github.com/libvirt/libvirt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	uri  = flag.String("uri", "test:///default", "libvirtd connection uri")
	bind = flag.String("bind", ":50051", "zvirt bind string")
)

//libvirt connection wrapper
type libvirtConnWrapper struct {
	conn *libvirt.Connect
}

func (conn *libvirtConnWrapper) Close() error {
	_, e := conn.conn.Close()
	return e
}

//new connection
func newConnection() (io.Closer, error) {
	conn, err := libvirt.NewConnect(*uri)
	if err != nil {
		return nil, err
	}
	return &libvirtConnWrapper{conn: conn}, nil
}

//Agent global agent server instance
type Agent struct {
	connectionPool *rpool.Pool  //resource pool
	listener       net.Listener //network binding
	rpc            *grpc.Server //rpc server instance

	domain *Domain //domain
	pool   *Pool   //storage pool
	Volume *Volume //volume
}

//Fatal log fatal message
func (agent *Agent) Fatal(args ...interface{}) {
	log.Fatal(args)
}

//start zvirt agent
func (agent *Agent) initInstance() {
	//create domain service instance
	domain := &Domain{agent: agent}
	agent.domain = domain
	pb.RegisterZvirtDomainServiceServer(agent.rpc, domain)

	pool := &Pool{agent: agent}
	agent.pool = pool
	pb.RegisterZvirtPoolServiceServer(agent.rpc, pool)

	volume := &Volume{agent:agent}
	agent.Volume = volume
	pb.RegisterZvirtVolumeServiceServer(agent.rpc,volume)

	//register grpc service
	reflection.Register(agent.rpc)

	//handle system singal
	go agent.handleSignal()
}

//Serve initialize agent instance and start rpc server
func (agent *Agent) Serve() {
	agent.initInstance()

	log.Println("starting rpc server....")
	if err := agent.rpc.Serve(agent.listener); err != nil {
		log.Println("failed to serve: ", err)
	}
}

//handle system singal
func (agent *Agent) handleSignal() {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	for {
		sig := <-ch
		log.Println("receive signal ", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// this ensures a subsequent INT/TERM will trigger standard go behaviour of
			// terminating.
			signal.Stop(ch)
			agent.close()
			return
		case syscall.SIGUSR2:
			// we only return here if there's an error, otherwise the new process
			// will send us a TERM when it's ready to trigger the actual shutdown.
		}
	}
}

//close zvirt agent
func (agent *Agent) close() error {
	log.Println("closing zvirt agent")

	log.Println("closing rpc server")
	agent.rpc.GracefulStop()

	if agent.listener != nil {
		log.Println("closing net listener")
		agent.listener.Close()
	}

	log.Println("closing pool")
	agent.connectionPool.Close()

	log.Println("zvirt agent closed")
	return nil
}
func (agent *Agent) executeInConnection(callback func(*libvirt.Connect) (interface{}, error)) (interface{}, error) {
	conn, err := agent.connectionPool.Acquire()
	ensure.Nil(agent, err)
	defer agent.connectionPool.Release(conn)
	libvirtConn := conn.(*libvirtConnWrapper).conn
	return callback(libvirtConn)
}

//NewServer create new server for zvirt
func NewServer() *Agent {
	//resource pool
	p := rpool.Pool{
		New:           newConnection,
		Max:           5,
		MinIdle:       1,
		IdleTimeout:   time.Hour,
		ClosePoolSize: 2,
	}
	log.Printf("starting zvirt agent for %v@%v ", *uri, *bind)
	//bind
	lis, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//new rpc server
	s := grpc.NewServer()
	return &Agent{connectionPool: &p, listener: lis, rpc: s}
}

//only for test
func newTestInstance() *Agent {
	//resource pool
	p := rpool.Pool{
		New:           newConnection,
		Max:           5,
		MinIdle:       1,
		IdleTimeout:   time.Hour,
		ClosePoolSize: 2,
	}
	//new rpc server
	s := grpc.NewServer()
	return &Agent{connectionPool: &p, listener: nil, rpc: s}
}
