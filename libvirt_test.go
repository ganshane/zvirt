package zvirt_test

import (
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
	"testing"
	"time"
)

func buildTestConnection() *libvirt.Connect {
	conn, err := libvirt.NewConnect("test:///default")
	if err != nil {
		panic(err)
	}
	return conn
}
func buildTestDomain() (*libvirt.Domain, *libvirt.Connect) {
	conn := buildTestConnection()
	dom, err := conn.DomainDefineXML(`<domain type="test">
		<name>` + time.Now().String() + `</name>
		<memory unit="KiB">8192</memory>
		<os>
			<type>hvm</type>
		</os>
	</domain>`)
	if err != nil {
		panic(err)
	}
	return dom, conn
}
func TestLibVirt_testlist(t *testing.T) {
	_, conn := buildTestDomain()
	if domains, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_RUNNING); err == nil {
		for i, dom := range domains {
			if state, _, err := dom.GetState(); err == nil {
				fmt.Println("state::", state)
			}
			if name, e := dom.GetName(); e == nil {
				fmt.Printf("dom %d %s \n", i, name)
			}

			dom.Free()
		}
	}
}
