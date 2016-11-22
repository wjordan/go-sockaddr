package sockaddr

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strconv"
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

// DescIfAddress is identical to AscIfAddress but reverse ordered.
func DescIfAddress(p1Ptr, p2Ptr *IfAddr) int {
	return -1 * AscAddress(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// DescIfName is identical to AscIfName but reverse ordered.
func DescIfName(p1Ptr, p2Ptr *IfAddr) int {
	return -1 * strings.Compare(p1Ptr.Name, p2Ptr.Name)
}

// DescIfNetworkSize is identical to AscIfNetworkSize but reverse ordered.
func DescIfNetworkSize(p1Ptr, p2Ptr *IfAddr) int {
	return -1 * AscNetworkSize(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// DescIfPort is identical to AscIfPort but reverse ordered.
func DescIfPort(p1Ptr, p2Ptr *IfAddr) int {
	return -1 * AscPort(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// DescIfPrivate is identical to AscIfPrivate but reverse ordered.
func DescIfPrivate(p1Ptr, p2Ptr *IfAddr) int {
	return -1 * AscPrivate(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// DescIfType is identical to AscIfType but reverse ordered.
func DescIfType(p1Ptr, p2Ptr *IfAddr) int {
	return -1 * AscType(&p1Ptr.SockAddr, &p2Ptr.SockAddr)
}

// IfAttr forwards the selector to IfAttr.Attr() for resolution
func IfAttr(selectorName string, ifAddrs IfAddrs) string {
	if len(ifAddrs) == 0 {
		return ""
	}

	attrName := AttrName(strings.ToLower(selectorName))
	return ifAddrs[0].Attr(attrName)
}

// FilterIfByType filters IfAddrs and returns a list of the matching type
func FilterIfByType(ifAddrs IfAddrs, type_ SockAddrType) (matchedIfs, excludedIfs IfAddrs) {
	excludedIfs = make(IfAddrs, 0, len(ifAddrs))
	matchedIfs = make(IfAddrs, 0, len(ifAddrs))

	for _, ifAddr := range ifAddrs {
		if ifAddr.SockAddr.Type()&type_ != 0 {
			matchedIfs = append(matchedIfs, ifAddr)
		} else {
			excludedIfs = append(excludedIfs, ifAddr)
		}
	}
	return matchedIfs, excludedIfs
}

// GetAllInterfaces iterates over all available network interfaces and finds all
// available IP addresses on each interface and converts them to
// sockaddr.IPAddrs, and returning the result as an array of IfAddr.
func GetAllInterfaces() (IfAddrs, error) {
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

// GetDefaultInterfaces returns IfAddrs of the addresses attached to the default
// route.
func GetDefaultInterfaces() (IfAddrs, error) {
	defaultIfName, err := getDefaultIfName()
	if err != nil {
		return nil, err
	}

	var ifs IfAddrs
	ifAddrs, err := GetAllInterfaces()
	for _, ifAddr := range ifAddrs {
		if ifAddr.Name == defaultIfName {
			ifs = append(ifs, ifAddr)
		}
	}

	return ifs, nil
}

// GetPrivateIP returns a string with a single IP address that is part of RFC
// 6890 and has a default route.  If the system can't determine its IP address
// or find an RFC 6890 IP address, an empty string will be returned instead.
// This function is the `eval` equivilant of:
//
// ```
// $ sockaddr eval -raw '{{GetPrivateInterfaces | limit 1 | join "address" " "}}'
/// ```
func GetPrivateIP() (string, error) {
	privateIfs, err := GetPrivateInterfaces()
	if err != nil {
		return "", err
	}
	if len(privateIfs) < 1 {
		return "", nil
	}

	ifAddr := privateIfs[0]
	ip := *ToIPAddr(ifAddr.SockAddr)
	return ip.NetIP().String(), nil
}

// GetPrivateInterfaces returns an IfAddrs that is part of RFC 6890 and has a
// default route.  If the system can't determine its IP address or find an RFC
// 6890 IP address, an empty IfAddrs will be returned instead.  This function is
// the `eval` equivilant of:
//
// ```
// $ sockaddr eval -raw '{{GetDefaultInterfaces | sort "type,size" | include "RFC" "6890" | limit 1 | join "address" " "}}'
/// ```
func GetPrivateInterfaces() (IfAddrs, error) {
	privateIfs, err := GetDefaultInterfaces()
	if err != nil {
		return IfAddrs{}, err
	}
	if len(privateIfs) == 0 {
		return IfAddrs{}, nil
	}

	privateIfs, _ = FilterIfByType(privateIfs, TypeIP)
	if len(privateIfs) == 0 {
		return IfAddrs{}, nil
	}

	OrderedIfAddrBy(AscIfType, AscIfNetworkSize).Sort(privateIfs)

	privateIfs, _, err = IfByRFC(6890, privateIfs)
	if err != nil {
		return IfAddrs{}, err
	} else if len(privateIfs) == 0 {
		return IfAddrs{}, nil
	}

	return privateIfs, nil
}

// GetPublicInterfaces returns an IfAddrs that is NOT part of RFC 6890 and has a
// default route.  If the system can't determine its IP address or find a non
// RFC 6890 IP address, an empty IfAddrs will be returned instead.  This
// function is the `eval` equivilant of:
//
// ```
// $ sockaddr eval -raw '{{GetDefaultInterfaces | sort "type,size" | exclude "RFC" "6890" }}'
/// ```
func GetPublicInterfaces() (IfAddrs, error) {
	publicIfs, err := GetDefaultInterfaces()
	if err != nil {
		return IfAddrs{}, err
	}
	if len(publicIfs) == 0 {
		return IfAddrs{}, nil
	}

	publicIfs, _ = FilterIfByType(publicIfs, TypeIP)
	if len(publicIfs) == 0 {
		return IfAddrs{}, nil
	}

	OrderedIfAddrBy(AscIfType, AscIfNetworkSize).Sort(publicIfs)

	_, publicIfs, err = IfByRFC(6890, publicIfs)
	if err != nil {
		return IfAddrs{}, err
	} else if len(publicIfs) == 0 {
		return IfAddrs{}, nil
	}

	return publicIfs, nil
}

// GetPublicIP returns a string with a single IP address that is NOT part of RFC
// 6890 and has a default route.  If the system can't determine its IP address
// or find a non RFC 6890 IP address, an empty string will be returned instead.
// This function is the `eval` equivilant of:
//
// ```
// $ sockaddr eval -raw '{{GetPublicInterfaces | limit 1 | join "address" " "}}'
/// ```
func GetPublicIP() (string, error) {
	publicIfs, err := GetPublicInterfaces()
	if err != nil {
		return "", err
	} else if len(publicIfs) < 1 {
		return "", nil
	}

	ifAddr := publicIfs[0]
	ip := *ToIPAddr(ifAddr.SockAddr)
	return ip.NetIP().String(), nil
}

// IfByAddress returns a list of matched and non-matched IfAddrs, or an error if
// the regexp fails to compile.
func IfByAddress(inputRe string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) {
	re, err := regexp.Compile(inputRe)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to compile address regexp %+q: %v", inputRe, err)
	}

	matchedAddrs := make(IfAddrs, 0, len(ifAddrs))
	excludedAddrs := make(IfAddrs, 0, len(ifAddrs))
	for _, addr := range ifAddrs {
		if re.MatchString(addr.SockAddr.String()) {
			matchedAddrs = append(matchedAddrs, addr)
		} else {
			excludedAddrs = append(excludedAddrs, addr)
		}
	}

	return matchedAddrs, excludedAddrs, nil
}

// IfByName returns a list of matched and non-matched IfAddrs, or an error if
// the regexp fails to compile.
func IfByName(inputRe string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) {
	re, err := regexp.Compile(inputRe)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to compile name regexp %+q: %v", inputRe, err)
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

// IfByPort returns a list of matched and non-matched IfAddrs, or an error if
// the regexp fails to compile.
func IfByPort(inputRe string, ifAddrs IfAddrs) (matchedIfs, excludedIfs IfAddrs, err error) {
	re, err := regexp.Compile(inputRe)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to compile port regexp %+q: %v", inputRe, err)
	}

	ipIfs, nonIfs := FilterIfByType(ifAddrs, TypeIP)
	matchedIfs = make(IfAddrs, 0, len(ipIfs))
	excludedIfs = append(IfAddrs(nil), nonIfs...)
	for _, addr := range ipIfs {
		ipAddr := ToIPAddr(addr.SockAddr)
		if ipAddr == nil {
			continue
		}

		port := strconv.FormatInt(int64((*ipAddr).IPPort()), 10)
		if re.MatchString(port) {
			matchedIfs = append(matchedIfs, addr)
		} else {
			excludedIfs = append(excludedIfs, addr)
		}
	}

	return matchedIfs, excludedIfs, nil
}

// IfByRFC returns a list of matched and non-matched IfAddrs, or an error if the
// regexp fails to compile, that contain the relevant RFC-specified traits.  The
// most common RFC is RFC1918.
func IfByRFC(inputRFC uint, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) {
	matchedIfAddrs := make(IfAddrs, 0, len(ifAddrs))
	remainingIfAddrs := make(IfAddrs, 0, len(ifAddrs))

	rfcNetMap := KnownRFCs()
	rfcNets, ok := rfcNetMap[inputRFC]
	if !ok {
		return nil, nil, fmt.Errorf("unsupported RFC %d", inputRFC)
	}

	for _, ifAddr := range ifAddrs {
		var contained bool
		for _, rfcNet := range rfcNets {
			if rfcNet.Contains(ifAddr.SockAddr) {
				matchedIfAddrs = append(matchedIfAddrs, ifAddr)
				contained = true
				break
			}
		}
		if !contained {
			remainingIfAddrs = append(remainingIfAddrs, ifAddr)
		}
	}

	return matchedIfAddrs, remainingIfAddrs, nil
}

// IfByMaskSize returns a list of matched and non-matched IfAddrs that have the
// matching mask size.
func IfByMaskSize(maskSize uint, ifAddrs IfAddrs) (matchedIfs, excludedIfs IfAddrs, err error) {
	ipIfs, nonIfs := FilterIfByType(ifAddrs, TypeIP)
	matchedIfs = make(IfAddrs, 0, len(ipIfs))
	excludedIfs = append(IfAddrs(nil), nonIfs...)
	for _, addr := range ipIfs {
		ipAddr := ToIPAddr(addr.SockAddr)
		if ipAddr == nil {
			continue
		}

		if (*ipAddr).Maskbits() == int(maskSize) {
			matchedIfs = append(matchedIfs, addr)
		} else {
			excludedIfs = append(excludedIfs, addr)
		}
	}

	return matchedIfs, excludedIfs, nil
}

// IfByType returns a list of matching and non-matching IfAddr that match the
// specified type.  For instance:
//
// include "type" "IPv4,IPv6"
//
// will include any IfAddrs that is either an IPv4 or IPv6 address.  Any
// addresses on those interfaces that don't match will be included in the
// remainder results.
func IfByType(inputTypes string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) {
	matchingIfAddrs := make(IfAddrs, 0, len(ifAddrs))
	remainingIfAddrs := make(IfAddrs, 0, len(ifAddrs))

	ifTypes := strings.Split(strings.ToLower(inputTypes), "|")
	for _, ifType := range ifTypes {
		if ifType != "ip" && ifType != "ipv4" && ifType != "ipv6" && ifType != "unix" {
			return nil, nil, fmt.Errorf("unsupported type %q %q", ifType, inputTypes)
		}
	}

	for _, ifAddr := range ifAddrs {
		for _, ifType := range ifTypes {
			var matched bool
			switch {
			case ifType == "ip" && ifAddr.SockAddr.Type()&TypeIP != 0:
				matched = true
			case ifType == "ipv4" && ifAddr.SockAddr.Type()&TypeIPv4 != 0:
				matched = true
			case ifType == "ipv6" && ifAddr.SockAddr.Type()&TypeIPv6 != 0:
				matched = true
			case ifType == "unix" && ifAddr.SockAddr.Type()&TypeUnix != 0:
				matched = true
			}

			if matched {
				matchingIfAddrs = append(matchingIfAddrs, ifAddr)
			} else {
				remainingIfAddrs = append(remainingIfAddrs, ifAddr)
			}
		}
	}

	return matchingIfAddrs, remainingIfAddrs, nil
}

// IfByFlag returns a list of matching and non-matching IfAddrs that match the
// specified type.  For instance:
//
// include "flag" "up,broadcast"
//
// will include any IfAddrs that have both the "up" and "broadcast" flags set.
// Any addresses on those interfaces that don't match will be omitted from the
// results.
func IfByFlag(inputFlags string, ifAddrs IfAddrs) (matched, remainder IfAddrs, err error) {
	matchedAddrs := make(IfAddrs, 0, len(ifAddrs))
	excludedAddrs := make(IfAddrs, 0, len(ifAddrs))

	var wantForwardable,
		wantGlobalUnicast,
		wantInterfaceLocalMulticast,
		wantLinkLocalMulticast,
		wantLinkLocalUnicast,
		wantLoopback,
		wantMulticast,
		wantUnspecified bool
	var ifFlag, ifWant net.Flags
	for _, flagName := range strings.Split(strings.ToLower(inputFlags), "|") {
		switch flagName {
		case "broadcast":
			ifFlag, ifWant = ifFlag|net.FlagBroadcast, ifWant|net.FlagBroadcast
		case "down":
			ifFlag, ifWant = ifFlag|net.FlagUp, ifWant|0
		case "forwardable":
			wantForwardable = true
		case "global unicast":
			wantGlobalUnicast = true
		case "interface-local multicast":
			wantInterfaceLocalMulticast = true
		case "link-local multicast":
			wantLinkLocalMulticast = true
		case "link-local unicast":
			wantLinkLocalUnicast = true
		case "loopback":
			ifFlag, ifWant = ifFlag|net.FlagLoopback, ifWant|net.FlagLoopback
			wantLoopback = true
		case "multicast":
			ifFlag, ifWant = ifFlag|net.FlagMulticast, ifWant|net.FlagMulticast
			wantMulticast = true
		case "point-to-point":
			ifFlag, ifWant = ifFlag|net.FlagPointToPoint, ifWant|net.FlagPointToPoint
		case "unspecified":
			wantUnspecified = true
		case "up":
			ifFlag, ifWant = ifFlag|net.FlagUp, ifWant|net.FlagUp
		default:
			return nil, nil, fmt.Errorf("Unknown interface flag: %+q", flagName)
		}
	}

	for _, ifAddr := range ifAddrs {
		var matched bool
		if ifFlag != 0 && ifWant != 0 && ifAddr.Interface.Flags&ifFlag == ifWant {
			matched = true
		} else if ip := ToIPAddr(ifAddr.SockAddr); ip != nil {
			netIP := (*ip).NetIP()
			switch {
			case wantForwardable && !IsRFC(ForwardingBlacklist, ifAddr.SockAddr):
				matched = true
			case wantGlobalUnicast && netIP.IsGlobalUnicast():
				matched = true
			case wantInterfaceLocalMulticast && netIP.IsInterfaceLocalMulticast():
				matched = true
			case wantLinkLocalMulticast && netIP.IsLinkLocalMulticast():
				matched = true
			case wantLinkLocalUnicast && netIP.IsLinkLocalUnicast():
				matched = true
			case wantLoopback && netIP.IsLoopback():
				matched = true
			case wantMulticast && netIP.IsMulticast():
				matched = true
			case wantUnspecified && netIP.IsUnspecified():
				matched = true
			}
		}

		if matched {
			matchedAddrs = append(matchedAddrs, ifAddr)
		} else {
			excludedAddrs = append(excludedAddrs, ifAddr)
		}
	}
	return matchedAddrs, excludedAddrs, nil
}

// IncludeIfs returns an IfAddrs based on the passed in selector.
func IncludeIfs(selectorName, selectorParam string, inputIfAddrs IfAddrs) (IfAddrs, error) {
	var includedIfs IfAddrs
	var err error

	switch strings.ToLower(selectorName) {
	case "address":
		includedIfs, _, err = IfByAddress(selectorParam, inputIfAddrs)
	case "flag", "flags":
		includedIfs, _, err = IfByFlag(selectorParam, inputIfAddrs)
	case "name":
		includedIfs, _, err = IfByName(selectorParam, inputIfAddrs)
	case "port":
		includedIfs, _, err = IfByPort(selectorParam, inputIfAddrs)
	case "rfc", "rfcs":
		rfcs := strings.Split(selectorParam, "|")
		for _, rfcStr := range rfcs {
			rfc, err := strconv.ParseUint(rfcStr, 10, 64)
			if err != nil {
				continue
			}

			includedRFCIfs, _, err := IfByRFC(uint(rfc), inputIfAddrs)
			if err != nil {
				continue
			}
			includedIfs = append(includedIfs, includedRFCIfs...)
		}
	case "size":
		maskSize, err := strconv.ParseUint(selectorParam, 10, 64)
		if err != nil {
			return IfAddrs{}, fmt.Errorf("invalid include size argument (%q): %v", selectorParam, err)
		}
		includedIfs, _, err = IfByMaskSize(uint(maskSize), inputIfAddrs)
	case "type":
		includedIfs, _, err = IfByType(selectorParam, inputIfAddrs)
	default:
		return IfAddrs{}, fmt.Errorf("invalid include selector %q", selectorName)
	}

	if err != nil {
		return IfAddrs{}, err
	}

	return includedIfs, nil
}

// ExcludeIfs returns an IfAddrs based on the passed in selector.
func ExcludeIfs(selectorName, selectorParam string, inputIfAddrs IfAddrs) (IfAddrs, error) {
	var excludedIfs IfAddrs
	var err error

	switch strings.ToLower(selectorName) {
	case "address":
		_, excludedIfs, err = IfByAddress(selectorParam, inputIfAddrs)
	case "flag", "flags":
		_, excludedIfs, err = IfByFlag(selectorParam, inputIfAddrs)
	case "name":
		_, excludedIfs, err = IfByName(selectorParam, inputIfAddrs)
	case "port":
		_, excludedIfs, err = IfByPort(selectorParam, inputIfAddrs)
	case "rfc", "rfcs":
		rfcs := strings.Split(selectorParam, "|")
		for _, rfcStr := range rfcs {
			rfc, err := strconv.ParseUint(rfcStr, 10, 64)
			if err != nil {
				continue
			}

			_, excludedRFCIfs, err := IfByRFC(uint(rfc), inputIfAddrs)
			if err != nil {
				continue
			}
			excludedIfs = append(excludedIfs, excludedRFCIfs...)
		}
	case "size":
		maskSize, err := strconv.ParseUint(selectorParam, 10, 64)
		if err != nil {
			return IfAddrs{}, fmt.Errorf("invalid exclude size argument (%q): %v", selectorParam, err)
		}
		_, excludedIfs, err = IfByMaskSize(uint(maskSize), inputIfAddrs)
	case "type":
		_, excludedIfs, err = IfByType(selectorParam, inputIfAddrs)
	default:
		return IfAddrs{}, fmt.Errorf("invalid exclude selector %q", selectorName)
	}

	if err != nil {
		return IfAddrs{}, err
	}

	return excludedIfs, nil
}

// SortIfBy returns an IfAddrs sorted based on the passed in selector.  Multiple
// sort clauses can be passed in as a comma delimited list without whitespace.
func SortIfBy(selectorParam string, inputIfAddrs IfAddrs) (IfAddrs, error) {
	sortedIfs := append(IfAddrs(nil), inputIfAddrs...)

	clauses := strings.Split(selectorParam, ",")
	sortFuncs := make([]CmpIfAddrFunc, len(clauses))

	for i, clause := range clauses {
		switch strings.ToLower(clause) {
		case "address", "+address":
			// The "address" selector returns an array of IfAddrs
			// ordered by the network address.  IfAddrs that are not
			// comparable will be at the end of the list and in a
			// non-deterministic order.
			sortFuncs[i] = AscIfAddress
		case "-address":
			sortFuncs[i] = DescIfAddress
		case "name", "+name":
			// The "name" selector returns an array of IfAddrs
			// ordered by the interface name.
			sortFuncs[i] = AscIfName
		case "-name":
			sortFuncs[i] = DescIfName
		case "port", "+port":
			// The "port" selector returns an array of IfAddrs
			// ordered by the port, if included in the IfAddr.
			// IfAddrs that are not comparable will be at the end of
			// the list and in a non-deterministic order.
			sortFuncs[i] = AscIfPort
		case "-port":
			sortFuncs[i] = DescIfPort
		case "private", "+private":
			// The "private" selector returns an array of IfAddrs
			// ordered by private addresses first.  IfAddrs that are
			// not comparable will be at the end of the list and in
			// a non-deterministic order.
			sortFuncs[i] = AscIfPrivate
		case "-private":
			sortFuncs[i] = DescIfPrivate
		case "size", "+size":
			// The "size" selector returns an array of IfAddrs
			// ordered by the size of the network mask, largest mask
			// (larger number of hosts per network) to smallest
			// (e.g. a /24 sorts before a /32).
			sortFuncs[i] = AscIfNetworkSize
		case "-size":
			sortFuncs[i] = DescIfNetworkSize
		case "type", "+type":
			// The "type" selector returns an array of IfAddrs
			// ordered by the type of the IfAddr.  The sort order is
			// Unix, IPv4, then IPv6.
			sortFuncs[i] = AscIfType
		case "-type":
			sortFuncs[i] = DescIfType
		default:
			// Return an empty list for invalid sort types.
			return IfAddrs{}, fmt.Errorf("unknown sort type: %q", clause)
		}
	}

	OrderedIfAddrBy(sortFuncs...).Sort(sortedIfs)

	return sortedIfs, nil
}

// UniqueIfAddrsBy creates a unique set of IfAddrs based on the matching
// selector.  UniqueIfAddrsBy assumes the input has already been sorted.
func UniqueIfAddrsBy(selectorName string, inputIfAddrs IfAddrs) (IfAddrs, error) {
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

	return ifs, nil
}

// JoinIfAddrs joins an IfAddrs and returns a string
func JoinIfAddrs(selectorName string, joinStr string, inputIfAddrs IfAddrs) (string, error) {
	outputs := make([]string, 0, len(inputIfAddrs))
	attrName := AttrName(strings.ToLower(selectorName))

	for _, ifAddr := range inputIfAddrs {
		outputs = append(outputs, ifAddr.Attr(attrName))
	}
	return strings.Join(outputs, joinStr), nil
}

// LimitIfAddrs returns a slice of IfAddrs based on the specified limit.
func LimitIfAddrs(lim uint, in IfAddrs) (IfAddrs, error) {
	// Clamp the limit to the length of the array
	if int(lim) > len(in) {
		lim = uint(len(in))
	}

	return in[0:lim], nil
}

// OffsetIfAddrs returns a slice of IfAddrs based on the specified offset.
func OffsetIfAddrs(off int, in IfAddrs) (IfAddrs, error) {
	var end bool
	if off < 0 {
		end = true
		off = off * -1
	}

	if off > len(in) {
		return IfAddrs{}, fmt.Errorf("unable to seek past the end of the interface array: offset (%d) exceeds the number of interfaces (%d)", off, len(in))
	}

	if end {
		return in[len(in)-off : len(in)], nil
	}
	return in[off:len(in)], nil
}

func (ifAddr IfAddr) String() string {
	return fmt.Sprintf("%s %v", ifAddr.SockAddr, ifAddr.Interface)
}

// SortIfByAddr returns an IfAddrs ordered by address.  IfAddr that are not
// comparable will be at the end of the list, however their order is
// non-deterministic.
func SortIfByAddr(inputIfAddrs IfAddrs) (IfAddrs, error) {
	sortedIfAddrs := append(IfAddrs(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfAddress).Sort(sortedIfAddrs)
	return sortedIfAddrs, nil
}

// SortIfByName returns an IfAddrs ordered by interface name.  IfAddr that are
// not comparable will be at the end of the list, however their order is
// non-deterministic.
func SortIfByName(inputIfAddrs IfAddrs) (IfAddrs, error) {
	sortedIfAddrs := append(IfAddrs(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfName).Sort(sortedIfAddrs)
	return sortedIfAddrs, nil
}

// SortIfByPort returns an IfAddrs ordered by their port number, if set.
// IfAddrs that don't have a port set and are therefore not comparable will be
// at the end of the list (note: the sort order is non-deterministic).
func SortIfByPort(inputIfAddrs IfAddrs) (IfAddrs, error) {
	sortedIfAddrs := append(IfAddrs(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfPort).Sort(sortedIfAddrs)
	return sortedIfAddrs, nil
}

// SortIfByType returns an IfAddrs ordered by their type.  IfAddr that share a
// type are non-deterministic re: their sort order.
func SortIfByType(inputIfAddrs IfAddrs) (IfAddrs, error) {
	sortedIfAddrs := append(IfAddrs(nil), inputIfAddrs...)
	OrderedIfAddrBy(AscIfType).Sort(sortedIfAddrs)
	return sortedIfAddrs, nil
}

// parseDefaultIfNameFromRoute parses standard route(8)'s output for the *BSDs
// and Solaris.
func parseDefaultIfNameFromRoute(routeOut string) (string, error) {
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

// parseDefaultIfNameFromIPCmd parses the default interface from ip(8) for
// Linux.
func parseDefaultIfNameFromIPCmd(routeOut string) (string, error) {
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
