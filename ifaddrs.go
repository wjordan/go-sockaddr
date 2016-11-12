package sockaddr

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
)

// IfAddrs is a slice of IfAddr
type IfAddrs []IfAddr

func (ifs IfAddrs) Len() int { return len(ifs) }

// CmpIfFunc is the function signature that must be met to be used in the
// OrderedIfAddrBy multiIfAddrSorter
type CmpIfAddrFunc func(p1, p2 *IfAddr) int

// multiIfAddrSorter implements the Sort interface, sorting the IfAddrs within.
type multiIfAddrSorter struct {
	ifAddrs IfAddrs
	cmp     []CmpIfAddrFunc
}

// Sort sorts the argument slice according to the Cmp functions passed to
// OrderedIfAddrBy.
func (ms *multiIfAddrSorter) Sort(ifAddrs IfAddrs) {
	ms.ifAddrs = ifAddrs
	sort.Sort(ms)
}

// OrderedIfAddrBy sorts SockAddr by the list of sort function pointers.
func OrderedIfAddrBy(cmpFuncs ...CmpIfAddrFunc) *multiIfAddrSorter {
	return &multiIfAddrSorter{
		cmp: cmpFuncs,
	}
}

// Len is part of sort.Interface.
func (ms *multiIfAddrSorter) Len() int {
	return len(ms.ifAddrs)
}

// Less is part of sort.Interface. It is implemented by looping along the Cmp()
// functions until it finds a comparison that is either less than, equal to, or
// greater than.
func (ms *multiIfAddrSorter) Less(i, j int) bool {
	p, q := &ms.ifAddrs[i], &ms.ifAddrs[j]
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
		panic("undefined sort order for remaining items in the list")
	}
}

// Swap is part of sort.Interface.
func (ms *multiIfAddrSorter) Swap(i, j int) {
	ms.ifAddrs[i], ms.ifAddrs[j] = ms.ifAddrs[j], ms.ifAddrs[i]
}

// AscIfAddress is a sorting function to sort IfAddrs by their respective
// address type.  Non-equal types are deferred in the sort.
func AscIfAddress(p1Ptr, p2Ptr *IfAddr) int {
	return AscAddress(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// AscIfName is a sorting function to sort IfAddrs by their interface names.
func AscIfName(p1Ptr, p2Ptr *IfAddr) int {
	return strings.Compare(p1Ptr.Name, p2Ptr.Name)
}

// AscIfNetworkSize is a sorting function to sort IfAddrs by their respective
// network mask size.
func AscIfNetworkSize(p1Ptr, p2Ptr *IfAddr) int {
	return AscNetworkSize(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// AscIfPort is a sorting function to sort IfAddrs by their respective
// port type.  Non-equal types are deferred in the sort.
func AscIfPort(p1Ptr, p2Ptr *IfAddr) int {
	return AscPort(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// AscIfPrivate is a sorting function to sort IfAddrs by "private" values before
// "public" values.  Both IPv4 and IPv6 are compared against RFC6890 (RFC6890
// includes, and is not limited to, RFC1918 and RFC6598 for IPv4, and IPv6
// includes RFC4193).
func AscIfPrivate(p1Ptr, p2Ptr *IfAddr) int {
	return AscPrivate(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// AscIfType is a sorting function to sort IfAddrs by their respective address
// type.  Non-equal types are deferred in the sort.
func AscIfType(p1Ptr, p2Ptr *IfAddr) int {
	return AscType(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// GetIfSockAddrs iterates over all available network interfaces and finds all
// available IP addresses on each interface and converts them to
// sockaddr.IPAddrs, and returning the result as an array of IfAddr.
func GetIfSockAddrs() (IfAddrs, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ifAddrs := make(IfAddrs, 0, len(ifs))
	for _, intf := range ifs {
		addrs, err := intf.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ipAddr, err := NewIPAddr(addr.String())
			if err != nil {
				return IfAddrs{}, fmt.Errorf("unable to create an IP address from %q", addr.String())
			}

			ifAddr := IfAddr{
				SockAddr:  ipAddr,
				Interface: intf,
			}
			ifAddrs = append(ifAddrs, ifAddr)
		}
	}

	return ifAddrs, nil
}

// GetDefaultInterfaces returns IfAddrs of the addresses attached to the
// default route.
func GetDefaultInterfaces() (IfAddrs, error) {
	defaultIfName, err := getDefaultIfName()
	if err != nil {
		return nil, err
	}

	var ifs IfAddrs
	ifAddrs, err := GetIfSockAddrs()
	for _, ifAddr := range ifAddrs {
		if ifAddr.Name == defaultIfName {
			ifs = append(ifs, ifAddr)
		}
	}

	return ifs, nil
}

// IfByName returns a list of matched and non-matched IfAddrs, or an error if
// the regexp fails to compile.
func IfByName(inputRe string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) {
	re, err := regexp.Compile(inputRe)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to compile *ByName regexp %+q: %v", inputRe, err)
	}

	matchedAddrs := make(IfAddrs, 0, len(ifAddrs))
	excludedAddrs := make(IfAddrs, 0, len(ifAddrs))
	for _, addr := range ifAddrs {
		if re.MatchString(addr.Name) {
			matchedAddrs = append(matchedAddrs, addr)
		} else {
			excludedAddrs = append(excludedAddrs, addr)
		}
	}

	return matchedAddrs, excludedAddrs, nil
}

// IfByNameExclude excludes any interface that matches the given regular
// expression (e.g. a regexp blacklist).
func IfByNameExclude(inputRe string, ifAddrs IfAddrs) (IfAddrs, error) {
	_, addrs, err := IfByName(inputRe, ifAddrs)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile excludeByName regexp %+q: %v", inputRe, err)
	}

	return addrs, nil
}

// IfByNameInclude includes any interface that matches the given regular
// expression (e.g. a regexp whitelist).
func IfByNameInclude(inputRe string, ifAddrs IfAddrs) (IfAddrs, error) {
	addrs, _, err := IfByName(inputRe, ifAddrs)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile includeByName regexp %+q: %v", inputRe, err)
	}

	return addrs, nil
}

// IfByRFC returns a list of matched and non-matched IfAddrs, or an error if the
// regexp fails to compile, that contain the relevant RFC-specified traits.  The
// most common RFC is RFC1918.
func IfByRFC(inputRFC uint, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) {
	matchedIfAddrs := make(IfAddrs, 0, len(ifAddrs))
	remainingIfAddrs := make(IfAddrs, 0, len(ifAddrs))

	rfcNets, ok := rfcNetMap[inputRFC]
	if !ok {
		return nil, nil, fmt.Errorf("unsupported RFC %d", inputRFC)
	}

	for _, ifAddr := range ifAddrs {
		for _, rfcNet := range rfcNets {
			if rfcNet.Contains(ifAddr.SockAddr) {
				matchedIfAddrs = append(matchedIfAddrs, ifAddr)
			} else {
				remainingIfAddrs = append(remainingIfAddrs, ifAddr)
			}
		}
	}

	return matchedIfAddrs, remainingIfAddrs, nil
}

// IfByRFCExclude excludes any interface that matches the given regular
// expression (e.g. a regexp blacklist).
func IfByRFCExclude(inputRFC uint, ifAddrs IfAddrs) (IfAddrs, error) {
	_, addrs, err := IfByRFC(inputRFC, ifAddrs)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

// IfByRFCInclude includes any interface that matches the given regular
// expression (e.g. a regexp whitelist).
func IfByRFCInclude(inputRFC uint, ifAddrs IfAddrs) (IfAddrs, error) {
	addrs, _, err := IfByRFC(inputRFC, ifAddrs)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

// IfByType returns a list of matching and non-matching IfAddr that match the
// specified type.  For instance:
//
// includeByType "^(IPv4|IPv6)$"
//
// will include any IfAddrs that contain at least one IPv4 or IPv6 address.  Any
// addresses on those interfaces that don't match will be omitted from the
// results.
func IfByType(inputRe string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) {
	re, err := regexp.Compile(inputRe)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to compile includeByType regexp %+q: %v", inputRe, err)
	}

	matchingIfAddrs := make(IfAddrs, 0, len(ifAddrs))
	remainingIfAddrs := make(IfAddrs, 0, len(ifAddrs))
	for _, ifAddr := range ifAddrs {
		if re.MatchString(ifAddr.SockAddr.Type().String()) {
			matchingIfAddrs = append(matchingIfAddrs, ifAddr)
		} else {
			remainingIfAddrs = append(remainingIfAddrs, ifAddr)
		}
	}

	return matchingIfAddrs, remainingIfAddrs, nil
}

// IfByTypeExclude excludes any addresses that are not of the input type.  For
// instance:
//
// excludeByType "^(IPv6)$"
//
// will only include IfAddrs that have at least one non-IPv6 address.  Any
// addresses on those interfaces that don't match will be omitted from the
// results.
func IfByTypeExclude(inputRe string, ifAddrs IfAddrs) (IfAddrs, error) {
	_, addrs, err := IfByType(inputRe, ifAddrs)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile excludeByType regexp %+q: {{err}}", inputRe, err)
	}

	return addrs, nil
}

// IfByFlag returns a list of matching and non-matching IfAddrs that match the
// specified type.  For instance:
//
// includeByFlag "up broadcast"
//
// will include any IfAddrs that have both the "up" and "broadcast" flags set.
// Any addresses on those interfaces that don't match will be omitted from the
// results.
func IfByFlag(inputFlags string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) {
	return nil, nil, fmt.Errorf("Unable to compile includeByFlag regexp %+q: %v", inputFlags, err)
}

// IfByFlagInclude includes any interface and only the matching addresses that
// are of the input type.  For instance:
//
// includeByFlag "up broadcast"
//
// will include any IfAddrs that have the flag "up" and "broadcast" set.  Any
// addresses on those interfaces that don't match will be omitted from the
// results.
func IfByFlagInclude(inputFlag string, ifAddrs IfAddrs) (IfAddrs, error) {
	addrs, _, err := IfByFlag(inputFlag, ifAddrs)
	if err != nil {
		return nil, fmt.Errorf("Invalid flag in includeByFlag %+q: %v", inputFlag, err)
	}

	return addrs, nil
}

// IfByFlagExclude excludes any interfaces that don't have the appropriate flag
// set.  For instance:
//
// excludeByFlag "up"
//
// will only include IfAddrs that don't have the "Up" flag set.  Any addresses
// on those interfaces that don't match will be omitted from the results.
func IfByFlagExclude(inputFlag string, ifAddrs IfAddrs) (IfAddrs, error) {
	_, addrs, err := IfByFlag(inputFlag, ifAddrs)
	if err != nil {
		return nil, fmt.Errorf("Invalid flag in excludeByFlag %+q: %v", inputFlag, err)
	}

	return addrs, nil
}

// IfByTypeInclude includes any interface and only the matching addresses that
// are of the input type.  For instance:
//
// includeByType "^(IPv4|IPv6)$"
//
// will include any IfAddrs that contain at least one IPv4 or IPv6 address.  Any
// addresses on those interfaces that don't match will be omitted from the
// results.
func IfByTypeInclude(inputRe string, ifAddrs IfAddrs) (IfAddrs, error) {
	addrs, _, err := IfByType(inputRe, ifAddrs)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile includeByType regexp %+q: %v", inputRe, err)
	}

	return addrs, nil
}

// SortIfBy returns an IfAddrs sorted based on the passed in selector.
func SortIfBy(selectorName string, inputIfAddrs IfAddrs) IfAddrs {
	sortedIfs := append(IfAddrs(nil), inputIfAddrs...)
	switch strings.ToLower(selectorName) {
	case "address":
		// The "address" selector returns an array of IfAddrs ordered by
		// the network address.  IfAddrs that are not comparable will be
		// at the end of the list and in a non-deterministic order.
		OrderedIfAddrBy(AscIfAddress).Sort(sortedIfs)
	case "name":
		// The "name" selector returns an array of IfAddrs ordered by
		// the interface name.
		OrderedIfAddrBy(AscIfName).Sort(sortedIfs)
	case "port":
		// The "port" selector returns an array of IfAddrs ordered by
		// the port, if included in the IfAddr.  IfAddrs that are not
		// comparable will be at the end of the list and in a
		// non-deterministic order.
		OrderedIfAddrBy(AscIfPort).Sort(sortedIfs)
	case "private":
		// The "private" selector returns an array of IfAddrs ordered by
		// private addresses first.  IfAddrs that are not comparable
		// will be at the end of the list and in a non-deterministic
		// order.
		OrderedIfAddrBy(AscIfPrivate).Sort(sortedIfs)
	case "size":
		// The "size" selector returns an array of IfAddrs ordered by
		// the size of the network mask, smallest mask (fewest number of
		// hosts per network) to largest (e.g. a /32 sorts before a
		// /24).
		OrderedIfAddrBy(AscIfNetworkSize).Sort(sortedIfs)
	case "type":
		// The "type" selector returns an array of IfAddrs ordered by
		// the type of the IfAddr.  The sort order is Unix, IPv4, then
		// IPv6.
		OrderedIfAddrBy(AscIfType).Sort(sortedIfs)
	default:
		// Return an empty list for invalid sort types.
		return IfAddrs{}
	}

	return sortedIfs
}

// UniqueIfAddrsBy creates a unique set of IfAddrs based on the matching
// selector.  UniqueIfAddrsBy assumes the input has already been sorted.
func UniqueIfAddrsBy(selectorName string, inputIfAddrs IfAddrs) IfAddrs {
	attrName := strings.ToLower(selectorName)

	ifs := make(IfAddrs, 0, len(inputIfAddrs))
	var lastMatch string
	for _, ifAddr := range inputIfAddrs {
		var out string
		switch attrName {
		case "address":
			out = ifAddr.SockAddr.String()
		case "name":
			out = ifAddr.Name
		default:
			out = fmt.Sprintf("<unsupported method %+q>", selectorName)
		}

		switch {
		case lastMatch == "", lastMatch != out:
			lastMatch = out
			ifs = append(ifs, ifAddr)
		case lastMatch == out:
			continue
		}
	}

	return ifs
}

// JoinIfAddrs joins an IfAddrs and returns a string
func JoinIfAddrs(selectorName string, joinStr string, inputIfAddrs IfAddrs) string {
	attrName := strings.ToLower(selectorName)

	outputs := make([]string, 0, len(inputIfAddrs))

	for _, ifAddr := range inputIfAddrs {
		var out string
		switch attrName {
		case "address":
			out = ifAddr.SockAddr.String()
		case "name":
			out = ifAddr.Name
		default:
			out = fmt.Sprintf("<unsupported method %+q>", selectorName)
		}
		outputs = append(outputs, out)
	}
	return strings.Join(outputs, joinStr)
}

// LimitIfAddrs returns a slice of IfAddrs based on the specified limit.
func LimitIfAddrs(lim uint, in IfAddrs) IfAddrs {
	// Clamp the limit to the length of the array
	if int(lim) > len(in) {
		lim = uint(len(in))
	}

	return in[0:lim]
}

// OffsetIfAddrs returns a slice of IfAddrs based on the specified offset.
func OffsetIfAddrs(off int, in IfAddrs) IfAddrs {
	var end bool
	if off < 0 {
		end = true
		off = off * -1
	}

	if off > len(in) {
		return IfAddrs{}
	}

	if end {
		return in[len(in)-off : len(in)]
	}
	return in[off:len(in)]
}

// ReverseIfAddrs reverses an IfAddrs.
func ReverseIfAddrs(inputIfAddrs IfAddrs) IfAddrs {
	reversedIfAddrs := append(IfAddrs(nil), inputIfAddrs...)
	for i := len(reversedIfAddrs)/2 - 1; i >= 0; i-- {
		opp := len(reversedIfAddrs) - 1 - i
		reversedIfAddrs[i], reversedIfAddrs[opp] = reversedIfAddrs[opp], reversedIfAddrs[i]
	}
	return reversedIfAddrs
}

func (ifAddr IfAddr) String() string {
	return fmt.Sprintf("%s %v", ifAddr.SockAddr, ifAddr.Interface)
}

// SortIfByAddr returns an IfAddrs ordered by address.  IfAddr that are not
// comparable will be at the end of the list, however their order is
// non-deterministic.
func SortIfByAddr(inputIfAddrs IfAddrs) IfAddrs {
	sortedIfAddrs := append(IfAddrs(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfAddress).Sort(sortedIfAddrs)
	return sortedIfAddrs
}

// SortIfByName returns an IfAddrs ordered by interface name.  IfAddr that are
// not comparable will be at the end of the list, however their order is
// non-deterministic.
func SortIfByName(inputIfAddrs IfAddrs) IfAddrs {
	sortedIfAddrs := append(IfAddrs(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfName).Sort(sortedIfAddrs)
	return sortedIfAddrs
}

// SortIfByPort returns an IfAddrs ordered by their port number, if set.
// IfAddrs that don't have a port set and are therefore not comparable will be
// at the end of the list (note: the sort order is non-deterministic).
func SortIfByPort(inputIfAddrs IfAddrs) IfAddrs {
	sortedIfAddrs := append(IfAddrs(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfPort).Sort(sortedIfAddrs)
	return sortedIfAddrs
}

// SortIfByType returns an IfAddrs ordered by their type.  IfAddr that share a
// type are non-deterministic re: their sort order.
func SortIfByType(inputIfAddrs IfAddrs) IfAddrs {
	sortedIfAddrs := append(IfAddrs(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfType).Sort(sortedIfAddrs)
	return sortedIfAddrs
}

// parseBSDDefaultIfName is a *BSD-specific parsing function for route(8)'s
// output.
func parseBSDDefaultIfName(routeOut string) (string, error) {
	lines := strings.Split(routeOut, "\n")
	for _, line := range lines {
		kvs := strings.SplitN(line, ":", 2)
		if len(kvs) != 2 {
			continue
		}

		if strings.TrimSpace(kvs[0]) == "interface" {
			ifName := strings.TrimSpace(kvs[1])
			return ifName, nil
		}
	}

	return "", errors.New("No default interface found")
}

// parseLinuxDefaultIfName is a Linux-specific parsing function for route(8)'s
// output.
func parseLinuxDefaultIfName(routeOut string) (string, error) {
	lines := strings.Split(routeOut, "\n")
	for _, line := range lines {
		kvs := strings.SplitN(line, " ", 5)
		if len(kvs) != 5 {
			continue
		}

		if kvs[0] == "default" &&
			kvs[1] == "via" &&
			kvs[3] == "dev" {
			ifName := strings.TrimSpace(kvs[4])
			return ifName, nil
		}
	}

	return "", errors.New("No default interface found")
}
