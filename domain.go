package zvirt

import (
	pb "github.com/ganshane/zvirt/protocol"
	"golang.org/x/net/context"
	"github.com/facebookgo/ensure"
)

// DomState implements zvirt_domain.DomState
func (s *ZvirtAgent) DomState(context.Context, *pb.DomStateRequest) (*pb.DomStateResponse, error){
	conn,err := s.pool.Acquire()
	ensure.Nil(s, err)
	defer s.pool.Release(conn)

	return &pb.DomStateResponse{State: pb.DomainState_VIR_DOMAIN_RUNNING}, nil
}
