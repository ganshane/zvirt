package zvirt

import (
	pb "github.com/ganshane/zvirt/protocol"
	"golang.org/x/net/context"
	"github.com/facebookgo/ensure"
	"github.com/libvirt/libvirt-go"
	"time"
)

//only for test
func (s*ZvirtAgent) buildTestDomain()(*libvirt.Domain) {
	conn,err := s.pool.Acquire()
	ensure.Nil(s, err)
	defer s.pool.Release(conn)

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
func (s*ZvirtAgent) buildTransientTestDomain() (*libvirt.Domain) {
	conn,err := s.pool.Acquire()
	ensure.Nil(s, err)
	defer s.pool.Release(conn)
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
// DomState implements zvirt_domain.DomState
func (s *ZvirtAgent) DomState(contxt context.Context, request *pb.DomStateRequest) (*pb.DomStateResponse, error){
	conn,err := s.pool.Acquire()
	ensure.Nil(s, err)
	defer s.pool.Release(conn)

	libvirtConn := conn.(*libvirtConnWrapper).conn
	dom,err :=libvirtConn.LookupDomainByUUIDString(request.GetVmUuid())
	if err != nil {
		return nil, err
	}else{
		defer dom.Free()
		domState,_,err:=dom.GetState()
		if err!= nil {
			return nil, err
		}else{
			response := pb.DomStateResponse{
				State: pb.DomainState(domState),
			}
			return &response,nil
		}
	}
}
