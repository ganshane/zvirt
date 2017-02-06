package zvirt

import (
	"testing"
	pb "github.com/ganshane/zvirt/protocol"
	"github.com/facebookgo/ensure"
)

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
