package sockaddr

import "net"

// IfAddrs is a slice of IPAddrs for per interface
type IfAddrs struct {
	Addrs []IPAddr
	net.Interface
}

// GetIfAddrs iterates over all available network interfaces and finds all
// available IP addresses on each interface and converts them to
// sockaddr.IPAddrs.
func GetIfAddrs() ([]IfAddrs, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ifAddrs := make([]IfAddrs, 0, len(ifs))
	for _, intf := range ifs {
		addrs, err := intf.Addrs()
		if err != nil {
			return nil, err
		}

		ipAddrs := make([]IPAddr, 0, len(addrs))
		for _, addr := range addrs {
			ipAddr, err := NewIPAddr(addr)
			if err != nil {
				continue
			}
			ipAddrs = append(ipAddrs, ipAddr)
		}

		ifAddr := IfAddrs{
			Addrs:     ipAddrs,
			Interface: intf,
		}
		ifAddrs = append(ifAddrs, ifAddr)
	}

	return ifAddrs, nil
}
