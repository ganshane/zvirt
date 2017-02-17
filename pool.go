package zvirt

import (
	. "github.com/ganshane/zvirt/protocol"
	"golang.org/x/net/context"
	"github.com/facebookgo/ensure"
	"github.com/libvirt/libvirt-go"
)


type ZvirtPool struct {
	agent *ZvirtAgent
}
func (zd *ZvirtPool) Define(ctx context.Context, request *PoolDefineRequest) (*PoolUUID, error){
	poolConn,err := zd.agent.pool.Acquire()
	ensure.Nil(zd.agent, err)
	defer zd.agent.pool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	pool,err :=conn.StoragePoolDefineXML(request.Xml,0)
	if err == nil {
		defer pool.Free()
		if uuid, err := pool.GetUUIDString(); err == nil {
			return &PoolUUID{Uuid:uuid}, nil
		}
	}
	return nil,err
}
func (zpool *ZvirtPool) Info(ctx context.Context, uuid *PoolUUID) (*PoolStateResponse, error){
	poolConn,err := zpool.agent.pool.Acquire()
	ensure.Nil(zpool.agent, err)
	defer zpool.agent.pool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	if pool,err :=conn.LookupStoragePoolByUUIDString(uuid.Uuid);err == nil{
		defer pool.Free()
		if info,err:=pool.GetInfo(); err == nil {
			return &PoolStateResponse{State:PoolState(info.State)},nil
		}
	}
	return nil,err
}
func (zpool *ZvirtPool) Start(ctx context.Context, poolUUID *PoolUUID) (*PoolStateResponse, error){
	poolConn,err := zpool.agent.pool.Acquire()
	ensure.Nil(zpool.agent, err)
	defer zpool.agent.pool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	pool,err :=conn.LookupStoragePoolByUUIDString(poolUUID.Uuid)
	if err == nil{
		defer pool.Free()
		if err:=pool.Create(libvirt.STORAGE_POOL_CREATE_NORMAL); err == nil {
			return &PoolStateResponse{State:PoolState_STORAGE_POOL_RUNNING},nil
		}
	}
	return nil,err
}
func (zpool *ZvirtPool) Destroy(ctx context.Context, poolUUID *PoolUUID) (*PoolStateResponse, error){
	poolConn,err := zpool.agent.pool.Acquire()
	ensure.Nil(zpool.agent, err)
	defer zpool.agent.pool.Release(poolConn)
	conn := poolConn.(*libvirtConnWrapper).conn

	pool,err :=conn.LookupStoragePoolByUUIDString(poolUUID.Uuid)
	if err == nil{
		defer pool.Free()
		if err:=pool.Destroy(); err == nil {
			return &PoolStateResponse{State:PoolState_STORAGE_POOL_INACCESSIBLE},nil
		}
	}
	return nil,err
}
