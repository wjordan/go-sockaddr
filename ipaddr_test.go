package sockaddr_test

import (
	"testing"

	"github.com/hashicorp/go-sockaddr"
)

func TestSockAddr_IPAddr_CmpPort(t *testing.T) {
	tests := []struct {
		a   string
		b   string
		cmp int
	}{
		{ // 0: Same port, same IPv4Addr
			a:   "208.67.222.222:0",
			b:   "208.67.222.222/32",
			cmp: 0,
		},
		{ // 1: Same port, different IPv4Addr
			a:   "208.67.220.220:0",
			b:   "208.67.222.222/32",
			cmp: 0,
		},
		{ // 2: Same IPv4Addr, different port
			a:   "208.67.222.222:80",
			b:   "208.67.222.222:443",
			cmp: -1,
		},
		{ // 3: Different IPv4Addr, different port
			a:   "208.67.220.220:8600",
			b:   "208.67.222.222:53",
			cmp: 1,
		},
		{ // 4: Same port, same IPv6Addr
			a:   "[::]:0",
			b:   "::/128",
			cmp: 0,
		},
		{ // 5: Same port, different IPv6Addr
			a:   "[::]:0",
			b:   "[2607:f0d0:1002:0051:0000:0000:0000:0004]:0",
			cmp: 0,
		},
		{ // 6: Same IPv6Addr, different port
			a:   "[::]:8400",
			b:   "[::]:8600",
			cmp: -1,
		},
		{ // 7: Different IPv6Addr, different port
			a:   "[::]:8600",
			b:   "[2607:f0d0:1002:0051:0000:0000:0000:0004]:53",
			cmp: 1,
		},
		{ // 8: Mixed IPAddr types, same port
			a:   "[::]:53",
			b:   "208.67.220.220:53",
			cmp: 0,
		},
		{ // 9: Mixed IPAddr types, different port
			a:   "[::]:53",
			b:   "128.95.120.1:123",
			cmp: -1,
		},
	}

	for idx, test := range tests {
		saA, err := sockaddr.NewSockAddr(test.a)
		if err != nil {
			t.Fatalf("[%d] Unable to create a SockAddr from %+q: %v", idx, test.a, err)
		}
		ipA, ok := saA.(sockaddr.IPAddr)
		if !ok {
			t.Fatalf("[%d] Unable to convert SockAddr %+q to an IPAddr", idx, test.a)
		}

		saB, err := sockaddr.NewSockAddr(test.b)
		if err != nil {
			t.Fatalf("[%d] Unable to create an SockAddr from %+q: %v", idx, test.b, err)
		}
		ipB, ok := saB.(sockaddr.IPAddr)
		if !ok {
			t.Fatalf("[%d] Unable to convert SockAddr %+q to an IPAddr", idx, test.b)
		}

		if x := ipA.CmpPort(ipB); x != test.cmp {
			t.Errorf("[%d] IPAddr.CmpPort() failed with %+q with %+q (expected %d, received %d)", idx, ipA, ipB, test.cmp, x)
		}

		if x := ipB.CmpPort(ipA); x*-1 != test.cmp {
			t.Errorf("[%d] IPAddr.CmpPort() failed with %+q with %+q (expected %d, received %d)", idx, ipA, ipB, test.cmp, x)
		}
	}
}
