package sockaddr

import (
	"fmt"
	"strings"
)

type SockAddrType int

const (
	TypeUnknown SockAddrType = 0x0
	TypeUnix                 = 0x1
	TypeIPv4                 = 0x2
	TypeIPv6                 = 0x4

	// TypeIP is the union of TypeIPv4 and TypeIPv6
	TypeIP = 0x6
)

type SockAddr interface {
	// CmpRFC returns 0 if SockAddr exactly matches one of the matched RFC
	// networks, -1 if the receiver is contained within the RFC network, or
	// 1 if the address is not contained within the RFC.
	CmpRFC(rfcNum uint, sa SockAddr) int

	// Contains returns true if the SockAddr arg is contained within the
	// receiver
	Contains(SockAddr) bool

	// Equal allows for the comparison of two SockAddrs
	Equal(SockAddr) bool

	// // ListenArgs returns the necessary arguments required for
	// // stream-based communication such as net.Listen().
	// ListenArgs() (net.Listener, error)

	// // ListenPacketArgs returns the necessary arguments required for
	// // packet-based communication such as net.ListenPacket().
	// ListenPacketArgs() (net.PacketConn, error)

	// String returns the string representation of SockAddr
	String() string

	// Type returns the SockAddrType
	Type() SockAddrType

	// // ToUint32 returns MaxUint32 for non-uint32
	// ToUint32() uint32
}

// New creates a new SockAddr from the string.  The order in which New()
// attempts to construct a SockAddr is: SockAddrUnix, IPv4Addr, IPv6Addr.
//
// NOTE: New() relies on the heuristic wherein the path begins with either a '.'
// or '/' character before creating a new UnixSock.  For UNIX sockets that are
// absolute paths or are nested within a sub-directory, this works as expected,
// however if the UNIX socket is contained in the current working directory,
// this will fail unless the path begins with "./" (e.g. "./my-local-socket").
// Calls directly to NewUnixSock() do not suffer this limitation.  Invalid IP
// addresses such as "256.0.0.0/-1" will run afoul of this heuristic and be
// assumed to be a valid UNIX socket path (which they are, but it is probably
// not what you want and you won't realize it until you stat(2) the file system
// to discover it doesn't exist).
func NewSockAddr(s string) (SockAddr, error) {
	ipv4Addr, err := NewIPv4Addr(s)
	if err == nil {
		return ipv4Addr, nil
	}

	ipv6Addr, err := NewIPv6Addr(s)
	if err == nil {
		return ipv6Addr, nil
	}

	// Check to make sure the string begins with either a '.' or '/', or
	// contains a '/'.  Not gunna do it.  Not gunna reverse that order.
	if len(s) > 1 && (strings.IndexAny(s[0:1], "./") != -1 || strings.IndexByte(s, '/') != -1) {
		unixSock, err := NewUnixSock(s)
		if err == nil {
			return unixSock, nil
		}
	}

	return nil, fmt.Errorf("Unable to convert %s to an IPv4 or IPv6 address, or a UNIX Socket", s)
}

// IsRFC tests to see if an SockAddr matches the specified RFC
func IsRFC(rfcNum uint, sa SockAddr) bool {
	rfcNets, ok := rfcNetMap[rfcNum]
	if !ok {
		return false
	}

	var contained bool
	for _, rfcNet := range rfcNets {
		if rfcNet.Contains(sa) {
			contained = true
			break
		}
	}
	return contained
}

// func NewPort(s string, port uint16) (ipa IPAddr, err error) {
// 	ipa, err = New(s)
// 	if err == nil {
// 		ipa.SetPort(port)
// 	}

// 	return ipa, err
// }

// func (ipa *IPAddr) ToIPAddr() *IPAddr {
// 	switch ipa.Type() {
// 	case TypeIPv4, TypeIPv6:
// 		ipa, ok := ipa.(*IPAddr)
// 		if !ok {
// 			return nil
// 		}
// 		return ipa.ToIPv4Addr()
// 	default:
// 		return nil
// 	}
// }

func ToIPv4Addr(sa SockAddr) *IPv4Addr {
	switch sa.Type() {
	case TypeIPv4:
		ipv4, ok := sa.(IPv4Addr)
		if !ok {
			return nil
		}
		return &ipv4
		//return ipv4.ToIPv4Addr()
	default:
		return nil
	}
}

// func (ipa *IPAddr) ToIPv6Addr() *IPv4Addr {
// 	switch ipa.Type() {
// 	case TypeIPv6:
// 		ipa, ok := ipa.(*IPAddr)
// 		if !ok {
// 			return nil
// 		}
// 		return ipa.ToIPv6Addr()
// 	default:
// 		return nil
// 	}
// }

// String() for SockAddrType returns a string representation of the
// SockAddrType (e.g. "IPv4", "IPv6", "UNIX", "IP", or "unknown").
func (sat SockAddrType) String() string {
	switch sat {
	case TypeIPv4:
		return "IPv4"
	case TypeIPv6:
		return "IPv6"
	case TypeIP:
		return "IP"
	case TypeUnix:
		return "UNIX"
	default:
		return "unknown"
	}
}
