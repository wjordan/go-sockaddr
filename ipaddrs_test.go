package sockaddr_test

import (
	"sort"
	"testing"

	"github.com/hashicorp/go-sockaddr"
)

type GoodTestIPAddrTest struct {
	sockAddrs               sockaddr.SockAddrs
	sortedBySpecificMasklen sockaddr.SockAddrs
	sortedByBroadMasklen    sockaddr.SockAddrs
	sortedByNetwork         sockaddr.SockAddrs
}
type GoodTestIPAddrTests []*GoodTestIPAddrTest

func makeTestIPAddrs(t *testing.T) GoodTestIPAddrTests {
	goodTestInputs := []struct {
		sockAddrs               []string
		sortedBySpecificMasklen []string
		sortedByBroadMasklen    []string
		sortedByNetwork         []string
	}{
		{
			sockAddrs: []string{
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
	gfs := make(GoodTestIPAddrTests, 0, len(goodTestInputs))
	for _, gfi := range goodTestInputs {
		gf := new(GoodTestIPAddrTest)
		gf.sockAddrs = make(sockaddr.SockAddrs, 0, len(gfi.sockAddrs))
		for _, n := range gfi.sockAddrs {
			sa, err := sockaddr.NewSockAddr(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sockAddrs = append(gf.sockAddrs, sa)
		}

		gf.sortedBySpecificMasklen = make(sockaddr.SockAddrs, 0, len(gfi.sortedBySpecificMasklen))
		for _, n := range gfi.sortedBySpecificMasklen {
			na, err := sockaddr.NewSockAddr(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sortedBySpecificMasklen = append(gf.sortedBySpecificMasklen, na)
		}

		if len(gf.sockAddrs) != len(gf.sortedBySpecificMasklen) {
			t.Fatalf("Expected same number of sortedBySpecificMasklen networks")
		}

		gf.sortedByBroadMasklen = make(sockaddr.SockAddrs, 0, len(gfi.sortedByBroadMasklen))
		for _, n := range gfi.sortedByBroadMasklen {
			na, err := sockaddr.NewSockAddr(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sortedByBroadMasklen = append(gf.sortedByBroadMasklen, na)
		}

		if len(gf.sockAddrs) != len(gf.sortedByBroadMasklen) {
			t.Fatalf("Expected same number of sortedByBroadMasklen networks")
		}

		gf.sortedByNetwork = make(sockaddr.SockAddrs, 0, len(gfi.sortedByNetwork))
		for _, n := range gfi.sortedByNetwork {
			na, err := sockaddr.NewSockAddr(n)
			if err != nil {
				t.Fatalf("Expected valid network")
			}
			gf.sortedByNetwork = append(gf.sortedByNetwork, na)
		}

		if len(gf.sockAddrs) != len(gf.sortedByNetwork) {
			t.Fatalf("Expected same number of sortedByNetwork networks")
		}
	}

	return gfs
}

func TestSockAddr_IPAddrs_BySpecificMaskLen(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"128.95.120.1/32",
				"192.168.1.10/24",
				"192.168.0.0/16",
				"172.16.1.3/12",
				"10.0.0.0/8",
				"240.0.0.1/4",
			},
		},
	}

	for _, test := range testInputs {
		inputAddrs := convertToSockAddrs(t, test.inputAddrs)
		sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)
		sockaddrs := append(sockaddr.SockAddrs(nil), inputAddrs...)
		filteredAddrs := sockaddrs.FilterByType(sockaddr.TypeIPv4)
		ipv4Addrs := make([]sockaddr.IPv4Addr, 0, len(filteredAddrs))
		for _, x := range filteredAddrs {
			switch v := x.(type) {
			case sockaddr.IPv4Addr:
				ipv4Addrs = append(ipv4Addrs, v)
			default:
				t.Fatalf("invalid type")
			}
		}

		ipAddrs := make([]sockaddr.IPAddr, 0, len(filteredAddrs))
		for _, x := range filteredAddrs {
			ipAddr, ok := x.(sockaddr.IPAddr)
			if !ok {
				t.Fatalf("Unable to typecast to IPAddr")
			}
			ipAddrs = append(ipAddrs, ipAddr)
		}
		sort.Sort(sockaddr.SortIPAddrsBySpecificMaskLen{ipAddrs})

		var lastLen int = 32
		for i, netaddr := range ipAddrs {
			maskLen := netaddr.Maskbits()
			if lastLen < maskLen {
				t.Fatalf("Sort by specific mask length failed")
			}
			lastLen = maskLen

			if sortedAddrs[i] != netaddr {
				t.Errorf("Expected %s, received %s in iteration %d", sortedAddrs[i], netaddr, i)
			}
		}
	}
}

func TestSockAddr_IPAddrs_ByBroadMaskLen(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"240.0.0.1/4",
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"128.95.120.1/32",
			},
		},
	}

	for _, test := range testInputs {
		inputAddrs := convertToSockAddrs(t, test.inputAddrs)
		sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)
		sockaddrs := append(sockaddr.SockAddrs(nil), inputAddrs...)
		filteredAddrs := sockaddrs.FilterByType(sockaddr.TypeIP)
		ipAddrs := make([]sockaddr.IPAddr, 0, len(filteredAddrs))
		for _, x := range filteredAddrs {
			ipAddr, ok := x.(sockaddr.IPAddr)
			if !ok {
				t.Fatalf("Unable to typecast to IPAddr")
			}
			ipAddrs = append(ipAddrs, ipAddr)
		}
		sort.Sort(sockaddr.SortIPAddrsByBroadMaskLen{ipAddrs})

		var lastLen int
		for i, netaddr := range ipAddrs {
			maskLen := netaddr.Maskbits()
			if lastLen > maskLen {
				t.Fatalf("Sort by specific mask length failed")
			}
			lastLen = maskLen

			if sortedAddrs[i] != netaddr {
				t.Errorf("Expected %s, received %s in iteration %d", sortedAddrs[i], netaddr, i)
			}
		}
	}
}

func TestSockAddr_IPAddrs_IPAddrsByNetwork(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"10.0.0.0/8",
				"128.95.120.1/32",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
		},
	}

	for _, test := range testInputs {
		inputAddrs := convertToSockAddrs(t, test.inputAddrs)
		sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)
		sockaddrs := append(sockaddr.SockAddrs(nil), inputAddrs...)
		ipaddrs := sockaddrs.FilterByTypeIPv4Addr()
		sort.Sort(sockaddr.SortIPAddrsByNetwork{ipaddrs})

		var lastIpUint sockaddr.IPv4Address
		for i, netaddr := range ipaddrs {
			if lastIpUint > netaddr.Address {
				t.Fatalf("Sort by network failed")
			}
			lastIpUint = netaddr.Address

			if !netaddr.Equal(sortedAddrs[i]) {
				t.Errorf("[%d] Sort equality failed: expected %s, received %s", i, sortedAddrs[i], netaddr)
			}
		}
	}
}

func TestSockAddr_IPAddrs_IPAddrsByNetworkSize(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"128.95.120.2:53",
				"128.95.120.2/32",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"128.95.120.2:8600",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"128.95.120.1/32",
				"128.95.120.2:53",
				"128.95.120.2:8600",
				"128.95.120.2/32",
				"192.168.1.10/24",
				"192.168.0.0/16",
				"172.16.1.3/12",
				"10.0.0.0/8",
				"240.0.0.1/4",
			},
		},
	}

	for _, test := range testInputs {
		inputAddrs := convertToSockAddrs(t, test.inputAddrs)
		sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)

		sockaddrs := append(sockaddr.SockAddrs(nil), inputAddrs...)
		filteredAddrs := sockaddrs.FilterByType(sockaddr.TypeIP)
		ipAddrs := make([]sockaddr.IPAddr, 0, len(filteredAddrs))
		for _, x := range filteredAddrs {
			ipAddr, ok := x.(sockaddr.IPAddr)
			if !ok {
				t.Fatalf("Unable to typecast to IPAddr")
			}
			ipAddrs = append(ipAddrs, ipAddr)
		}
		sort.Sort(sockaddr.SortIPAddrsByNetworkSize{ipAddrs})

		// var prevAddr sockaddr.IPAddr
		for i, ipAddr := range ipAddrs {
			// if i == 0 {
			// 	prevAddr = ipAddr
			// 	continue
			// }

			// if prevAddr.Cmp(ipAddr) > 0 {
			// 	t.Logf("[%d] Prev:\t%v", i, prevAddr)
			// 	t.Logf("[%d] ipAddr:\t%v", i, ipAddr)
			// 	t.Fatalf("Sort by network failed")
			// }
			// prevAddr = ipAddr

			if !ipAddr.Equal(sortedAddrs[i]) {
				t.Errorf("[%d] Sort equality failed: expected %s, received %s", i, sortedAddrs[i], ipAddr)
			}
		}
	}
}

// func TestSockAddr_IPAddrs_IPAddrsByCmp(t *testing.T) {
// 	testInputs := testIPAddrsInputs{
// 		{
// 			sockAddrs: []string{
// 				"10.0.0.0/8",
// 				"172.16.1.3/12",
// 				"128.95.120.2:53",
// 				"128.95.120.2/32",
// 				"192.168.0.0/16",
// 				"128.95.120.1/32",
// 				"192.168.1.10/24",
// 				"128.95.120.2:8600",
// 				"240.0.0.1/4",
// 			},
// 			sortedSockAddrs: []string{
// 				"128.95.120.1/32",
// 				"128.95.120.2:53",
// 				"128.95.120.2:8600",
// 				"128.95.120.2/32",
// 				"192.168.1.10/24",
// 				"192.168.0.0/16",
// 				"172.16.1.3/12",
// 				"10.0.0.0/8",
// 				"240.0.0.1/4",
// 			},
// 		},
// 	}

// 	for _, test := range makeTestsFromInput(t, testInputs) {
// 		sockaddrs := append(sockaddr.SockAddrs(nil), test.sockAddrs...)
// 		ipAddrs := sockaddrs.FilterByTypeIPAddr()
// 		sort.Sort(sockaddr.SortIPAddrsByCmp{ipAddrs})
// 		t.Logf("Here: %+v", ipAddrs)

// 		var prevAddr sockaddr.IPAddr
// 		for i, ipAddr := range ipAddrs {
// 			if i == 0 {
// 				prevAddr = ipAddr
// 				continue
// 			}

// 			if prevAddr.Cmp(ipAddr) > 0 {
// 				t.Logf("[%d] Prev:\t%v", i, prevAddr)
// 				t.Logf("[%d] ipAddr:\t%v", i, ipAddr)
// 				t.Fatalf("Sort by network failed")
// 			}
// 			prevAddr = ipAddr

// 			if !ipAddr.Equal(test.sortedSockAddrs[i]) {
// 				t.Errorf("[%d] Sort equality failed: expected %s, received %s", i, test.sortedSockAddrs[i], ipAddr)
// 			}
// 		}
// 	}
// }

func TestSockAddr_IPAddrs_IPAddrsByCmp(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"128.95.120.2:53",
				"128.95.120.2/32",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"128.95.120.2:8600",
				"0:0:0:0:0:0:0:0",
				"0:0:0:0:0:0:0:1",
				"2607:f0d0:1002:0051:0000:0000:0000:0004",
				"2607:f0d0:1002:0051:0000:0000:0000:0003",
				"2607:f0d0:1002:0051:0000:0000:0000:0005",
				"[2607:f0d0:1002:0051:0000:0000:0000:0004]:8600",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"240.0.0.1/4",
				"128.95.120.1/32",
				"128.95.120.2/32",
				"128.95.120.2:53",
				"128.95.120.2:8600",
				"0:0:0:0:0:0:0:0",
				"0:0:0:0:0:0:0:1",
				"2607:f0d0:1002:0051:0000:0000:0000:0003",
				"2607:f0d0:1002:0051:0000:0000:0000:0004",
				"[2607:f0d0:1002:0051:0000:0000:0000:0004]:8600",
				"2607:f0d0:1002:0051:0000:0000:0000:0005",
			},
		},
	}

	for _, test := range testInputs {
		shuffleStrings(test.inputAddrs)

		inputAddrs := convertToSockAddrs(t, test.inputAddrs)
		sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)

		sockaddr.OrderedBy(sockaddr.AscType, sockaddr.AscPrivate, sockaddr.AscAddress, sockaddr.AscPort).Sort(inputAddrs)

		for i, sockAddr := range inputAddrs {
			if !sockAddr.Equal(sortedAddrs[i]) {
				t.Errorf("[%d] Sort equality failed: expected %s, received %s", i, sortedAddrs[i], sockAddr)
			}
		}
	}
}
