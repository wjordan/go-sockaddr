package netaddr_test

import (
	"sort"
	"testing"

	"github.com/hashicorp/go-netaddr"
)

func makeTestNetAddrs(t *testing.T) []*netaddr.NetAddr {
	nets := []string{
		"10.0.0.0/8",
		"172.16.1.3/12",
		"192.168.0.0/16",
		"192.168.1.10/24",
		"240.0.0.1/4",
	}
	nas := make([]*netaddr.NetAddr, 0, len(nets))
	for _, n := range nets {
		na, err := netaddr.New(n)
		if err != nil {
			t.Fatalf("Expected valid network")
		}
		nas = append(nas, na)
	}

	return nas
}

func TestNetAddr_Netaddrs_BySpecificMaskLen(t *testing.T) {
	nets := makeTestNetAddrs(t)

	sort.Sort(netaddr.BySpecificMaskLen{nets})

	var lastLen int = 32
	for _, net := range nets {
		maskLen := net.Maskbits()
		if lastLen < maskLen {
			t.Fatalf("Sort by specific mask length failed")
		}
		lastLen = maskLen
	}
}

func TestNetAddr_Netaddrs_ByBroadMaskLen(t *testing.T) {
	nets := makeTestNetAddrs(t)

	sort.Sort(netaddr.ByBroadMaskLen{nets})

	var lastLen int
	for _, net := range nets {
		maskLen := net.Maskbits()
		if lastLen > maskLen {
			t.Fatalf("Sort by specific mask length failed")
		}
		lastLen = maskLen
	}
}

func TestNetAddr_Netaddrs_ByNetwork(t *testing.T) {
	nets := makeTestNetAddrs(t)

	sort.Sort(netaddr.ByNetwork{nets})

	var lastIpUint uint32
	for _, net := range nets {
		ipuint, ok := net.ToUint32()
		if !ok {
			t.Errorf("Unable to create a Uin32 from a network")
		}
		if lastIpUint > ipuint {
			t.Fatalf("Sort by network failed")
		}
		lastIpUint = ipuint
	}
}
