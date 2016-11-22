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
func (ifAddr IfAddr) Attr(attrName AttrName) (string, error) {
	sa := ifAddr.SockAddr
	switch sockType := sa.Type(); {
	case sockType&TypeIP != 0:
		ip := *ToIPAddr(sa)
		attrVal := IPAddrAttr(ip, attrName)
		if attrVal != "" {
			return attrVal, nil
		}

		if sa.Type() == TypeIPv4 {
			ipv4 := *ToIPv4Addr(sa)
			attrVal := IPv4AddrAttr(ipv4, attrName)
			if attrVal != "" {
				return attrVal, nil
			}
		}

		if sa.Type() == TypeIPv6 {
			ipv6 := *ToIPv6Addr(sa)
			attrVal := IPv6AddrAttr(ipv6, attrName)
			if attrVal != "" {
				return attrVal, nil
			}
		}

		// Random attribute names that are Interface specific
		switch attrName {
		case "name":
			return ifAddr.Interface.Name, nil
		case "flags":
			return ifAddr.Interface.Flags.String(), nil
		}
	case sockType == TypeUnix:
		us := *ToUnixSock(sa)
		attrVal := UnixSockAttr(us, attrName)
		if attrVal != "" {
			return attrVal, nil
		}
	}

	// Non type-specific attributes
	switch attrName {
	case "string":
		return sa.String(), nil
	case "type":
		return sa.Type().String(), nil
	}

	return "", fmt.Errorf("unsupported attribute name %q", attrName)
}
