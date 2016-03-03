package netaddr_test

import (
	"sort"
	"testing"

	"github.com/hashicorp/go-netaddr"
)

type GoodFixtureInput struct {
	inputNets               []string
	sortedBySpecificMasklen []string
	sortedByBroadMasklen    []string
	sortedByNetwork         []string
}
type GoodFixture struct {
	inputNets               netaddr.NetAddrs
	sortedBySpecificMasklen netaddr.NetAddrs
	sortedByBroadMasklen    netaddr.NetAddrs
	sortedByNetwork         netaddr.NetAddrs
}
type GoodFixtures []*GoodFixture

func makeTestNetAddrs(t *testing.T) GoodFixtures {
	goodFixtureInputs := []GoodFixtureInput{
		{
			inputNets: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedBySpecificMasklen: []string{
				"128.95.120.1/32",
				"192.168.1.10/24",
				"192.168.0.0/16",
				"172.16.1.3/12",
				"10.0.0.0/8",
				"240.0.0.1/4",
			},
			sortedByBroadMasklen: []string{
				"240.0.0.1/4",
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"128.95.120.1/32",
			},
			sortedByNetwork: []string{
				"10.0.0.0/8",
				"128.95.120.1/32",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
		},
	}
	gfs := make(GoodFixtures, 0, len(goodFixtureInputs))
	for _, gfi := range goodFixtureInputs {
		gf := new(GoodFixture)
		gf.inputNets = make(netaddr.NetAddrs, 0, len(gfi.inputNets))
		for _, n := range gfi.inputNets {
			na, err := netaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.inputNets = append(gf.inputNets, na)
		}

		gf.sortedBySpecificMasklen = make(netaddr.NetAddrs, 0, len(gfi.sortedBySpecificMasklen))
		for _, n := range gfi.sortedBySpecificMasklen {
			na, err := netaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sortedBySpecificMasklen = append(gf.sortedBySpecificMasklen, na)
		}

		if len(gf.inputNets) != len(gf.sortedBySpecificMasklen) {
			t.Fatalf("Expected same number of sortedBySpecificMasklen networks")
		}

		gf.sortedByBroadMasklen = make(netaddr.NetAddrs, 0, len(gfi.sortedByBroadMasklen))
		for _, n := range gfi.sortedByBroadMasklen {
			na, err := netaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sortedByBroadMasklen = append(gf.sortedByBroadMasklen, na)
		}

		if len(gf.inputNets) != len(gf.sortedByBroadMasklen) {
			t.Fatalf("Expected same number of sortedByBroadMasklen networks")
		}

		gf.sortedByNetwork = make(netaddr.NetAddrs, 0, len(gfi.sortedByNetwork))
		for _, n := range gfi.sortedByNetwork {
			na, err := netaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sortedByNetwork = append(gf.sortedByNetwork, na)
		}

		if len(gf.inputNets) != len(gf.sortedByNetwork) {
			t.Fatalf("Expected same number of sortedByNetwork networks")
		}
	}

	return gfs
}

func TestNetAddr_Netaddrs_BySpecificMaskLen(t *testing.T) {
	goodFixtures := makeTestNetAddrs(t)

	for _, goodFixture := range goodFixtures {
		inputNets := append(netaddr.NetAddrs(nil), goodFixture.inputNets...)
		sort.Sort(netaddr.BySpecificMaskLen{inputNets})

		var lastLen int = 32
		for i, net := range inputNets {
			maskLen := net.Maskbits()
			if lastLen < maskLen {
				t.Fatalf("Sort by specific mask length failed")
			}
			lastLen = maskLen

			if goodFixture.sortedBySpecificMasklen[i] != net {
				t.Errorf("Expected %s, received %s in iteration %d", goodFixture.sortedBySpecificMasklen[i], net, i)
			}
		}
	}
}

func TestNetAddr_Netaddrs_ByBroadMaskLen(t *testing.T) {
	goodFixtures := makeTestNetAddrs(t)

	for _, goodFixture := range goodFixtures {
		inputNets := append(netaddr.NetAddrs(nil), goodFixture.inputNets...)
		sort.Sort(netaddr.ByBroadMaskLen{inputNets})

		var lastLen int
		for i, net := range inputNets {
			maskLen := net.Maskbits()
			if lastLen > maskLen {
				t.Fatalf("Sort by specific mask length failed")
			}
			lastLen = maskLen

			if goodFixture.sortedByBroadMasklen[i] != net {
				t.Errorf("Expected %s, received %s in iteration %d", goodFixture.sortedByBroadMasklen[i], net, i)
			}
		}
	}
}

func TestNetAddr_Netaddrs_ByNetwork(t *testing.T) {
	goodFixtures := makeTestNetAddrs(t)

	for _, goodFixture := range goodFixtures {
		inputNets := append(netaddr.NetAddrs(nil), goodFixture.inputNets...)
		sort.Sort(netaddr.ByNetwork{inputNets})

		var lastIpUint uint32
		for i, net := range inputNets {
			ipuint, ok := net.ToUint32()
			if !ok {
				t.Errorf("Unable to create a Uin32 from a network")
			}
			if lastIpUint > ipuint {
				t.Fatalf("Sort by network failed")
			}
			lastIpUint = ipuint

			if goodFixture.sortedByNetwork[i] != net {
				t.Errorf("Expected %s, received %s in iteration %d", goodFixture.sortedByNetwork[i], net, i)
			}
		}
	}
}
