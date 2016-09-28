package sockaddr

import (
	"fmt"
	"net"
)

// Constants for the sizes of IPv3, IPv4, and IPv6 address types.
const (
	IPv3len = 6
	IPv4len = 4
	IPv6len = 16
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

// IPPort is the type for an IP port number for the TCP and UDP IP transports.
type IPPort uint16

// IPPrefixLen is a typed integer representing the prefix length for a given
// IPAddr.
type IPPrefixLen byte

// NewIPAddr creates a new IPAddr from a string.  Returns nil if the string is
// not an IPv4 or an IPv6 address.
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
