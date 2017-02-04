package zvirt

import (
	pb "github.com/ganshane/zvirt/protocol"
	"golang.org/x/net/context"
)

// DomState implements zvirt_domain.DomState
func (s *ZvirtAgent) DomState(context.Context, *pb.DomStateRequest) (*pb.DomStateResponse, error){
	return &pb.DomStateResponse{State: pb.DomainState_VIR_DOMAIN_RUNNING}, nil
}
