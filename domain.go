package zvirt

import (
	. "github.com/ganshane/zvirt/protocol"
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
type ZvirtDomain struct {
	agent *ZvirtAgent
}
// DomState implements zvirt_domain.DomState
func (zd *ZvirtDomain) Domstate(contxt context.Context, request *DomainUUID) (*DomainStateResponse, error){
	poolConn,err := zd.agent.pool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.pool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	dom,err :=conn.LookupDomainByUUIDString(request.GetUuid())
	if err != nil {
		return nil, err
	}else {
		defer dom.Free()
		domState, _, err := dom.GetState()
		if err != nil {
			return nil, err
		} else {
			response := DomainStateResponse{State: DomainState(domState)}
			return &response, nil
		}
	}
}
func (zd *ZvirtDomain) Define(context.Context, *DomainDefineRequest) (*DomainUUID, error){
	return nil,nil
}
func (zd *ZvirtDomain) Start(context.Context, *DomainUUID) (*DomainStateResponse, error){
	return nil,nil
}
func (zd *ZvirtDomain) Stop(context.Context, *DomainUUID) (*DomainStateResponse, error){
	return nil,nil
}
func (zd *ZvirtDomain) Destroy(context.Context, *DomainUUID) (*DomainStateResponse, error){
	return nil,nil
}
