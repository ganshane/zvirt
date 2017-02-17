package zvirt

import (
	"github.com/facebookgo/ensure"
	pb "github.com/ganshane/zvirt/protocol"
	"testing"
)

func TestPool_Define(t *testing.T) {
	zvirt := newTestInstance()
	defer zvirt.close()
	zvirt.initInstance()

	request := pb.PoolDefineRequest{Xml: `
  <pool type="logical">
    <name>virt</name>
  </pool> `}
	response, err := zvirt.pool.Define(nil, &request)

	ensure.Nil(t, err)
	ensure.NotNil(t, response.Uuid)

	poolUUID := pb.PoolUUID{Uuid: response.Uuid}
	_, err = zvirt.pool.Info(nil, &poolUUID)
	ensure.Nil(t, err)

	_, err = zvirt.pool.Start(nil, &poolUUID)
	ensure.Nil(t, err)
	_, err = zvirt.pool.Destroy(nil, &poolUUID)
	ensure.Nil(t, err)
}
