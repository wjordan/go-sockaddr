package sockaddr_test

import (
	"math/rand"
	"testing"

	"github.com/hashicorp/consul/lib"
	"github.com/hashicorp/go-sockaddr"
)

func init() {
	lib.SeedMathRand()
}

// sockAddrStringInputs allows for easy test creation by developers.
// Parallel arrays of string inputs are converted to their SockAddr
// equivalents for use by unit tests.
type sockAddrStringInputs []struct {
	inputAddrs    []string
	sortedAddrs   []string
	sortedTypes   []sockaddr.SockAddrType
	sortFuncs     []sockaddr.CmpFunc
	numIPv4Inputs int
	numIPv6Inputs int
	numUnixInputs int
}

func convertToSockAddrs(t *testing.T, inputs []string) sockaddr.SockAddrs {
	sockAddrs := make(sockaddr.SockAddrs, 0, len(inputs))
	for i, input := range inputs {
		sa, err := sockaddr.NewSockAddr(input)
		if err != nil {
			t.Fatalf("[%d] Invalid SockAddr input for %+q: %v", i, input, err)
		}
		sockAddrs = append(sockAddrs, sa)
	}

	return sockAddrs
}

// shuffleStrings randomly shuffles the list of strings
func shuffleStrings(list []string) {
	for i := range list {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
}

func TestSockAddr_SockAddrs_AscAddress(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{ // testNum: 0
			sortFuncs: []sockaddr.CmpFunc{
				sockaddr.AscAddress,
			},
			numIPv4Inputs: 9,
			numIPv6Inputs: 1,
			numUnixInputs: 0,
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
				"::",
			},
			sortedAddrs: []string{
				"10.0.0.0/8",
				"128.95.120.1/32",
				"128.95.120.2:53",
				"128.95.120.2/32",
				"128.95.120.2:8600",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"240.0.0.1/4",
				"::",
			},
		},
	}

	for testNum, test := range testInputs {
		shuffleStrings(test.inputAddrs)
		inputSockAddrs := convertToSockAddrs(t, test.inputAddrs)
		sortedIPv4Addrs, nonIPv4Addrs := convertToSockAddrs(t, test.sortedAddrs).OnlyIPv4()
		if l := len(sortedIPv4Addrs); l != test.numIPv4Inputs {
			t.Fatal("[%d] Missing IPv4Addrs: expected %d, received %d", testNum, test.numIPv4Inputs, l)
		}
		if len(nonIPv4Addrs) != test.numIPv6Inputs+test.numUnixInputs {
			t.Fatal("[%d] Non-IPv4 Address in input", testNum)
		}

		// Copy inputAddrs so we can manipulate it. wtb const.
		sockAddrs := append(sockaddr.SockAddrs(nil), inputSockAddrs...)
		filteredAddrs := sockAddrs.FilterByType(sockaddr.TypeIPv4)
		sockaddr.OrderedBy(test.sortFuncs...).Sort(filteredAddrs)
		ipv4Addrs, nonIPv4s := filteredAddrs.OnlyIPv4()
		if len(nonIPv4s) != 0 {
			t.Fatalf("[%d] bad", testNum)
		}

		for i, ipv4Addr := range ipv4Addrs {
			if ipv4Addr.Address != sortedIPv4Addrs[i].Address {
				t.Errorf("[%d/%d] Sort equality failed: expected %s, received %s", testNum, i, sortedIPv4Addrs[i], ipv4Addr)
			}
		}
	}
}

func TestSockAddr_SockAddrs_AscPrivate(t *testing.T) {
	testInputs := []struct {
		sortFuncs   []sockaddr.CmpFunc
		inputAddrs  []string
		sortedAddrs []string
	}{
		{ // testNum: 0
			sortFuncs: []sockaddr.CmpFunc{
				sockaddr.AscPrivate,
				sockaddr.AscAddress,
				sockaddr.AscType,
				sockaddr.AscAddress,
				sockaddr.AscPort,
			},
			inputAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"128.95.120.1/32",
				"128.95.120.2/32",
				"128.95.120.2:53",
				"128.95.120.2:8600",
				"240.0.0.1/4",
				"::",
			},
			sortedAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"128.95.120.1/32",
				"128.95.120.2/32",
				// "128.95.120.2:53",
				// "128.95.120.2:8600",
				// "240.0.0.1/4",
				// "::",
			},
		},
	}

	for testNum, test := range testInputs {
		sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)

		inputAddrs := append([]string(nil), test.inputAddrs...)
		shuffleStrings(inputAddrs)
		inputSockAddrs := convertToSockAddrs(t, inputAddrs)

		sockaddr.OrderedBy(test.sortFuncs...).Sort(inputSockAddrs)

		for i, sockAddr := range sortedAddrs {
			if !sockAddr.Equal(inputSockAddrs[i]) {
				t.Logf("Input Addrs:\t%+v", inputAddrs)
				t.Logf("Sorted Addrs:\t%+v", inputSockAddrs)
				t.Logf("Expected Addrs:\t%+v", test.sortedAddrs)
				t.Fatalf("[%d/%d] Sort AscType/AscAddress failed: expected %+q, received %+q", testNum, i, sockAddr, inputSockAddrs[i])
			}
		}
	}
}

func TestSockAddr_SockAddrs_AscType(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{ // testNum: 0
			sortFuncs: []sockaddr.CmpFunc{
				sockaddr.AscType,
			},
			inputAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"128.95.120.2:53",
				"::",
				"128.95.120.2/32",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"128.95.120.2:8600",
				"240.0.0.1/4",
			},
			sortedTypes: []sockaddr.SockAddrType{
				sockaddr.TypeIPv4,
				sockaddr.TypeIPv4,
				sockaddr.TypeIPv4,
				sockaddr.TypeIPv4,
				sockaddr.TypeIPv4,
				sockaddr.TypeIPv4,
				sockaddr.TypeIPv4,
				sockaddr.TypeIPv4,
				sockaddr.TypeIPv4,
				sockaddr.TypeIPv6,
			},
		},
	}

	for testNum, test := range testInputs {
		shuffleStrings(test.inputAddrs)

		inputSockAddrs := convertToSockAddrs(t, test.inputAddrs)
		sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)

		sockaddr.OrderedBy(test.sortFuncs...).Sort(inputSockAddrs)

		for i, sockAddr := range sortedAddrs {
			if sockAddr.Type() != sortedAddrs[i].Type() {
				t.Errorf("[%d/%d] Sort AscType failed: expected %+q, received %+q", testNum, i, sortedAddrs[i], sockAddr)
			}
		}
	}
}
