package sockaddr

// IPv4Addrs is a collection of IPv4Addrs
type IPv4Addrs []IPv4Addr

func (s IPv4Addrs) Len() int      { return len(s) }
func (s IPv4Addrs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort IPv4Addrs by smallest network/most specific to largest network
type SortIPv4AddrsBySpecificMaskLen struct{ IPv4Addrs }

func (s SortIPv4AddrsBySpecificMaskLen) Less(i, j int) bool {
	return s.IPv4Addrs[i].Maskbits() > s.IPv4Addrs[j].Maskbits()
}

// Sort by largest to smallest network size
type SortIPv4AddrsByBroadMaskLen struct{ IPv4Addrs }

func (s SortIPv4AddrsByBroadMaskLen) Less(i, j int) bool {
	return s.IPv4Addrs[i].Maskbits() < s.IPv4Addrs[j].Maskbits()
}

// Sort by network address
type SortIPv4AddrsByNetwork struct{ IPv4Addrs }

func (s SortIPv4AddrsByNetwork) Less(i, j int) bool {
	return s.IPv4Addrs[i].ToUint32() < s.IPv4Addrs[j].ToUint32()
}

// FilterByTypeIPv4Addr filters Sockaddrs and returns a list of IPv4Addrs
func (sas Sockaddrs) FilterByTypeIPv4Addr() (ipv4addrs IPv4Addrs) {
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
