package sockaddr

// IPv4Addrs is a collection of IPv4Addrs
type IPv4Addrs []IPv4Addr

// Len returns the length of the IPv4Addrs array
func (s IPv4Addrs) Len() int { return len(s) }

// Swap swaps the elements with indexes i and j.
func (s IPv4Addrs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// SortIPv4AddrsByNetwork is a type that satisfies sort.Interface and can be
// used by the routines in this package.  The SortIPv4AddrsByNetwork type is
// used to sort IPv4Addrs by numerical order from lowest to highest.
type SortIPAddrsByNetwork struct{ IPv4Addrs }

// Less reports whether the element with index i should sort before the
// element with index j.
func (s SortIPAddrsByNetwork) Less(i, j int) bool {
	iAddr := s.IPv4Addrs[i].NetworkAddress()
	jAddr := s.IPv4Addrs[j].NetworkAddress()
	if iAddr != jAddr {
		return iAddr < jAddr
	}

	// Sort smaller networks first
	return s.IPv4Addrs[i].Mask > s.IPv4Addrs[j].Mask
}

// FilterByTypeIPv4Addr filters SockAddrs and returns a list of IPv4Addrs
func (sas SockAddrs) FilterByTypeIPv4Addr() (ipv4addrs IPv4Addrs) {
	ipv4addrs = make(IPv4Addrs, 0, sas.Len())
	for _, sa := range sas {
		if sa.Type() == TypeIPv4 {
			na, ok := sa.(IPv4Addr)
			if ok {
				ipv4addrs = append(ipv4addrs, na)
			}
		}
	}
	return ipv4addrs
}
