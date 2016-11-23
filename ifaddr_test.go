package sockaddr_test

import (
	"net"
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

func TestGetPublicIP(t *testing.T) {
	ip, err := sockaddr.GetPublicIP()
	if err != nil {
		t.Fatalf("unable to get a public IP")
	}

	if ip == "" {
		t.Skip("it's hard to test this reliably")
	}
}

func TestAscIfType(t *testing.T) {
	tests := []struct {
		name       string
		t1         sockaddr.IfAddr // must come before t2 according to the ascOp
		t2         sockaddr.IfAddr
		ascOp      sockaddr.CmpIfAddrFunc
		ascResult  int
		descOp     sockaddr.CmpIfAddrFunc
		descResult int
	}{
		{
			name:       "empty test",
			t1:         sockaddr.IfAddr{},
			t2:         sockaddr.IfAddr{},
			ascOp:      sockaddr.AscIfAddress,
			descOp:     sockaddr.DescIfAddress,
			ascResult:  0,
			descResult: 0,
		},
		{
			name: "ipv4 address less",
			t1: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("1.2.3.3"),
			},
			t2: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("1.2.3.4"),
			},
			ascOp:      sockaddr.AscIfAddress,
			descOp:     sockaddr.DescIfAddress,
			ascResult:  -1,
			descResult: -1,
		},
		{
			name: "ipv4 private",
			t1: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("10.1.2.3"),
			},
			t2: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("203.0.113.3"),
			},
			ascOp:      sockaddr.AscIfPrivate,
			descOp:     sockaddr.DescIfPrivate,
			ascResult:  0, // not both private, can't complete the test
			descResult: 0,
		},
		{
			name: "IfAddr name",
			t1: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("10.1.2.3"),
				Interface: net.Interface{
					Name: "abc0",
				},
			},
			t2: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("203.0.113.3"),
				Interface: net.Interface{
					Name: "xyz0",
				},
			},
			ascOp:      sockaddr.AscIfName,
			descOp:     sockaddr.DescIfName,
			ascResult:  -1,
			descResult: -1,
		},
		{
			name: "IfAddr network size",
			t1: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("10.0.0.0/8"),
			},
			t2: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("127.0.0.0/24"),
			},
			ascOp:      sockaddr.AscIfNetworkSize,
			descOp:     sockaddr.DescIfNetworkSize,
			ascResult:  -1,
			descResult: -1,
		},
		{
			name: "IfAddr port",
			t1: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("10.0.0.0:80"),
			},
			t2: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("127.0.0.0:8600"),
			},
			ascOp:      sockaddr.AscIfPort,
			descOp:     sockaddr.DescIfPort,
			ascResult:  -1,
			descResult: -1,
		},
		{
			name: "IfAddr type",
			t1: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv4Addr("10.0.0.0:80"),
			},
			t2: sockaddr.IfAddr{
				SockAddr: sockaddr.MustIPv6Addr("[::1]:80"),
			},
			ascOp:      sockaddr.AscIfType,
			descOp:     sockaddr.DescIfType,
			ascResult:  -1,
			descResult: -1,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d must have a name", i)
		}

		// Test ascending operation
		ascExpected := test.ascResult
		ascResult := test.ascOp(&test.t1, &test.t2)
		if ascResult != ascExpected {
			t.Errorf("%s: Unexpected result %d, expected %d when comparing %v and %v using %v", test.name, ascResult, ascExpected, test.t1, test.t2, test.ascOp)
		}

		// Test descending operation
		descExpected := test.descResult
		descResult := test.descOp(&test.t2, &test.t1)
		if descResult != descExpected {
			t.Errorf("%s: Unexpected result %d, expected %d when comparing %v and %v using %v", test.name, descResult, descExpected, test.t1, test.t2, test.descOp)
		}

		if ascResult != descResult {
			t.Fatalf("bad")
		}

		// Reverse the args
		ascExpected = -1 * test.ascResult
		ascResult = test.ascOp(&test.t2, &test.t1)
		if ascResult != ascExpected {
			t.Errorf("%s: Unexpected result %d, expected %d when comparing %v and %v using %v", test.name, ascResult, ascExpected, test.t1, test.t2, test.ascOp)
		}

		descExpected = -1 * test.descResult
		descResult = test.descOp(&test.t1, &test.t2)
		if descResult != descExpected {
			t.Errorf("%s: Unexpected result %d, expected %d when comparing %v and %v using %v", test.name, descResult, descExpected, test.t1, test.t2, test.descOp)
		}

		if ascResult != descResult {
			t.Fatalf("bad")
		}

		// Test equality
		ascExpected = 0
		ascResult = test.ascOp(&test.t1, &test.t1)
		if ascResult != ascExpected {
			t.Errorf("%s: Unexpected result %d, expected %d when comparing %v and %v using %v", test.name, ascResult, ascExpected, test.t1, test.t2, test.ascOp)
		}

		descExpected = 0
		descResult = test.descOp(&test.t1, &test.t1)
		if descResult != descExpected {
			t.Errorf("%s: Unexpected result %d, expected %d when comparing %v and %v using %v", test.name, descResult, descExpected, test.t1, test.t2, test.descOp)
		}
	}
}

func TestIfAddrAttr(t *testing.T) {
	tests := []struct {
		name     string
		ifAddr   sockaddr.IfAddr
		attr     string
		expected string
	}{
		{
			name: "name",
			ifAddr: sockaddr.IfAddr{
				Interface: net.Interface{
					Name: "abc0",
				},
			},
			attr:     "name",
			expected: "abc0",
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d must have a name", i)
		}

		result, err := sockaddr.IfAttr(test.attr, sockaddr.IfAddrs{test.ifAddr})
		if err != nil {
			t.Errorf("failed to get attr %q from %v", test.name, test.ifAddr)
		}

		if result != test.expected {
			t.Errorf("unexpected result")
		}
	}

	// Test an empty array
	result, err := sockaddr.IfAttr("name", sockaddr.IfAddrs{})
	if err != nil {
		t.Error(`failed to get attr "name" from an empty array`)
	}

	if result != "" {
		t.Errorf("unexpected result")
	}
}
