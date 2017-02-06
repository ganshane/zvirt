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
func (zd *ZvirtDomain) Define(ctx context.Context, request *DomainDefineRequest) (*DomainUUID, error){
	poolConn,err := zd.agent.pool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.pool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	if dom,err :=conn.DomainDefineXML(request.Xml);err == nil{
		defer dom.Free()
		if uuid,err:=dom.GetUUIDString(); err == nil {
			return &DomainUUID{Uuid:uuid},nil
		}else{
			return nil,err
		}
	}else{
		return nil,err
	}
}
func (zd *ZvirtDomain) Start(ctx context.Context, request *DomainUUID) (*DomainStateResponse, error){
	var err error

	poolConn,err := zd.agent.pool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.pool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	if dom,err := conn.LookupDomainByUUIDString(request.Uuid);err == nil{
		/*
		//see https://github.com/libvirt/libvirt/blob/master/tools/virsh-domain.c#L4097
		if id,_ :=dom.GetID();id != -1 {
			return nil,errors.New("Domain is already active")
		}
		*/
		flag := libvirt.DOMAIN_NONE
		if err = dom.CreateWithFlags(flag);err == nil{
			return &DomainStateResponse{State:DomainState_VIR_DOMAIN_RUNNING},nil
		}
	}
	return nil,err
}
func (zd *ZvirtDomain) Shutdown(ctx context.Context, request *DomainUUID) (*DomainStateResponse, error){
	var err error

	poolConn,err := zd.agent.pool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.pool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn


	if dom,err := conn.LookupDomainByUUIDString(request.Uuid);err == nil{
		if err= dom.Shutdown();err == nil{
			return &DomainStateResponse{State:DomainState_VIR_DOMAIN_SHUTDOWN},nil
		}
	}
	return nil,err
}
func (zd *ZvirtDomain) Destroy(ctx context.Context, request *DomainUUID) (*DomainStateResponse, error){
	var err error

	poolConn,err := zd.agent.pool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.pool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	if dom,err := conn.LookupDomainByUUIDString(request.Uuid);err == nil{
		if err = dom.Destroy();err == nil {
			return &DomainStateResponse{State:DomainState_VIR_DOMAIN_NOSTATE}, nil
		}
	}
	return nil,err
}
