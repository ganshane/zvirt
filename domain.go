package zvirt

import (
	"time"

	"github.com/facebookgo/ensure"
	pb "github.com/ganshane/zvirt/protocol"
	"github.com/libvirt/libvirt-go"
	"golang.org/x/net/context"
)

//only for test
func (s *Agent) buildTestDomain() *libvirt.Domain {
	conn, err := s.connectionPool.Acquire()
	ensure.Nil(s, err)
	defer s.connectionPool.Release(conn)

	libvirtConn := conn.(*libvirtConnWrapper).conn
	dom, err := libvirtConn.DomainDefineXML(`<domain type="test">
		<name>` + time.Now().String() + `</name>
		<memory unit="KiB">8192</memory>
		<os>
			<type>hvm</type>
		</os>
	</domain>`)
	if err != nil {
		panic(err)
	}
	return dom
}

//only for test
func (s *Agent) buildTransientTestDomain() *libvirt.Domain {
	conn, err := s.connectionPool.Acquire()
	ensure.Nil(s, err)
	defer s.connectionPool.Release(conn)
	libvirtConn := conn.(*libvirtConnWrapper).conn

	dom, err := libvirtConn.DomainCreateXML(`<domain type="test">
		<name>`+time.Now().String()+`</name>
		<memory unit="KiB">8192</memory>
		<os>
			<type>hvm</type>
		</os>
	</domain>`, libvirt.DOMAIN_NONE)
	if err != nil {
		panic(err)
	}
	return dom
}

//Domain is a node domain manager .
type Domain struct {
	agent *Agent
}

// Domstate implements zvirt_domain.DomState
func (zd *Domain) Domstate(contxt context.Context, request *pb.DomainUUID) (*pb.DomainStateResponse, error) {
	poolConn, err := zd.agent.connectionPool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	dom, err := conn.LookupDomainByUUIDString(request.GetUuid())
	if err == nil {
		defer dom.Free()
		if domState, _, err := dom.GetState(); err == nil {
			response := pb.DomainStateResponse{State: pb.DomainState(domState)}
			return &response, nil
		}
	}
	return nil, err
}

//Define define (but don't start) a domain from an XML file
func (zd *Domain) Define(ctx context.Context, request *pb.DomainDefineRequest) (*pb.DomainUUID, error) {
	poolConn, err := zd.agent.connectionPool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	dom, err := conn.DomainDefineXML(request.Xml)
	if err == nil {
		defer dom.Free()
		if uuid, err := dom.GetUUIDString(); err == nil {
			return &pb.DomainUUID{Uuid: uuid}, nil
		}
	}
	return nil, err
}

//Start - start a (previously defined) inactive domain
func (zd *Domain) Start(ctx context.Context, request *pb.DomainUUID) (*pb.DomainStateResponse, error) {
	poolConn, err := zd.agent.connectionPool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	dom, err := conn.LookupDomainByUUIDString(request.Uuid)
	if err == nil {
		/*
			//see https://github.com/libvirt/libvirt/blob/master/tools/virsh-domain.c#L4097
			if id,_ :=dom.GetID();id != -1 {
				return nil,errors.New("Domain is already active")
			}
		*/
		flag := libvirt.DOMAIN_NONE
		if err = dom.CreateWithFlags(flag); err == nil {
			return &pb.DomainStateResponse{State: pb.DomainState_VIR_DOMAIN_RUNNING}, nil
		}
	}
	return nil, err
}

//Shutdown - gracefully shutdown a domain
func (zd *Domain) Shutdown(ctx context.Context, request *pb.DomainUUID) (*pb.DomainStateResponse, error) {

	poolConn, err := zd.agent.connectionPool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	dom, err := conn.LookupDomainByUUIDString(request.Uuid)
	if err == nil {
		if err = dom.Shutdown(); err == nil {
			return &pb.DomainStateResponse{State: pb.DomainState_VIR_DOMAIN_SHUTDOWN}, nil
		}
	}
	return nil, err
}

//Destroy - destroy (stop) a domain
func (zd *Domain) Destroy(ctx context.Context, request *pb.DomainUUID) (*pb.DomainStateResponse, error) {
	poolConn, err := zd.agent.connectionPool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	dom, err := conn.LookupDomainByUUIDString(request.Uuid)
	if err == nil {
		if err = dom.Destroy(); err == nil {
			return &pb.DomainStateResponse{State: pb.DomainState_VIR_DOMAIN_NOSTATE}, nil
		}
	}
	return nil, err
}
