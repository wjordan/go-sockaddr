package sockaddr

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
)

// IfAddr is a combined SockAddr and Interface
type IfAddr struct {
	SockAddr
	net.Interface
}

//type IfAddrs []IfAddr

// CmpIfFunc is the function signature that must be met to be used in the
// OrderedIfAddrBy multiIfAddrSorter
type CmpIfAddrFunc func(p1, p2 *IfAddr) int

// multiIfAddrSorter implements the Sort interface, sorting the []IfAddr within.
type multiIfAddrSorter struct {
	ifAddrs []IfAddr
	cmp     []CmpIfAddrFunc
}

// Sort sorts the argument slice according to the Cmp functions passed to
// OrderedIfAddrBy.
func (ms *multiIfAddrSorter) Sort(ifAddrs []IfAddr) {
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

// AscIfAddress is a sorting function to sort []IfAddr by their respective
// address type.  Non-equal types are deferred in the sort.
func AscIfAddress(p1Ptr, p2Ptr *IfAddr) int {
	return AscAddress(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// AscIfName is a sorting function to sort []IfAddr by their interface names.
func AscIfName(p1Ptr, p2Ptr *IfAddr) int {
	return strings.Compare(p1Ptr.Name, p2Ptr.Name)
}

// AscIfPort is a sorting function to sort []IfAddr by their respective
// port type.  Non-equal types are deferred in the sort.
func AscIfPort(p1Ptr, p2Ptr *IfAddr) int {
	return AscPort(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// AscIfType is a sorting function to sort []IfAddr by their respective address
// type.  Non-equal types are deferred in the sort.
func AscIfType(p1Ptr, p2Ptr *IfAddr) int {
	return AscType(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// GetIfSockAddrs iterates over all available network interfaces and finds all
// available IP addresses on each interface and converts them to
// sockaddr.IPAddrs, and returning the result as an array of IfAddr.
func GetIfSockAddrs() ([]IfAddr, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ifAddrs := make([]IfAddr, 0, len(ifs))
	for _, intf := range ifs {
		addrs, err := intf.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ipAddr, err := NewIPAddr(addr.String())
			if err != nil {
				return []IfAddr{}, fmt.Errorf("unable to create an IP address from %q", addr.String())
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

// GetDefaultInterfaces returns []IfAddr of the addresses attached to the
// default route.
func GetDefaultInterfaces() ([]IfAddr, error) {
	defaultIfName, err := getDefaultIfName()
	if err != nil {
		return nil, err
	}

	var ifs []IfAddr
	ifAddrs, err := GetIfSockAddrs()
	for _, ifAddr := range ifAddrs {
		if ifAddr.Name == defaultIfName {
			ifs = append(ifs, ifAddr)
		}
	}

	return ifs, nil
}

// IfByName returns a list of matched and non-matched IfAddr, or an error if
// the regexp fails to compile.
func IfByName(inputRe string, ifAddrs []IfAddr) (matched, remainder []IfAddr, err error) {
	re, err := regexp.Compile(inputRe)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to compile *ByName regexp %+q: %v", inputRe, err)
	}

	matchedAddrs := make([]IfAddr, 0, len(ifAddrs))
	excludedAddrs := make([]IfAddr, 0, len(ifAddrs))
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
func IfByNameExclude(inputRe string, ifAddrs []IfAddr) ([]IfAddr, error) {
	_, addrs, err := IfByName(inputRe, ifAddrs)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile excludeByName regexp %+q: %v", inputRe, err)
	}

	return addrs, nil
}

// IfByNameInclude includes any interface that matches the given regular
// expression (e.g. a regexp whitelist).
func IfByNameInclude(inputRe string, ifAddrs []IfAddr) ([]IfAddr, error) {
	addrs, _, err := IfByName(inputRe, ifAddrs)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile includeByName regexp %+q: %v", inputRe, err)
	}

	return addrs, nil
}

// IfByRFC returns a list of matched and non-matched IfAddrs, or an error if the
// regexp fails to compile, that contain the relevant RFC-specified traits.  The
// most common RFC is RFC1918.
func IfByRFC(inputRFC uint, ifAddrs []IfAddr) (matched, remainder []IfAddr, err error) {
	matchedIfAddrs := make([]IfAddr, 0, len(ifAddrs))
	remainingIfAddrs := make([]IfAddr, 0, len(ifAddrs))

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
func IfByRFCExclude(inputRFC uint, ifAddrs []IfAddr) ([]IfAddr, error) {
	_, addrs, err := IfByRFC(inputRFC, ifAddrs)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

// IfByRFCInclude includes any interface that matches the given regular
// expression (e.g. a regexp whitelist).
func IfByRFCInclude(inputRFC uint, ifAddrs []IfAddr) ([]IfAddr, error) {
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
func IfByType(inputRe string, ifAddrs []IfAddr) (matched, remainder []IfAddr, err error) {
	re, err := regexp.Compile(inputRe)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to compile includeByType regexp %+q: %v", inputRe, err)
	}

	matchingIfAddrs := make([]IfAddr, 0, len(ifAddrs))
	remainingIfAddrs := make([]IfAddr, 0, len(ifAddrs))
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
func IfByTypeExclude(inputRe string, ifAddrs []IfAddr) ([]IfAddr, error) {
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
func IfByFlag(inputFlags string, ifAddrs []IfAddr) (matched, remainder []IfAddr, err error) {
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
func IfByFlagInclude(inputFlag string, ifAddrs []IfAddr) ([]IfAddr, error) {
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
func IfByFlagExclude(inputFlag string, ifAddrs []IfAddr) ([]IfAddr, error) {
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
func IfByTypeInclude(inputRe string, ifAddrs []IfAddr) ([]IfAddr, error) {
	addrs, _, err := IfByType(inputRe, ifAddrs)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile includeByType regexp %+q: %v", inputRe, err)
	}

	return addrs, nil
}

// GroupIfAddrsBy groups an []IfAddr based on the matching selector
func GroupIfAddrsBy(selectorName string, inputIfAddrs []IfAddr) []IfAddr {
	attrName := strings.ToLower(selectorName)

	hash := make(map[string][]IfAddr, len(inputIfAddrs))

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

		if _, found := hash[out]; found {
			hash[out] = append(hash[out], ifAddr)
		} else {
			hash[out] = []IfAddr{ifAddr}
		}
	}

	ifs := make([]IfAddr, 0, len(inputIfAddrs))
	for _, v := range hash {
		for _, ifAddr := range v {
			ifs = append(ifs, ifAddr)
		}
	}
	return ifs
}

// UniqueIfAddrsBy creates a unique set of []IfAddr based on the matching
// selector.  UniqueIfAddrsBy assumes the input has already been sorted.
func UniqueIfAddrsBy(selectorName string, inputIfAddrs []IfAddr) []IfAddr {
	attrName := strings.ToLower(selectorName)

	ifs := make([]IfAddr, 0, len(inputIfAddrs))
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

// JoinIfAddrs joins an []IfAddr and returns a string
func JoinIfAddrs(selectorName string, joinStr string, inputIfAddrs []IfAddr) string {
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

// LimitIfAddrs returns a slice of []IfAddrs based on the limitIfAddrs
func LimitIfAddrs(limitIfAddrs uint, inputIfAddrs []IfAddr) []IfAddr {
	// Clamp the limit to the length of the array
	if int(limitIfAddrs) > len(inputIfAddrs) {
		limitIfAddrs = uint(len(inputIfAddrs))
	}

	return inputIfAddrs[0:limitIfAddrs]
}

// ReverseIfAddrs reverses an []IfAddr.
func ReverseIfAddrs(inputIfAddrs []IfAddr) []IfAddr {
	reversedIfAddrs := append([]IfAddr(nil), inputIfAddrs...)
	for i := len(reversedIfAddrs)/2 - 1; i >= 0; i-- {
		opp := len(reversedIfAddrs) - 1 - i
		reversedIfAddrs[i], reversedIfAddrs[opp] = reversedIfAddrs[opp], reversedIfAddrs[i]
	}
	return reversedIfAddrs
}

func (ifAddr IfAddr) String() string {
	return fmt.Sprintf("%s %v", ifAddr.SockAddr, ifAddr.Interface)
}

// SortIfByAddr returns an []IfAddr ordered by address.  IfAddr that are not
// comparable will be at the end of the list, however their order is
// non-deterministic.
func SortIfByAddr(inputIfAddrs []IfAddr) []IfAddr {
	sortedIfAddrs := append([]IfAddr(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfAddress).Sort(sortedIfAddrs)
	return sortedIfAddrs
}

// SortIfByName returns an []IfAddr ordered by interface name.  IfAddr that are
// not comparable will be at the end of the list, however their order is
// non-deterministic.
func SortIfByName(inputIfAddrs []IfAddr) []IfAddr {
	sortedIfAddrs := append([]IfAddr(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfName).Sort(sortedIfAddrs)
	return sortedIfAddrs
}

// SortIfByPort returns an []IfAddr ordered by their port number, if set.
// IfAddrs that don't have a port set and are therefore not comparable will be
// at the end of the list (note: the sort order is non-deterministic).
func SortIfByPort(inputIfAddrs []IfAddr) []IfAddr {
	sortedIfAddrs := append([]IfAddr(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfPort).Sort(sortedIfAddrs)
	return sortedIfAddrs
}

// SortIfByType returns an []IfAddr ordered by their type.  IfAddr that share a
// type are non-deterministic re: their sort order.
func SortIfByType(inputIfAddrs []IfAddr) []IfAddr {
	sortedIfAddrs := append([]IfAddr(nil), inputIfAddrs...)
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
