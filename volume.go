package zvirt

import (
	"github.com/facebookgo/ensure"
	pb "github.com/ganshane/zvirt/protocol"
	"golang.org/x/net/context"
)

//Pool storage pool for libvirt
type Volume struct {
	agent *Agent
}

//Create create volume
func (zvolume *Volume) Create(ctx context.Context, request *pb.VolumeDefineRequest) (*pb.VolumeInfo, error) {
	poolConn, err := zvolume.agent.connectionPool.Acquire()
	ensure.Nil(zvolume.agent, err)
	defer zvolume.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	pool, err := conn.LookupStoragePoolByUUIDString(request.PoolUuid)
	if err != nil {
		return nil, err
	}
	defer pool.Free()
	vol, err := pool.StorageVolCreateXML(request.Xml, 0)
	if err != nil {
		return nil, err
	}
	defer vol.Free()
	info, err := vol.GetInfo()
	if (err != nil) {
		return nil, err
	}
	name, _ := vol.GetName()
	return &pb.VolumeInfo{
		Name:name,
		Capacity:info.Capacity,
		Allocation:info.Allocation,
	}, nil
}

