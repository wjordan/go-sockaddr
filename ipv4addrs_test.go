package sockaddr_test

import (
	"sort"
	"testing"

	"github.com/hashicorp/go-sockaddr"
)

type GoodTestIPv4AddrsInput struct {
	ipv4Nets                []string
	sortedBySpecificMasklen []string
	sortedByBroadMasklen    []string
	sortedByNetwork         []string
}
type GoodTestIPv4AddrTest struct {
	ipv4Nets                sockaddr.Sockaddrs
	sortedBySpecificMasklen sockaddr.Sockaddrs
	sortedByBroadMasklen    sockaddr.Sockaddrs
	sortedByNetwork         sockaddr.Sockaddrs
}
type GoodTestIPv4AddrTests []*GoodTestIPv4AddrTest

func makeTestIPv4Addrs(t *testing.T) GoodTestIPv4AddrTests {
	goodTestInputs := []GoodTestIPv4AddrsInput{
		{
			ipv4Nets: []string{
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
	gfs := make(GoodTestIPv4AddrTests, 0, len(goodTestInputs))
	for _, gfi := range goodTestInputs {
		gf := new(GoodTestIPv4AddrTest)
		gf.ipv4Nets = make(sockaddr.Sockaddrs, 0, len(gfi.ipv4Nets))
		for _, n := range gfi.ipv4Nets {
			sa, err := sockaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.ipv4Nets = append(gf.ipv4Nets, sa)
		}

		gf.sortedBySpecificMasklen = make(sockaddr.Sockaddrs, 0, len(gfi.sortedBySpecificMasklen))
		for _, n := range gfi.sortedBySpecificMasklen {
			na, err := sockaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sortedBySpecificMasklen = append(gf.sortedBySpecificMasklen, na)
		}

		if len(gf.ipv4Nets) != len(gf.sortedBySpecificMasklen) {
			t.Fatalf("Expected same number of sortedBySpecificMasklen networks")
		}

		gf.sortedByBroadMasklen = make(sockaddr.Sockaddrs, 0, len(gfi.sortedByBroadMasklen))
		for _, n := range gfi.sortedByBroadMasklen {
			na, err := sockaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sortedByBroadMasklen = append(gf.sortedByBroadMasklen, na)
		}

		if len(gf.ipv4Nets) != len(gf.sortedByBroadMasklen) {
			t.Fatalf("Expected same number of sortedByBroadMasklen networks")
		}

		gf.sortedByNetwork = make(sockaddr.Sockaddrs, 0, len(gfi.sortedByNetwork))
		for _, n := range gfi.sortedByNetwork {
			na, err := sockaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sortedByNetwork = append(gf.sortedByNetwork, na)
		}

		if len(gf.ipv4Nets) != len(gf.sortedByNetwork) {
			t.Fatalf("Expected same number of sortedByNetwork networks")
		}
	}

	return gfs
}

type testIPv4AddrsInputs []struct {
	ipv4Nets       []string
	sortedIPv4Nets []string
}

type testIPv4AddrsEntry struct {
	ipv4Nets       sockaddr.Sockaddrs
	sortedIPv4Nets sockaddr.Sockaddrs
}

type testIPv4AddrsTable []testIPv4AddrsEntry

func makeTestsFromInput(t *testing.T, inputs testIPv4AddrsInputs) (tt testIPv4AddrsTable) {
	tt = make(testIPv4AddrsTable, 0, len(inputs))
	for testNum, input := range inputs {
		ipv4Nets := make(sockaddr.Sockaddrs, 0, len(input.ipv4Nets))
		for _, n := range input.ipv4Nets {
			sa, err := sockaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network from %s", n)
			}
			ipv4Nets = append(ipv4Nets, sa)
		}

		sortedIPv4Nets := make(sockaddr.Sockaddrs, 0, len(input.sortedIPv4Nets))
		for _, n := range input.sortedIPv4Nets {
			sa, err := sockaddr.New(n)
			if err != nil {
				t.Fatalf("Expected valid network from %s", n)
			}
			sortedIPv4Nets = append(sortedIPv4Nets, sa)
		}

		if len(ipv4Nets) != len(sortedIPv4Nets) {
			t.Fatalf("Fixtures in row %d/%d different lengths: %d vs %d", testNum, len(inputs), len(ipv4Nets), len(sortedIPv4Nets))
		}
		tt = append(tt, testIPv4AddrsEntry{ipv4Nets, sortedIPv4Nets})
	}

	return tt
}

func TestSockaddr_IPv4Addrs_BySpecificMaskLen(t *testing.T) {
	testInputs := testIPv4AddrsInputs{
		{
			ipv4Nets: []string{"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedIPv4Nets: []string{
				"128.95.120.1/32",
				"192.168.1.10/24",
				"192.168.0.0/16",
				"172.16.1.3/12",
				"10.0.0.0/8",
				"240.0.0.1/4",
			},
		},
	}

	tests := makeTestsFromInput(t, testInputs)

	for _, test := range tests {
		sockaddrs := append(sockaddr.Sockaddrs(nil), test.ipv4Nets...)
		ipaddrs := sockaddrs.FilterByTypeIPAddr()
		sort.Sort(sockaddr.SortIPAddrsBySpecificMaskLen{ipaddrs})

		var lastLen int = 32
		for i, netaddr := range ipaddrs {
			maskLen := netaddr.Maskbits()
			if lastLen < maskLen {
				t.Fatalf("Sort by specific mask length failed")
			}
			lastLen = maskLen

			if test.sortedIPv4Nets[i] != netaddr {
				t.Errorf("Expected %s, received %s in iteration %d", test.sortedIPv4Nets[i], netaddr, i)
			}
		}
	}
}

func TestSockaddr_IPv4Addrs_ByBroadMaskLen(t *testing.T) {
	testInputs := testIPv4AddrsInputs{
		{
			ipv4Nets: []string{"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedIPv4Nets: []string{
				"240.0.0.1/4",
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"128.95.120.1/32",
			},
		},
	}

	tests := makeTestsFromInput(t, testInputs)

	for _, test := range tests {
		sockaddrs := append(sockaddr.Sockaddrs(nil), test.ipv4Nets...)
		ipaddrs := sockaddrs.FilterByTypeIPAddr()
		sort.Sort(sockaddr.SortIPAddrsByBroadMaskLen{ipaddrs})

		var lastLen int
		for i, netaddr := range ipaddrs {
			maskLen := netaddr.Maskbits()
			if lastLen > maskLen {
				t.Fatalf("Sort by specific mask length failed")
			}
			lastLen = maskLen

			if test.sortedIPv4Nets[i] != netaddr {
				t.Errorf("Expected %s, received %s in iteration %d", test.sortedIPv4Nets[i], netaddr, i)
			}
		}
	}
}

func TestSockaddr_IPv4Addrs_IPv4AddrsByNetwork(t *testing.T) {
	testInputs := testIPv4AddrsInputs{
		{
			ipv4Nets: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedIPv4Nets: []string{
				"10.0.0.0/8",
				"128.95.120.1/32",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
		},
	}

	tests := makeTestsFromInput(t, testInputs)

	for _, test := range tests {
		sockaddrs := append(sockaddr.Sockaddrs(nil), test.ipv4Nets...)
		ipv4addrs := sockaddrs.FilterByTypeIPv4Addr()
		sort.Sort(sockaddr.SortIPv4AddrsByNetwork{ipv4addrs})

		var lastIpUint uint32
		for i, netv4addr := range ipv4addrs {
			ipuint := netv4addr.ToUint32()
			if lastIpUint > ipuint {
				t.Fatalf("Sort by network failed")
			}
			lastIpUint = ipuint

			if !netv4addr.Equal(test.sortedIPv4Nets[i]) {
				t.Errorf("[%d] Sort equality failed: expected %s, received %s", i, test.sortedIPv4Nets[i], netv4addr)
			}
		}
	}
}
