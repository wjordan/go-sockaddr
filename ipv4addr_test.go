package sockaddr_test

import (
	"testing"

	"github.com/hashicorp/go-sockaddr"
)

func TestSockAddr_IPv4Addr(t *testing.T) {
	tests := []struct {
		input         string
		address       string
		uintAddr      sockaddr.IPv4Address
		uintNet       sockaddr.IPv4Address
		uintMask      sockaddr.IPv4Mask
		octets        []int
		broadcast     string
		firstUsable   string
		lastUsable    string
		listenTCPArgs []string
		listenUDPArgs []string
		ipMaskStr     string
		ipNetStr      string
		networkStr    string
		addrStr       string
		port          sockaddr.IPPort
		maskbits      int
		addrBinStr    string
		addrHexStr    string
		pass          bool
	}{
		{ // 0
			input:         "0.0.0.0",
			uintAddr:      0,
			uintNet:       0,
			uintMask:      sockaddr.IPv4HostMask,
			octets:        []int{0, 0, 0, 0},
			address:       "0.0.0.0",
			broadcast:     "0.0.0.0",
			firstUsable:   "0.0.0.0",
			lastUsable:    "0.0.0.0",
			listenTCPArgs: []string{"tcp4", "0.0.0.0:0"},
			listenUDPArgs: []string{"udp4", "0.0.0.0:0"},
			ipMaskStr:     "ffffffff",
			ipNetStr:      "0.0.0.0/32",
			networkStr:    "0.0.0.0",
			addrStr:       "0.0.0.0",
			maskbits:      32,
			addrBinStr:    "00000000000000000000000000000000",
			addrHexStr:    "00000000",
			pass:          true,
		},
		{ // 1
			input:         "0.0.0.0:80",
			uintAddr:      0,
			uintNet:       0,
			uintMask:      sockaddr.IPv4HostMask,
			octets:        []int{0, 0, 0, 0},
			address:       "0.0.0.0",
			broadcast:     "0.0.0.0",
			firstUsable:   "0.0.0.0",
			lastUsable:    "0.0.0.0",
			listenTCPArgs: []string{"tcp4", "0.0.0.0:80"},
			listenUDPArgs: []string{"udp4", "0.0.0.0:80"},
			ipMaskStr:     "ffffffff",
			ipNetStr:      "0.0.0.0/32",
			networkStr:    "0.0.0.0",
			addrStr:       "0.0.0.0:80",
			port:          80,
			maskbits:      32,
			addrBinStr:    "00000000000000000000000000000000",
			addrHexStr:    "00000000",
			pass:          true,
		},
		{ // 2
			input:         "0.0.0.0/0",
			uintAddr:      0,
			uintNet:       0,
			uintMask:      0,
			octets:        []int{0, 0, 0, 0},
			address:       "0.0.0.0",
			broadcast:     "255.255.255.255",
			firstUsable:   "0.0.0.1",
			lastUsable:    "255.255.255.254",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "00000000",
			ipNetStr:      "0.0.0.0/0",
			networkStr:    "0.0.0.0/0",
			addrStr:       "0.0.0.0/0",
			addrBinStr:    "00000000000000000000000000000000",
			addrHexStr:    "00000000",
			pass:          true,
		},
		{ // 3
			input:         "0.0.0.1",
			uintAddr:      1,
			uintNet:       1,
			uintMask:      sockaddr.IPv4HostMask,
			octets:        []int{0, 0, 0, 1},
			address:       "0.0.0.1",
			broadcast:     "0.0.0.1",
			firstUsable:   "0.0.0.1",
			lastUsable:    "0.0.0.1",
			listenTCPArgs: []string{"tcp4", "0.0.0.1:0"},
			listenUDPArgs: []string{"udp4", "0.0.0.1:0"},
			ipMaskStr:     "ffffffff",
			ipNetStr:      "0.0.0.1/32",
			networkStr:    "0.0.0.1",
			addrStr:       "0.0.0.1",
			maskbits:      32,
			addrBinStr:    "00000000000000000000000000000001",
			addrHexStr:    "00000001",
			pass:          true,
		},
		{ // 4
			input:         "0.0.0.1/1",
			uintAddr:      1,
			uintNet:       0,
			uintMask:      2147483648,
			octets:        []int{0, 0, 0, 1},
			address:       "0.0.0.1",
			broadcast:     "127.255.255.255",
			firstUsable:   "0.0.0.1",
			lastUsable:    "127.255.255.254",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "80000000",
			ipNetStr:      "0.0.0.0/1",
			networkStr:    "0.0.0.0/1",
			addrStr:       "0.0.0.1/1",
			maskbits:      1,
			addrBinStr:    "00000000000000000000000000000001",
			addrHexStr:    "00000001",
			pass:          true,
		},
		{ // 5
			input:         "1.2.3.4",
			uintAddr:      16909060,
			uintNet:       16909060,
			uintMask:      sockaddr.IPv4HostMask,
			octets:        []int{1, 2, 3, 4},
			address:       "1.2.3.4",
			broadcast:     "1.2.3.4",
			firstUsable:   "1.2.3.4",
			lastUsable:    "1.2.3.4",
			listenTCPArgs: []string{"tcp4", "1.2.3.4:0"},
			listenUDPArgs: []string{"udp4", "1.2.3.4:0"},
			ipMaskStr:     "ffffffff",
			ipNetStr:      "1.2.3.4/32",
			networkStr:    "1.2.3.4",
			addrStr:       "1.2.3.4",
			maskbits:      32,
			addrBinStr:    "00000001000000100000001100000100",
			addrHexStr:    "01020304",
			pass:          true,
		},
		{ // 6
			input:         "10.0.0.0/8",
			uintAddr:      167772160,
			uintNet:       167772160,
			uintMask:      4278190080,
			octets:        []int{10, 0, 0, 0},
			address:       "10.0.0.0",
			broadcast:     "10.255.255.255",
			firstUsable:   "10.0.0.1",
			lastUsable:    "10.255.255.254",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "ff000000",
			ipNetStr:      "10.0.0.0/8",
			networkStr:    "10.0.0.0/8",
			addrStr:       "10.0.0.0/8",
			maskbits:      8,
			addrBinStr:    "00001010000000000000000000000000",
			addrHexStr:    "0a000000",
			pass:          true,
		},
		{ // 7
			input:         "128.0.0.0",
			uintAddr:      2147483648,
			uintNet:       2147483648,
			uintMask:      sockaddr.IPv4HostMask,
			octets:        []int{128, 0, 0, 0},
			address:       "128.0.0.0",
			broadcast:     "128.0.0.0",
			firstUsable:   "128.0.0.0",
			lastUsable:    "128.0.0.0",
			listenTCPArgs: []string{"tcp4", "128.0.0.0:0"},
			listenUDPArgs: []string{"udp4", "128.0.0.0:0"},
			ipMaskStr:     "ffffffff",
			ipNetStr:      "128.0.0.0/32",
			networkStr:    "128.0.0.0",
			addrStr:       "128.0.0.0",
			maskbits:      32,
			addrBinStr:    "10000000000000000000000000000000",
			addrHexStr:    "80000000",
			pass:          true,
		},
		{ // 8
			input:         "128.95.120.1/32",
			uintAddr:      2153740289,
			uintNet:       2153740289,
			uintMask:      sockaddr.IPv4HostMask,
			octets:        []int{128, 95, 120, 1},
			address:       "128.95.120.1",
			broadcast:     "128.95.120.1",
			firstUsable:   "128.95.120.1",
			lastUsable:    "128.95.120.1",
			listenTCPArgs: []string{"tcp4", "128.95.120.1:0"},
			listenUDPArgs: []string{"udp4", "128.95.120.1:0"},
			ipMaskStr:     "ffffffff",
			ipNetStr:      "128.95.120.1/32",
			networkStr:    "128.95.120.1",
			addrStr:       "128.95.120.1",
			maskbits:      32,
			addrBinStr:    "10000000010111110111100000000001",
			addrHexStr:    "805f7801",
			pass:          true,
		},
		{ // 9
			input:         "172.16.1.3/12",
			uintAddr:      2886729987,
			uintNet:       2886729728,
			uintMask:      4293918720,
			octets:        []int{172, 16, 1, 3},
			address:       "172.16.1.3",
			broadcast:     "172.31.255.255",
			firstUsable:   "172.16.0.1",
			lastUsable:    "172.31.255.254",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "fff00000",
			ipNetStr:      "172.16.0.0/12",
			networkStr:    "172.16.0.0/12",
			addrStr:       "172.16.1.3/12",
			maskbits:      12,
			addrBinStr:    "10101100000100000000000100000011",
			addrHexStr:    "ac100103",
			pass:          true,
		},
		{ // 10
			input:         "192.168.0.0/16",
			uintAddr:      3232235520,
			uintNet:       3232235520,
			uintMask:      4294901760,
			octets:        []int{192, 168, 0, 0},
			address:       "192.168.0.0",
			broadcast:     "192.168.255.255",
			firstUsable:   "192.168.0.1",
			lastUsable:    "192.168.255.254",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "ffff0000",
			ipNetStr:      "192.168.0.0/16",
			networkStr:    "192.168.0.0/16",
			addrStr:       "192.168.0.0/16",
			maskbits:      16,
			addrBinStr:    "11000000101010000000000000000000",
			addrHexStr:    "c0a80000",
			pass:          true,
		},
		{ // 11
			input:         "192.168.0.1",
			uintAddr:      3232235521,
			uintNet:       3232235521,
			uintMask:      sockaddr.IPv4HostMask,
			octets:        []int{192, 168, 0, 1},
			address:       "192.168.0.1",
			broadcast:     "192.168.0.1",
			firstUsable:   "192.168.0.1",
			lastUsable:    "192.168.0.1",
			listenTCPArgs: []string{"tcp4", "192.168.0.1:0"},
			listenUDPArgs: []string{"udp4", "192.168.0.1:0"},
			ipMaskStr:     "ffffffff",
			ipNetStr:      "192.168.0.1/32",
			networkStr:    "192.168.0.1",
			addrStr:       "192.168.0.1",
			maskbits:      32,
			addrBinStr:    "11000000101010000000000000000001",
			addrHexStr:    "c0a80001",
			pass:          true,
		},
		{ // 12
			input:         "192.168.0.2/31",
			uintAddr:      3232235522,
			uintNet:       3232235522,
			uintMask:      4294967294,
			octets:        []int{192, 168, 0, 2},
			address:       "192.168.0.2",
			broadcast:     "192.168.0.3",
			firstUsable:   "192.168.0.2",
			lastUsable:    "192.168.0.3",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "fffffffe",
			ipNetStr:      "192.168.0.2/31",
			networkStr:    "192.168.0.2/31",
			addrStr:       "192.168.0.2/31",
			maskbits:      31,
			addrBinStr:    "11000000101010000000000000000010",
			addrHexStr:    "c0a80002",
			pass:          true,
		},
		{ // 13
			input:         "192.168.1.10/24",
			uintAddr:      3232235786,
			uintNet:       3232235776,
			uintMask:      4294967040,
			octets:        []int{192, 168, 1, 10},
			address:       "192.168.1.10",
			broadcast:     "192.168.1.255",
			firstUsable:   "192.168.1.1",
			lastUsable:    "192.168.1.254",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "ffffff00",
			ipNetStr:      "192.168.1.0/24",
			networkStr:    "192.168.1.0/24",
			addrStr:       "192.168.1.10/24",
			maskbits:      24,
			addrBinStr:    "11000000101010000000000100001010",
			addrHexStr:    "c0a8010a",
			pass:          true,
		},
		{ // 14
			input:         "192.168.10.10/16",
			uintAddr:      3232238090,
			uintNet:       3232235520,
			uintMask:      4294901760,
			octets:        []int{192, 168, 10, 10},
			address:       "192.168.10.10",
			broadcast:     "192.168.255.255",
			firstUsable:   "192.168.0.1",
			lastUsable:    "192.168.255.254",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "ffff0000",
			ipNetStr:      "192.168.0.0/16",
			networkStr:    "192.168.0.0/16",
			addrStr:       "192.168.10.10/16",
			maskbits:      16,
			addrBinStr:    "11000000101010000000101000001010",
			addrHexStr:    "c0a80a0a",
			pass:          true,
		},
		{ // 15
			input:         "240.0.0.0/4",
			uintAddr:      4026531840,
			uintNet:       4026531840,
			uintMask:      4026531840,
			octets:        []int{240, 0, 0, 0},
			address:       "240.0.0.0",
			broadcast:     "255.255.255.255",
			firstUsable:   "240.0.0.1",
			lastUsable:    "255.255.255.254",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "f0000000",
			ipNetStr:      "240.0.0.0/4",
			networkStr:    "240.0.0.0/4",
			addrStr:       "240.0.0.0/4",
			maskbits:      4,
			addrBinStr:    "11110000000000000000000000000000",
			addrHexStr:    "f0000000",
			pass:          true,
		},
		{ // 16
			input:         "240.0.0.1/4",
			uintAddr:      4026531841,
			uintNet:       4026531840,
			uintMask:      4026531840,
			octets:        []int{240, 0, 0, 1},
			address:       "240.0.0.1",
			broadcast:     "255.255.255.255",
			firstUsable:   "240.0.0.1",
			lastUsable:    "255.255.255.254",
			listenTCPArgs: []string{"tcp4", ""},
			listenUDPArgs: []string{"udp4", ""},
			ipMaskStr:     "f0000000",
			ipNetStr:      "240.0.0.0/4",
			networkStr:    "240.0.0.0/4",
			addrStr:       "240.0.0.1/4",
			maskbits:      4,
			addrBinStr:    "11110000000000000000000000000001",
			addrHexStr:    "f0000001",
			pass:          true,
		},
		{ // 17
			input:         "255.255.255.255",
			uintAddr:      4294967295,
			uintNet:       4294967295,
			uintMask:      sockaddr.IPv4HostMask,
			octets:        []int{255, 255, 255, 255},
			address:       "255.255.255.255",
			broadcast:     "255.255.255.255",
			firstUsable:   "255.255.255.255",
			lastUsable:    "255.255.255.255",
			listenTCPArgs: []string{"tcp4", "255.255.255.255:0"},
			listenUDPArgs: []string{"udp4", "255.255.255.255:0"},
			ipMaskStr:     "ffffffff",
			ipNetStr:      "255.255.255.255/32",
			networkStr:    "255.255.255.255",
			addrStr:       "255.255.255.255",
			maskbits:      32,
			addrBinStr:    "11111111111111111111111111111111",
			addrHexStr:    "ffffffff",
			pass:          true,
		},
		{ // 18
			input: "www.hashicorp.com",
			pass:  false,
		},
		{ // 19
			input: "2001:DB8::/48",
			pass:  false,
		},
		{ // 20
			input: "2001:DB8::",
			pass:  false,
		},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4Addr(test.input)
		if test.pass && err != nil {
			t.Fatalf("[%d] Unable to create an IPv4Addr from %+q: %v", idx, test, err)
		} else if !test.pass && err == nil {
			t.Fatalf("[%d] Expected test to fail for %+q", idx, test.input)
		} else if !test.pass && err != nil {
			continue
		}

		if type_ := ipv4.Type(); type_ != sockaddr.TypeIPv4 {
			t.Errorf("[%d] Expected new IPv4Addr to be Type %d, received %d (int)", idx, sockaddr.TypeIPv4, type_)
		}

		if c := cap(*ipv4.NetIP()); c != sockaddr.IPv4len {
			t.Errorf("[%d] Expected new IPv4Addr's Address capacity to be %d bytes, received %d", idx, sockaddr.IPv4len, c)
		}

		if l := len(*ipv4.NetIP()); l != sockaddr.IPv4len {
			t.Errorf("[%d] Expected new IPv4Addr's Address length to be %d bytes, received %d", idx, sockaddr.IPv4len, l)
		}

		if a := ipv4.Address; a != test.uintAddr {
			t.Errorf("[%d] Expected %+q's Address to return %d, received %d", idx, test.input, test.uintAddr, a)
		}

		if n, ok := ipv4.Network().(sockaddr.IPv4Addr); !ok || n.Address != test.uintNet {
			t.Errorf("[%d] Expected %+q's Network to return %d, received %d", idx, test.input, test.uintNet, n.Address)

			if n.Mask != test.uintMask {
				t.Errorf("[%d] Expected %+q's Network's Mask to return %d, received %d", idx, test.input, test.uintMask, n.Mask)
			}
		}

		if m := ipv4.Mask; m != test.uintMask {
			t.Errorf("[%d] Expected %+q's Mask to return %d, received %d", idx, test.input, test.uintMask, m)
		}

		if p := ipv4.IPPort(); sockaddr.IPPort(p) != test.port || sockaddr.IPPort(p) != test.port {
			t.Errorf("[%d] Expected %+q's port to be %d, received %d", idx, test.input, test.port, p)
		}

		if o := ipv4.Octets(); len(o) != 4 || o[0] != test.octets[0] || o[1] != test.octets[1] || o[2] != test.octets[2] || o[3] != test.octets[3] {
			t.Errorf("[%d] Expected %+q's Octets to be %+v, received %+v", idx, test.input, test.octets, o)
		}

		if h, ok := ipv4.Host().(sockaddr.IPv4Addr); !ok || h.Address != ipv4.Address || h.Mask != sockaddr.IPv4HostMask || h.Port != ipv4.Port {
			t.Errorf("[%d] Expected %+q's Host() to return identical IPv4Addr except mask, received %+q", idx, test.input, h.String())
		}

		if s := ipv4.NetIP().String(); s != test.address {
			t.Errorf("[%d] Expected %+q's address to be %+q, received %+q", idx, test.input, test.address, s)
		}

		if s := ipv4.AddressBinString(); s != test.addrBinStr {
			t.Errorf("[%d] Expected address %+q's binary representation to be %+q, received %+q", idx, test.input, test.addrBinStr, s)
		}

		if s := ipv4.AddressHexString(); s != test.addrHexStr {
			t.Errorf("[%d] Expected address %+q's hexadecimal representation to be %+q, received %+q", idx, test.input, test.addrHexStr, s)
		}

		if b := ipv4.Broadcast().String(); b != test.broadcast {
			t.Errorf("[%d] Expected %+q's broadcast to be %+q, received %+q", idx, test.input, test.broadcast, b)
		}

		if f := ipv4.FirstUsable().String(); f != test.firstUsable {
			t.Errorf("[%d] Expected %+q's FirstUsable() to be %+q, received %+q", idx, test.input, test.firstUsable, f)
		}

		if listenNet, listenArgs := ipv4.ListenTCPArgs(); listenNet != test.listenTCPArgs[0] || listenArgs != test.listenTCPArgs[1] {
			t.Errorf("[%d] Expected %+q's ListenArgs() to be %+q, received %+q, %+q", idx, test.input, test.listenTCPArgs, listenNet, listenArgs)
		}

		if listenNet, listenArgs := ipv4.ListenUDPArgs(); listenNet != test.listenUDPArgs[0] || listenArgs != test.listenUDPArgs[1] {
			t.Errorf("[%d] Expected %+q's ListenArgs() to be %+q, received %+q, %+q", idx, test.input, test.listenUDPArgs, listenNet, listenArgs)
		}

		if l := ipv4.LastUsable().String(); l != test.lastUsable {
			t.Errorf("[%d] Expected %+q's LastUsable() to be %+q, received %+q", idx, test.input, test.lastUsable, l)
		}

		if m := ipv4.NetIPMask().String(); m != test.ipMaskStr {
			t.Errorf("[%d] Expected %+q's mask to be %+q, received %+q", idx, test.input, test.ipMaskStr, m)
		}

		if n := ipv4.NetIPNet().String(); n != test.ipNetStr {
			t.Errorf("[%d] Expected %+q's network to be %+q, received %+q", idx, test.input, test.ipNetStr, n)
		}

		if n := ipv4.Network().String(); n != test.networkStr {
			t.Errorf("[%d] Expected %+q's Network() to be %+q, received %+q", idx, test.input, test.networkStr, n)
		}

		if m := ipv4.Maskbits(); m != test.maskbits {
			t.Errorf("[%dr] Expected %+q's port to be %d, received %d", idx, test.input, test.maskbits, m)
		}

		if s := ipv4.String(); s != test.addrStr {
			t.Errorf("[%d] Expected %+q's String to be %+q, received %+q", idx, test.input, test.addrStr, s)
		}
	}
}

func TestSockAddr_IPv4Addr_Equal(t *testing.T) {
	tests := []struct {
		input string
		pass  []string
		fail  []string
	}{
		{
			input: "208.67.222.222/32",
			pass:  []string{"208.67.222.222", "208.67.222.222/32"},
			fail:  []string{"208.67.222.222/31", "208.67.220.220", "208.67.220.220/32"},
		},
		{
			input: "4.2.2.1",
			pass:  []string{"4.2.2.1", "4.2.2.1/32"},
			fail:  []string{"4.2.2.1/0", "4.2.2.2", "4.2.2.2/32"},
		},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4Addr(test.input)
		if err != nil {
			t.Fatalf("[%d] Unable to create an IPv4Addr from %+q: %v", idx, test.input, err)
		}

		for goodIdx, passInput := range test.pass {
			good, err := sockaddr.NewIPv4Addr(passInput)
			if err != nil {
				t.Fatalf("[%d] Unable to create an IPv4Addr from %+q: %v", idx, passInput, err)
			}

			if !ipv4.Equal(good) {
				t.Errorf("[%d/%d] Expected %+q to be equal to %+q: %+q/%+q", idx, goodIdx, test.input, passInput, ipv4.String(), good.String())
			}
		}

		for failIdx, failInput := range test.fail {
			fail, err := sockaddr.NewIPv4Addr(failInput)
			if err != nil {
				t.Fatalf("[%d] Unable to create an IPv4Addr from %+q: %v", idx, failInput, err)
			}

			if ipv4.Equal(fail) {
				t.Errorf("[%d/%d] Expected %+q to be not equal to %+q", idx, failIdx, test.input, failInput)
			}
		}
	}
}
