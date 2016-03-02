package netaddr_test

import (
	"testing"

	"github.com/hashicorp/go-netaddr"
)

func TestNetAddr_New(t *testing.T) {
	type NetAddrFixture struct {
		IPv4               string
		NetworkAddress     string
		BroadcastAddress   string
		IPUint32           uint32
		Maskbits           int
		BinString          string
		HexString          string
		FirstUsableAddress string
		LastUsableAddress  string
	}
	type NetAddrFixtures []NetAddrFixtures

	goodResults := []NetAddrFixture{
		{
			IPv4:               "0.0.0.0",
			NetworkAddress:     "0.0.0.0",
			BroadcastAddress:   "0.0.0.0",
			Maskbits:           32,
			IPUint32:           0,
			BinString:          "00000000000000000000000000000000",
			HexString:          "00000000",
			FirstUsableAddress: "0.0.0.0",
			LastUsableAddress:  "0.0.0.0",
		},
		{
			IPv4:               "0.0.0.0/0",
			NetworkAddress:     "0.0.0.0",
			BroadcastAddress:   "255.255.255.255",
			Maskbits:           0,
			IPUint32:           0,
			BinString:          "00000000000000000000000000000000",
			HexString:          "00000000",
			FirstUsableAddress: "0.0.0.1",
			LastUsableAddress:  "255.255.255.254",
		},
		{
			IPv4:               "0.0.0.1",
			NetworkAddress:     "0.0.0.1",
			BroadcastAddress:   "0.0.0.1",
			Maskbits:           32,
			IPUint32:           1,
			BinString:          "00000000000000000000000000000001",
			HexString:          "00000001",
			FirstUsableAddress: "0.0.0.1",
			LastUsableAddress:  "0.0.0.1",
		},
		{
			IPv4:               "0.0.0.1/1",
			NetworkAddress:     "0.0.0.0",
			BroadcastAddress:   "127.255.255.255",
			Maskbits:           1,
			IPUint32:           1,
			BinString:          "00000000000000000000000000000001",
			HexString:          "00000001",
			FirstUsableAddress: "0.0.0.1",
			LastUsableAddress:  "127.255.255.254",
		},
		{
			IPv4:               "128.0.0.0",
			NetworkAddress:     "128.0.0.0",
			BroadcastAddress:   "128.0.0.0",
			Maskbits:           32,
			IPUint32:           2147483648,
			BinString:          "10000000000000000000000000000000",
			HexString:          "80000000",
			FirstUsableAddress: "128.0.0.0",
			LastUsableAddress:  "128.0.0.0",
		},
		{
			IPv4:               "255.255.255.255",
			NetworkAddress:     "255.255.255.255",
			BroadcastAddress:   "255.255.255.255",
			Maskbits:           32,
			IPUint32:           4294967295,
			BinString:          "11111111111111111111111111111111",
			HexString:          "ffffffff",
			FirstUsableAddress: "255.255.255.255",
			LastUsableAddress:  "255.255.255.255",
		},
		{
			IPv4:               "1.2.3.4",
			NetworkAddress:     "1.2.3.4",
			BroadcastAddress:   "1.2.3.4",
			Maskbits:           32,
			IPUint32:           16909060,
			BinString:          "00000001000000100000001100000100",
			HexString:          "01020304",
			FirstUsableAddress: "1.2.3.4",
			LastUsableAddress:  "1.2.3.4",
		},
		{
			IPv4:               "192.168.10.10/16",
			NetworkAddress:     "192.168.0.0",
			BroadcastAddress:   "192.168.255.255",
			Maskbits:           16,
			IPUint32:           3232238090,
			BinString:          "11000000101010000000101000001010",
			HexString:          "c0a80a0a",
			FirstUsableAddress: "192.168.0.1",
			LastUsableAddress:  "192.168.255.254",
		},
		{
			IPv4:               "192.168.1.10/24",
			NetworkAddress:     "192.168.1.0",
			BroadcastAddress:   "192.168.1.255",
			Maskbits:           24,
			IPUint32:           3232235786,
			BinString:          "11000000101010000000000100001010",
			HexString:          "c0a8010a",
			FirstUsableAddress: "192.168.1.1",
			LastUsableAddress:  "192.168.1.254",
		},
		{
			IPv4:               "192.168.0.1",
			NetworkAddress:     "192.168.0.1",
			BroadcastAddress:   "192.168.0.1",
			Maskbits:           32,
			IPUint32:           3232235521,
			BinString:          "11000000101010000000000000000001",
			HexString:          "c0a80001",
			FirstUsableAddress: "192.168.0.1",
			LastUsableAddress:  "192.168.0.1",
		},
		{
			IPv4:               "192.168.0.2/31",
			NetworkAddress:     "192.168.0.2",
			BroadcastAddress:   "192.168.0.3",
			Maskbits:           31,
			IPUint32:           3232235522,
			BinString:          "11000000101010000000000000000010",
			HexString:          "c0a80002",
			FirstUsableAddress: "192.168.0.2",
			LastUsableAddress:  "192.168.0.3",
		},
		{
			IPv4:               "240.0.0.0/4",
			NetworkAddress:     "240.0.0.0",
			BroadcastAddress:   "255.255.255.255",
			Maskbits:           4,
			IPUint32:           4026531840,
			BinString:          "11110000000000000000000000000000",
			HexString:          "f0000000",
			FirstUsableAddress: "240.0.0.1",
			LastUsableAddress:  "255.255.255.254",
		},
	}

	for _, r := range goodResults {
		var (
			addr *netaddr.NetAddr
			str  string
		)

		na, err := netaddr.New(r.IPv4)
		if err != nil {
			t.Fatalf("Failed parse %s", r.IPv4)
		}

		maskbits := na.Maskbits()
		if maskbits != r.Maskbits {
			t.Fatalf("Failed Maskbits %s: %d != %d", r.IPv4, maskbits, r.Maskbits)
		}

		ipuint, ok := na.ToUint32()
		if !ok || ipuint != r.IPUint32 {
			t.Fatalf("Failed ToUint32() %s: %d != %d", r.IPv4, ipuint, r.IPUint32)
		}

		str = na.ToBinString()
		if str != r.BinString {
			t.Fatalf("Failed BinString %s: %s != %s", r.IPv4, str, r.BinString)
		}

		str = na.ToHexString()
		if str != r.HexString {
			t.Fatalf("Failed HexString %s: %s != %s", r.IPv4, str, r.HexString)
		}

		addr = na.BroadcastAddress()
		if addr == nil || addr.Address.To4().String() != r.BroadcastAddress {
			t.Fatalf("Failed BroadcastAddress %s: %s != %s", r.IPv4, addr.Address.To4().String(), r.BroadcastAddress)
		}

		addr = na.NetworkAddress()
		if addr == nil || addr.Address.To4().String() != r.NetworkAddress {
			t.Fatalf("Failed NetworkAddress %s: %s != %s", r.IPv4, addr.Address.To4().String(), r.NetworkAddress)
		}

		addr = na.FirstUsableAddress()
		if addr == nil || addr.Address.To4().String() != r.FirstUsableAddress {
			t.Fatalf("Failed FirstUsableAddress %s: %s != %s", r.IPv4, addr.Address.To4().String(), r.FirstUsableAddress)
		}

		addr = na.LastUsableAddress()
		if addr == nil || addr.Address.To4().String() != r.LastUsableAddress {
			t.Fatalf("Failed LastUsableAddress %s: %s != %s", r.IPv4, addr.Address.To4().String(), r.LastUsableAddress)
		}
	}

	badResults := []string{
		"256.0.0.0",
		"0.0.0.0.0",
	}

	for _, badIP := range badResults {
		na, err := netaddr.New(badIP)
		if err == nil {
			t.Fatalf("Failed should have failed to parse %s", badIP)
		}
		if na != nil {
			t.Fatalf("NetAddr should be nil")
		}
	}
}
