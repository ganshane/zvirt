package zvirt

import (
	"testing"
	pb "github.com/ganshane/zvirt/protocol"
	"github.com/facebookgo/ensure"
)
func TestPool_Define(t *testing.T) {
	zvirt := newTestInstance()
	defer zvirt.close()
	zvirt.initInstance()

	request := pb.PoolDefineRequest{Xml:`
  <pool type="logical">
    <name>virt</name>
  </pool> `}
	response, err := zvirt.zpool.Define(nil, &request)

	ensure.Nil(t,err)
	ensure.NotNil(t,response.Uuid)

	poolUUID := pb.PoolUUID{Uuid:response.Uuid}
	_, err = zvirt.zpool.Info(nil, &poolUUID)
	ensure.Nil(t,err)

	_, err =zvirt.zpool.Start(nil,&poolUUID)
	ensure.Nil(t,err)
	_, err =zvirt.zpool.Destroy(nil,&poolUUID)
	ensure.Nil(t,err)
}
