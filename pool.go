package zvirt

import (
	"github.com/facebookgo/ensure"
	pb "github.com/ganshane/zvirt/protocol"
	"github.com/libvirt/libvirt-go"
	"golang.org/x/net/context"
)

//Pool storage pool for libvirt
type Pool struct {
	agent *Agent
}

//Define define an inactive persistent storage pool or modify an existing persistent one from an XML file
func (zpool *Pool) Define(ctx context.Context, request *pb.PoolDefineRequest) (*pb.PoolUUID, error) {
	poolConn, err := zpool.agent.connectionPool.Acquire()
	ensure.Nil(zpool.agent, err)
	defer zpool.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	pool, err := conn.StoragePoolDefineXML(request.Xml, 0)
	if err != nil {
		return nil, err
	}
	defer pool.Free()
	uuid, err := pool.GetUUIDString()
	if err != nil {
		return nil,err
	}
	return &pb.PoolUUID{Uuid: uuid}, nil
}

//Info - storage pool information
func (zpool *Pool) Info(ctx context.Context, uuid *pb.PoolUUID) (*pb.PoolStateResponse, error) {
	poolConn, err := zpool.agent.connectionPool.Acquire()
	ensure.Nil(zpool.agent, err)
	defer zpool.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	pool, err := conn.LookupStoragePoolByUUIDString(uuid.Uuid)
	if err != nil {
		return nil, err
	}
	defer pool.Free()
	info, err := pool.GetInfo()
	if err != nil {
		return nil,err
	}
	return &pb.PoolStateResponse{State: pb.PoolState(info.State)}, nil
}

//Start start a (previously defined) inactive pool
func (zpool *Pool) Start(ctx context.Context, poolUUID *pb.PoolUUID) (*pb.PoolStateResponse, error) {
	poolConn, err := zpool.agent.connectionPool.Acquire()
	ensure.Nil(zpool.agent, err)
	defer zpool.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	pool, err := conn.LookupStoragePoolByUUIDString(poolUUID.Uuid)
	if err != nil {
		return nil, err
	}
	defer pool.Free()
	err = pool.Create(libvirt.STORAGE_POOL_CREATE_NORMAL)
	if err != nil {
		return nil, err
	}
	return &pb.PoolStateResponse{State: pb.PoolState_STORAGE_POOL_RUNNING}, nil
}

//Destroy destroy (stop) a pool
func (zpool *Pool) Destroy(ctx context.Context, poolUUID *pb.PoolUUID) (*pb.PoolStateResponse, error) {
	poolConn, err := zpool.agent.connectionPool.Acquire()
	ensure.Nil(zpool.agent, err)
	defer zpool.agent.connectionPool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	pool, err := conn.LookupStoragePoolByUUIDString(poolUUID.Uuid)
	if err != nil {
		return nil, err
	}
	defer pool.Free()
	err = pool.Destroy()
	if err != nil {
		return nil, err
	}
	return &pb.PoolStateResponse{State: pb.PoolState_STORAGE_POOL_INACCESSIBLE}, nil
}
