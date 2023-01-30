package sockaddr

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"
)

// GetAllInterfaces iterates over all available network interfaces and finds all
// available IP addresses on each interface and converts them to
// sockaddr.IPAddrs, and returning the result as an array of IfAddr.
func GetAllInterfaces() (IfAddrs, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Calling Addrs() on each net.Interface has poor O(n^2) performance over a
	// large number of interfaces: https://github.com/golang/go/issues/53660
	// Instead, get all addresses in a single netlink call and map them to
	// their associated interfaces.
	ifMap := make(map[int]net.Interface)
	for _, intf := range ifs {
		ifMap[intf.Index] = intf
	}

	tab, err := syscall.NetlinkRIB(syscall.RTM_GETADDR, syscall.AF_UNSPEC)
	if err != nil {
		return nil, os.NewSyscallError("netlinkrib", err)
	}
	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return nil, os.NewSyscallError("parsenetlinkmessage", err)
	}

	ifAddrs := make(IfAddrs, 0, len(msgs))
	for _, m := range msgs {
		if m.Header.Type == syscall.RTM_NEWADDR {
			ifam := (*syscall.IfAddrmsg)(unsafe.Pointer(&m.Data[0]))
			attrs, err := syscall.ParseNetlinkRouteAttr(&m)
			if err != nil {
				return nil, os.NewSyscallError("parsenetlinkrouteattr", err)
			}
			addr := newAddr(ifam, attrs)
			ipAddr, err := NewIPAddr(addr.String())
			if err != nil {
				return IfAddrs{}, fmt.Errorf("unable to create an IP address from %q", addr.String())
			}
			if intf, ok := ifMap[int(ifam.Index)]; ok {
				ifAddrs = append(ifAddrs, IfAddr{
					SockAddr:  ipAddr,
					Interface: intf,
				})
			}
		}
	}

	return ifAddrs, nil
}

// Vendored unexported function:
// https://github.com/golang/go/blob/8bcc490667d4dd44c633c536dd463bbec0a3838f/src/net/interface_linux.go#L178-L203
func newAddr(ifam *syscall.IfAddrmsg, attrs []syscall.NetlinkRouteAttr) net.Addr {
	var ipPointToPoint bool
	// Seems like we need to make sure whether the IP interface
	// stack consists of IP point-to-point numbered or unnumbered
	// addressing.
	for _, a := range attrs {
		if a.Attr.Type == syscall.IFA_LOCAL {
			ipPointToPoint = true
			break
		}
	}
	for _, a := range attrs {
		if ipPointToPoint && a.Attr.Type == syscall.IFA_ADDRESS {
			continue
		}
		switch ifam.Family {
		case syscall.AF_INET:
			return &net.IPNet{IP: net.IPv4(a.Value[0], a.Value[1], a.Value[2], a.Value[3]), Mask: net.CIDRMask(int(ifam.Prefixlen), 8*IPv4len)}
		case syscall.AF_INET6:
			ifa := &net.IPNet{IP: make(net.IP, IPv6len), Mask: net.CIDRMask(int(ifam.Prefixlen), 8*IPv6len)}
			copy(ifa.IP, a.Value[:])
			return ifa
		}
	}
	return nil
}
