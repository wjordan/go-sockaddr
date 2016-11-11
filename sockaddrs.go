package sockaddr

import (
	"bytes"
	"sort"
)

// SockAddrs is a collection of SockAddrs
type SockAddrs []SockAddr

func (s SockAddrs) Len() int      { return len(s) }
func (s SockAddrs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// CmpAddrFunc is the function signature that must be met to be used in the
// OrderedAddrBy multiAddrSorter
type CmpAddrFunc func(p1, p2 *SockAddr) int

// multiAddrSorter implements the Sort interface, sorting the SockAddrs within.
type multiAddrSorter struct {
	addrs SockAddrs
	cmp   []CmpAddrFunc
}

// Sort sorts the argument slice according to the Cmp functions passed to
// OrderedAddrBy.
func (ms *multiAddrSorter) Sort(sockAddrs SockAddrs) {
	ms.addrs = sockAddrs
	sort.Sort(ms)
}

// OrderedAddrBy sorts SockAddr by the list of sort function pointers.
func OrderedAddrBy(cmpFuncs ...CmpAddrFunc) *multiAddrSorter {
	return &multiAddrSorter{
		cmp: cmpFuncs,
	}
}

// Len is part of sort.Interface.
func (ms *multiAddrSorter) Len() int {
	return len(ms.addrs)
}

// Less is part of sort.Interface. It is implemented by looping along the
// Cmp() functions until it finds a comparison that is either less than,
// equal to, or greater than.
func (ms *multiAddrSorter) Less(i, j int) bool {
	p, q := &ms.addrs[i], &ms.addrs[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.cmp)-1; k++ {
		cmp := ms.cmp[k]
		x := cmp(p, q)
		switch x {
		case -1:
			// p < q, so we have a decision.
			return true
		case 1:
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever the
	// final comparison reports.
	switch ms.cmp[k](p, q) {
	case -1:
		return true
	case 1:
		return false
	default:
		// Still a tie! Now what?
		return false
	}
}

// Swap is part of sort.Interface.
func (ms *multiAddrSorter) Swap(i, j int) {
	ms.addrs[i], ms.addrs[j] = ms.addrs[j], ms.addrs[i]
}

// sortOrderDifferentTypes provides a constant to describe the strategy taken
// when two different types are compared.  The strategy is always to defer
// comparison and allow the next sort function attempt to sort.
//
// NOTE: This constant is here only for code readability.
const sortOrderDifferentTypes = 0

// AscAddress is a sorting function to sort SockAddrs by their respective
// address type.  Non-equal types are deferred in the sort.
func AscAddress(p1Ptr, p2Ptr *SockAddr) int {
	p1 := *p1Ptr
	p2 := *p2Ptr

	switch v := p1.(type) {
	case IPv4Addr:
		return v.CmpAddress(p2)
	case IPv6Addr:
		return v.CmpAddress(p2)
	}
	return sortOrderDifferentTypes
}

// AscPort is a sorting function to sort SockAddrs by their respective address
// type.  Non-equal types are deferred in the sort.
func AscPort(p1Ptr, p2Ptr *SockAddr) int {
	p1 := *p1Ptr
	p2 := *p2Ptr

	switch v := p1.(type) {
	case IPv4Addr:
		return v.CmpPort(p2)
	case IPv6Addr:
		return v.CmpPort(p2)
	}
	return sortOrderDifferentTypes
}

// AscPrivate is a sorting function to sort "more secure" private values before
// "more public" values.  Both IPv4 and IPv6 are compared against RFC6890
// (RFC6890 includes, and is not limited to, RFC1918 and RFC6598 for IPv4, and
// IPv6 includes RFC4193).
func AscPrivate(p1Ptr, p2Ptr *SockAddr) int {
	p1 := *p1Ptr
	p2 := *p2Ptr

	switch v := p1.(type) {
	case IPv4Addr, IPv6Addr:
		return v.CmpRFC(6890, p2)
	}
	return sortOrderDifferentTypes
}

// AscNetworkSize is a sorting function to sort SockAddrs based on their network
// size.  Non-equal types are deferred in the sort.
func AscNetworkSize(p1Ptr, p2Ptr *SockAddr) int {
	p1 := *p1Ptr
	p2 := *p2Ptr
	p1Type := p1.Type()
	p2Type := p2.Type()

	// Network size operations on non-IP types make no sense
	if p1Type != p2Type && p1Type != TypeIP {
		return sortOrderDifferentTypes
	}

	ipA := p1.(IPAddr)
	ipB := p2.(IPAddr)

	return bytes.Compare([]byte(*ipA.NetIPMask()), []byte(*ipB.NetIPMask()))
}

// AscType is a sorting function to sort "more secure" types before
// "less-secure" types.
func AscType(p1Ptr, p2Ptr *SockAddr) int {
	p1 := *p1Ptr
	p2 := *p2Ptr
	p1Type := p1.Type()
	p2Type := p2.Type()
	if p1Type < p2Type {
		return -1
	} else if p1Type == p2Type {
		return 0
	} else if p1Type > p2Type {
		return 1
	}
	return 0
}

// FilterByType filters SockAddrs and returns a list of the matching type
func (sas SockAddrs) FilterByType(type_ SockAddrType) SockAddrs {
	x := make(SockAddrs, 0, sas.Len())
	for _, sa := range sas {
		if sa.Type()&type_ != 0 {
			x = append(x, sa)
		}
	}
	return x
}

// OnlyIPv4 filters an array of SockAddrs and returns two slices: one of the
// IPv4Addr members of SockAddrs, and the list of non-IPv4Addrs.
func (sas SockAddrs) OnlyIPv4() (IPv4Addrs, SockAddrs) {
	ipv4Addrs := make(IPv4Addrs, 0, len(sas))
	nonIPv4Addrs := make(SockAddrs, 0, len(sas))

	for _, sa := range sas {
		switch v := sa.(type) {
		case IPv4Addr:
			ipv4Addrs = append(ipv4Addrs, v)
		default:
			nonIPv4Addrs = append(nonIPv4Addrs, sa)
		}
	}

	return ipv4Addrs, nonIPv4Addrs
}

// OnlyIPv6 filters an array of SockAddrs and returns two slices: one of the
// IPv6Addr members of SockAddrs, and the list of non-IPv6Addrs.
func (sas SockAddrs) OnlyIPv6() (IPv6Addrs, SockAddrs) {
	ipv6Addrs := make(IPv6Addrs, 0, len(sas))
	nonIPv6Addrs := make(SockAddrs, 0, len(sas))

	for _, sa := range sas {
		switch v := sa.(type) {
		case IPv6Addr:
			ipv6Addrs = append(ipv6Addrs, v)
		default:
			nonIPv6Addrs = append(nonIPv6Addrs, sa)
		}
	}

	return ipv6Addrs, nonIPv6Addrs
}

// SortByAddr returns an array of SockAddrs ordered by address.  SockAddrs that
// are not comparable will be at the end of the list, however their order is
// non-deterministic.
func SortByAddr(inputAddrs SockAddrs) SockAddrs {
	sortedAddrs := append([]SockAddr(nil), inputAddrs...)
	OrderedAddrBy(AscAddress).Sort(sortedAddrs)
	return sortedAddrs
}

// SortByPort returns an array of SockAddrs ordered by their port number, if
// set.  SockAddrs that don't have a port set and are therefore not comparable
// will be at the end of the list (note: the sort order is non-deterministic).
func SortByPort(inputAddrs SockAddrs) SockAddrs {
	sortedAddrs := append([]SockAddr(nil), inputAddrs...)
	OrderedAddrBy(AscPort).Sort(sortedAddrs)
	return sortedAddrs
}

// SortByType returns an array of SockAddrs ordered by their type.  SockAddrs
// that share a type are non-deterministic re: their sort order.
func SortByType(inputAddrs SockAddrs) SockAddrs {
	sortedAddrs := append([]SockAddr(nil), inputAddrs...)
	OrderedAddrBy(AscType).Sort(sortedAddrs)
	return sortedAddrs
}

// SortByNetworkSize returns an array of SockAddrs ordered by the size of the
// subnet (smallest to largest).
func SortByNetworkSize(inputAddrs SockAddrs) SockAddrs {
	sortedAddrs := append([]SockAddr(nil), inputAddrs...)
	OrderedAddrBy(AscNetworkSize).Sort(sortedAddrs)
	return sortedAddrs
}
