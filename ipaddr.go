package sockaddr

import (
	"fmt"
	"net"
)

// IPAddr is a generic IP address interface for IPv4 and IPv6 addresses,
// networks, and socket endpoints.
type IPAddr interface {
	SockAddr
	AddressBinString() string
	AddressHexString() string
	Cmp(SockAddr) int
	CmpAddress(SockAddr) int
	CmpPort(SockAddr) int
	DialPacketArgs() (string, string)
	DialStreamArgs() (string, string)
	FirstUsable() IPAddr
	Host() IPAddr
	IPPort() IPPort
	LastUsable() IPAddr
	ListenPacketArgs() (string, string)
	ListenStreamArgs() (string, string)
	Maskbits() int
	NetIP() *net.IP
	NetIPMask() *net.IPMask
	NetIPNet() *net.IPNet
	Network() IPAddr
}

type IPPort uint16

// NewIPAddr creates a new IPAddr from a string.
func NewIPAddr(addr string) (IPAddr, error) {
	ipv4Addr, err := NewIPv4Addr(addr)
	if err == nil {
		return ipv4Addr, nil
	}

	ipv6Addr, err := NewIPv6Addr(addr)
	if err == nil {
		return ipv6Addr, nil
	}

	return nil, fmt.Errorf("invalid IPAddr %v", addr)
}
