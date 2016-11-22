package sockaddr

import (
	"fmt"
	"net"
)

// IfAddr is a union of a SockAddr and a net.Interface.
type IfAddr struct {
	SockAddr
	net.Interface
}

// Attr returns the named attribute as a string
func (ifAddr IfAddr) Attr(attrName AttrName) string {
	sa := ifAddr.SockAddr
	switch sockType := sa.Type(); {
	case sockType&TypeIP != 0:
		ip := *ToIPAddr(sa)
		attrVal := IPAddrAttr(ip, attrName)
		if attrVal != "" {
			return attrVal
		}

		if sa.Type() == TypeIPv4 {
			ipv4 := *ToIPv4Addr(sa)
			attrVal := IPv4AddrAttr(ipv4, attrName)
			if attrVal != "" {
				return attrVal
			}
		}

		if sa.Type() == TypeIPv6 {
			ipv6 := *ToIPv6Addr(sa)
			attrVal := IPv6AddrAttr(ipv6, attrName)
			if attrVal != "" {
				return attrVal
			}
		}

		// Random attribute names that are Interface specific
		switch attrName {
		case "name":
			return ifAddr.Interface.Name
		case "flags":
			return ifAddr.Interface.Flags.String()
		}
	case sockType == TypeUnix:
		us := *ToUnixSock(sa)
		attrVal := UnixSockAttr(us, attrName)
		if attrVal != "" {
			return attrVal
		}
	}

	// Non type-specific attributes
	switch attrName {
	case "string":
		return sa.String()
	case "type":
		return sa.Type().String()
	}

	return fmt.Sprintf("<unsupported attribute name %q>", attrName)
}
