package zvirt
import (
	libvirt "github.com/libvirt/libvirt-go"
)

//global ZvirtAgent server
type ZvirtAgent struct {
	conn *libvirt.Connect //libvirt connection
}
//create new server for zvirt
func NewServer(conn *libvirt.Connect) *ZvirtAgent {
	return &ZvirtAgent{conn:conn}
}
