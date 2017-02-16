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
	"flag"
	"github.com/facebookgo/ensure"
)
var(
	uri = flag.String("uri", "test:///default", "libvirtd connection uri")
	bind = flag.String("bind", ":50051", "zvirt bind string")
)

//libvirt connection wrapper
type libvirtConnWrapper struct {
	conn *libvirt.Connect
}
func (conn *libvirtConnWrapper)Close() error  {
	_, e :=conn.conn.Close()
	return e;
}

//new connection
func  newConnection() (io.Closer, error) {
	conn, err := libvirt.NewConnect(*uri)
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

	//zd domain
	domain *ZvirtDomain
	zpool *ZvirtPool
}
//fatal method
func(s *ZvirtAgent) Fatal(args ...interface{})()  {
	log.Fatal(args)
}
//start zvirt agent
func (agent *ZvirtAgent) initInstance() {
	//create domain service instance
	domain := &ZvirtDomain{agent:agent}
	agent.domain = domain
	pb.RegisterZvirtDomainServiceServer(agent.rpc, domain)

	pool := &ZvirtPool{agent:agent}
	agent.zpool = pool
	pb.RegisterZvirtPoolServiceServer(agent.rpc,pool)

	//register grpc service
	reflection.Register(agent.rpc)

	//handle system singal
	go agent.handleSignal()
}
func (agent *ZvirtAgent) Serve(){
	agent.initInstance()

	log.Println("starting rpc server....")
	if err := agent.rpc.Serve(agent.listener); err != nil {
		log.Println("failed to serve: ", err)
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

	if(s.listener != nil) {
		log.Println("closing net listener")
		s.listener.Close()
	}

	log.Println("closing pool")
	s.pool.Close()

	log.Println("zvirt agent closed")
	return nil
}
func (agent*ZvirtAgent) executeInConnection(callback func(*libvirt.Connect)(interface{},error))(interface{},error){
	conn,err := agent.pool.Acquire()
	ensure.Nil(agent, err)
	defer agent.pool.Release(conn)
	libvirtConn := conn.(*libvirtConnWrapper).conn
	return callback(libvirtConn)
}
//create new server for zvirt
func NewServer() *ZvirtAgent {
	//resource pool
	p := rpool.Pool{
		New:           newConnection,
		Max:           5,
		MinIdle:       1,
		IdleTimeout:   time.Hour,
		ClosePoolSize: 2,
	}
	log.Printf("starting zvirt agent for %v@%v ",*uri,*bind)
	//bind
	lis, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//new rpc server
	s := grpc.NewServer()
	return &ZvirtAgent{pool:&p,listener:lis,rpc:s}
}
//only for test
func newTestInstance() *ZvirtAgent {
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
	return &ZvirtAgent{pool:&p,listener:nil,rpc:s}
}
