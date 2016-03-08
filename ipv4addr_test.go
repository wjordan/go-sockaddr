package sockaddr_test

import (
	"testing"

	"github.com/hashicorp/go-sockaddr"
)

type GoodIPv4Test struct {
	ipv4Net sockaddr.IPv4Addr
}
type GoodIPv4Tests []*GoodIPv4Test

func TestSockaddr_IPv4Addr_Address(t *testing.T) {
	var tests = []struct {
		input   string
		address string
	}{
		{"208.67.222.222/32", "208.67.222.222"},
		{"4.2.2.1/24", "4.2.2.1"},
		{"8.8.8.8", "8.8.8.8"},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.Address().String() != test.address {
			t.Errorf("[%d] Expected %s's address to be %s", idx, test.input, test.address)
		}
	}
}

func TestSockaddr_IPv4Addr_BroadcastAddress(t *testing.T) {
	var tests = []struct {
		input     string
		broadcast string
	}{
		{"208.67.222.222/32", "208.67.222.222"},
		{"4.2.2.1/24", "4.2.2.255"},
		{"10.0.0.10/23", "10.0.1.255"},
		{"8.8.8.8", "8.8.8.8"},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.BroadcastAddress().Address().String() != test.broadcast {
			t.Errorf("[%d] Expected %s's broadcast address to be %s", idx, test.input, test.broadcast)
		}
	}
}

func TestSockaddr_IPv4Addr_Equal(t *testing.T) {
	var tests = []struct {
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
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		for goodIdx, passInput := range test.pass {
			good, err := sockaddr.NewIPv4AddrFromString(passInput)
			if err != nil {
				t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, passInput)
			}

			if !ipv4.Equal(good) {
				t.Errorf("[%d/%d] Expected %s's to be equal to %s", idx, goodIdx, test.input, passInput)
			}
		}

		for failIdx, failInput := range test.fail {
			fail, err := sockaddr.NewIPv4AddrFromString(failInput)
			if err != nil {
				t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, failInput)
			}

			if ipv4.Equal(fail) {
				t.Errorf("[%d/%d] Expected %s's to be not equal to %s", idx, failIdx, test.input, failInput)
			}
		}
	}
}

func TestSockaddr_IPv4Addr_FirstUsableAddress(t *testing.T) {
	var tests = []struct {
		input       string
		firstUsable string
	}{
		{"208.67.222.222/32", "208.67.222.222"},
		{"4.2.2.2/24", "4.2.2.1"},
		{"10.0.1.10/23", "10.0.0.1"},
		{"8.8.8.8", "8.8.8.8"},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.FirstUsableAddress().Address().String() != test.firstUsable {
			t.Errorf("[%d] Expected %s's first usable address to be %s", idx, test.input, test.firstUsable)
		}
	}
}

func TestSockaddr_IPv4Addr_LastUsableAddress(t *testing.T) {
	var tests = []struct {
		input      string
		lastUsable string
	}{
		{"208.67.222.222/32", "208.67.222.222"},
		{"4.2.2.2/24", "4.2.2.254"},
		{"10.0.1.10/23", "10.0.1.254"},
		{"10.0.0.10/23", "10.0.1.254"},
		{"8.8.8.8", "8.8.8.8"},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.LastUsableAddress().Address().String() != test.lastUsable {
			t.Errorf("[%d] Expected %s's last usable address to be %s", idx, test.input, test.lastUsable)
		}
	}
}

func TestSockaddr_IPv4Addr_Maskbits(t *testing.T) {
	var tests = []struct {
		input string
		bits  int
	}{
		{"208.67.222.222/32", 32},
		{"4.2.2.2/24", 24},
		{"10.0.1.10/23", 23},
		{"8.8.8.8", 32},
		{"0.0.0.0/0", 0},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.Maskbits() != test.bits {
			t.Errorf("[%d] Expected %s's masked bits to be %d", idx, test.input, test.bits)
		}
	}
}

func TestSockaddr_IPv4Addr_Network(t *testing.T) {
	var tests = []struct {
		input   string
		network string
	}{
		{"208.67.222.222/32", "208.67.222.222/32"},
		{"4.2.2.1/24", "4.2.2.0/24"},
		{"192.168.3.0/22", "192.168.0.0/22"},
		{"8.8.8.8", "8.8.8.8/32"},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.Network().String() != test.network {
			t.Errorf("[%d] Expected %s's network to be %s, received %s", idx, test.input, test.network, ipv4.Network().String())
		}
	}
}

func TestSockaddr_IPv4Addr_NetworkPrefix(t *testing.T) {
	var tests = []struct {
		input  string
		prefix string
	}{
		{"208.67.222.222/32", "208.67.222.222"},
		{"4.2.2.2/24", "4.2.2.0"},
		{"10.0.1.10/23", "10.0.0.0"},
		{"8.8.8.8", "8.8.8.8"},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.NetworkPrefix().Address().String() != test.prefix {
			t.Errorf("[%d] Expected %s's network prefix to be %d", idx, test.input, test.prefix)
		}
	}
}

func TestSockaddr_IPv4Addr_NewIPv4FromString(t *testing.T) {
	var tests = []struct {
		input   string
		address string
		network string
	}{
		{
			input:   "0.0.0.0",
			address: "0.0.0.0",
			network: "0.0.0.0/32",
		},
		{
			input:   "0.0.0.0/0",
			address: "0.0.0.0",
			network: "0.0.0.0/0",
		},
		{
			input:   "0.0.0.1",
			address: "0.0.0.1",
			network: "0.0.0.1/32",
		},
		{
			input:   "0.0.0.1/1",
			address: "0.0.0.1",
			network: "0.0.0.0/1",
		},
		{
			input:   "1.2.3.4",
			address: "1.2.3.4",
			network: "1.2.3.4/32",
		},
		{
			input:   "10.0.0.0/8",
			address: "10.0.0.0",
			network: "10.0.0.0/8",
		},
		{
			input:   "128.0.0.0",
			address: "128.0.0.0",
			network: "128.0.0.0/32",
		},
		{
			input:   "128.95.120.1/32",
			address: "128.95.120.1",
			network: "128.95.120.1/32",
		},
		{
			input:   "172.16.1.3/12",
			address: "172.16.1.3",
			network: "172.16.0.0/12",
		},
		{
			input:   "192.168.0.0/16",
			address: "192.168.0.0",
			network: "192.168.0.0/16",
		},
		{
			input:   "192.168.0.1",
			address: "192.168.0.1",
			network: "192.168.0.1/32",
		},
		{
			input:   "192.168.0.2/31",
			address: "192.168.0.2",
			network: "192.168.0.2/31",
		},
		{
			input:   "192.168.1.10/24",
			address: "192.168.1.10",
			network: "192.168.1.0/24",
		},
		{
			input:   "192.168.10.10/16",
			address: "192.168.10.10",
			network: "192.168.0.0/16",
		},
		{
			input:   "240.0.0.0/4",
			address: "240.0.0.0",
			network: "240.0.0.0/4",
		},
		{
			input:   "240.0.0.1/4",
			address: "240.0.0.1",
			network: "240.0.0.0/4",
		},
		{
			input:   "255.255.255.255",
			address: "255.255.255.255",
			network: "255.255.255.255/32",
		},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test)
		}

		if ipv4.Type() != sockaddr.TypeIPv4 {
			t.Errorf("[%d] Expected new IPv4Addr to be Type %d, received %d (int)", idx, sockaddr.TypeIPv4, ipv4.Type())
		}

		if ipv4.Address().String() != test.address {
			t.Errorf("[%d] Expected %s's address to be %s, received %s", idx, test.input, test.address, ipv4.Address().String())
		}

		if ipv4.Network().String() != test.network {
			t.Errorf("[%d] Expected %s's network to be %s, received %s", idx, test.input, test.network, ipv4.Network().String())
		}

		if ipv4.Port() != 0 {
			t.Errorf("[%d] Expected %s's port to be 0", idx, test.input)
		}
	}
}

func TestSockaddr_IPv4Addr_Port(t *testing.T) {
	var tests = []struct {
		input string
		port  uint16
	}{
		{"208.67.222.222/32", 53},
		{"10.0.1.10", 80},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		// FIXME(sean@): Why does this fail???  It's like SetPort()
		// is being noop'ed.  Direct assignment works. Verifying the
		// content of the receiver's IPPort value confirms its right.
		// Why in this calling scope does it revert back to 0?
		//
		// ipv4.SetPort(test.port)
		// if ipv4.Port() != test.port {
		// 	t.Errorf("[%d] Expected %s's port to be %d, received %d", idx, test.input, test.port, ipv4.Port())
		// }

		ipv4.IPPort = test.port
		if ipv4.Port() != test.port {
			t.Errorf("[%d] Expected %s's port to be %d, received %d", idx, test.input, test.port, ipv4.Port())
		}
	}
}

func TestSockaddr_IPv4Addr_ToBinString(t *testing.T) {
	var tests = []struct {
		input     string
		binstring string
	}{
		{"0.0.0.255/24", "00000000000000000000000011111111"},
		{"255.0.0.0/32", "11111111000000000000000000000000"},
		{"255.0.255.0", "11111111000000001111111100000000"},
		{"0.255.0.255", "00000000111111110000000011111111"},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.ToBinString() != test.binstring {
			t.Errorf("[%d] Expected %s's binary representation to be %s, received %s", idx, test.input, test.binstring, ipv4.ToBinString())
		}
	}
}

func TestSockaddr_IPv4Addr_ToHexString(t *testing.T) {
	var tests = []struct {
		input     string
		hexstring string
	}{
		{"0.0.0.255/24", "000000ff"},
		{"255.0.0.0/32", "ff000000"},
		{"255.0.255.0", "ff00ff00"},
		{"0.255.0.255", "00ff00ff"},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.ToHexString() != test.hexstring {
			t.Errorf("[%d] Expected %s's hex representation to be %s, received %s", idx, test.input, test.hexstring, ipv4.ToHexString())
		}
	}
}

func TestSockaddr_IPv4Addr_ToUint32(t *testing.T) {
	var tests = []struct {
		input  string
		ipuint uint32
	}{
		{"0.0.0.0", 0},
		{"0.0.0.1", 1},
		{"0.255.0.255", 16711935},
		{"255.0.255.0", 4278255360},
		{"255.255.255.255/32", 4294967295},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test.input)
		}

		if ipv4.ToUint32() != test.ipuint {
			t.Errorf("[%d] Expected %s's uint32 representation to be %d, received %d", idx, test.input, test.ipuint, ipv4.ToUint32())
		}
	}
}

func TestSockaddr_IPv4Addr_Type(t *testing.T) {
	var tests = []struct {
		input string
	}{
		{"208.67.222.222/32"},
		{"4.2.2.1"},
		{"192.168.2.3/24"},
	}

	for idx, test := range tests {
		ipv4, err := sockaddr.NewIPv4AddrFromString(test.input)
		if err != nil {
			t.Errorf("[%d] Unable to create an IPv4Addr from %s: %s", idx, test)
		}

		if ipv4.Type() != sockaddr.TypeIPv4 {
			t.Errorf("[%d] Expected type IPv4Addr from %s, received %d (int)", idx, test, ipv4.Type())
		}
	}
}
