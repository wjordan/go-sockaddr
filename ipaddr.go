package sockaddr

import "net"

// IPAddr is a generic IP address interface for IPv4 and IPv6 addresses,
// networks, and endpoints.
type IPAddr interface {
	SockAddr
	AddressBinString() string
	AddressHexString() string
	DialPacketArgs() (string, string)
	DialStreamArgs() (string, string)
	FirstUsable() IPAddr
	Host() IPAddr
	IPPort() uint16
	LastUsable() IPAddr
	ListenPacketArgs() (string, string)
	ListenStreamArgs() (string, string)
	Maskbits() int
	NetIP() *net.IP
	NetIPNet() *net.IPNet
	Network() IPAddr
	SetPort(uint16) IPAddr
	String() string
}

type IPPort uint16
