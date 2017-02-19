package zvirt

import (
	"github.com/facebookgo/ensure"
	pb "github.com/ganshane/zvirt/protocol"
	"github.com/libvirt/libvirt-go"
	"golang.org/x/net/context"
)

//Volume storage volume for libvirt
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
	if err != nil {
		return nil, err
	}
	key, err := vol.GetKey()
	if err != nil {
		return nil, err
	}

	return &pb.VolumeInfo{
		Key:        key,
		Capacity:   info.Capacity,
		Allocation: info.Allocation,
	}, nil
}

//Info Get volume information
func (zvolume *Volume) Info(ctx context.Context, request *pb.VolumeRequest) (*pb.VolumeInfo, error) {
	poolConn, err := zvolume.agent.connectionPool.Acquire()
	ensure.Nil(zvolume.agent, err)
	defer zvolume.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	vol, err := conn.LookupStorageVolByKey(request.Key)
	if err != nil {
		return nil, err
	}
	defer vol.Free()

	key, err := vol.GetKey()
	if err != nil {
		return nil, err
	}

	info, err := vol.GetInfo()
	if err != nil {
		return nil, err
	}

	return &pb.VolumeInfo{
		Key:        key,
		Capacity:   info.Capacity,
		Allocation: info.Allocation,
	}, nil

}

//Delete delete volume
func (zvolume *Volume) Delete(ctx context.Context, request *pb.VolumeRequest) (*pb.VolumeInfo, error) {
	poolConn, err := zvolume.agent.connectionPool.Acquire()
	ensure.Nil(zvolume.agent, err)
	defer zvolume.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	vol, err := conn.LookupStorageVolByKey(request.Key)
	if err != nil {
		return nil, err
	}
	defer vol.Free()

	err = vol.Delete(libvirt.STORAGE_VOL_DELETE_NORMAL)
	if err != nil {
		return nil, err
	}

	return &pb.VolumeInfo{}, nil
}
