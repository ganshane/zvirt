package zvirt

import (
	"testing"
	pb "github.com/ganshane/zvirt/protocol"
	"github.com/facebookgo/ensure"
	"time"
)
func Test_Define(t *testing.T) {
	zvirt := newTestInstance()
	defer zvirt.close()
	zvirt.initInstance()

	request := pb.DomainDefineRequest{Xml:`
	<domain type="test">
		<name>` + time.Now().String() + `</name>
		<memory unit="KiB">8192</memory>
		<os>
			<type>hvm</type>
		</os>
	</domain>`}
	response,err :=zvirt.domain.Define(nil,&request)

	ensure.Nil(t,err)
	ensure.NotNil(t,response.Uuid)

	_,err =zvirt.domain.Start(nil,&pb.DomainUUID{Uuid:response.Uuid})
	ensure.Nil(t,err)
	_,err =zvirt.domain.Shutdown(nil,&pb.DomainUUID{Uuid:response.Uuid})
	ensure.Nil(t,err)
	_,err =zvirt.domain.Destroy(nil,&pb.DomainUUID{Uuid:response.Uuid})
	ensure.Nil(t,err)
}

func Test_Domstate(t *testing.T) {
	zvirt := newTestInstance()
	defer zvirt.close()
	zvirt.initInstance()



	request :=pb.DomainUUID{Uuid:"asdf"}
	response,err :=zvirt.domain.Domstate(nil,&request)

	ensure.NotNil(t,err)

	dom := zvirt.buildTestDomain()
	defer dom.Destroy()


	uuid,err := dom.GetUUIDString()
	ensure.Nil(t,err)

	request =pb.DomainUUID{Uuid:uuid}
	response,err =zvirt.domain.Domstate(nil,&request)
	ensure.DeepEqual(t,response.State,pb.DomainState_VIR_DOMAIN_SHUTOFF)


	dom = zvirt.buildTransientTestDomain()
	defer dom.Destroy()

	uuid,err = dom.GetUUIDString()
	ensure.Nil(t,err)

	request =pb.DomainUUID{Uuid:uuid}
	response,err =zvirt.domain.Domstate(nil,&request)
	ensure.DeepEqual(t,response.State,pb.DomainState_VIR_DOMAIN_RUNNING)

}
