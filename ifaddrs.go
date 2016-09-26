package sockaddr

import (
	"fmt"
	"net"
	"regexp"

	"github.com/hashicorp/errwrap"
)

// IfAddrs is a slice of IPAddrs for per interface
type IfAddrs struct {
	Addrs []SockAddr
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

		ipAddrs := make([]SockAddr, 0, len(addrs))
		for _, addr := range addrs {
			ipAddr, err := NewIPAddr(addr.String())
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

// IfByName returns a list of matched and non-matched IfAddrs, or an error if
// the regexp fails to compile.
func IfByName(inputRe string, ifAddrs []IfAddrs) (matched, remainder []IfAddrs, err error) {
	re, err := regexp.Compile(inputRe)
	if err != nil {
		return nil, nil, errwrap.Wrapf(fmt.Sprintf("Unable to compile *ByName regexp %+q: {{err}}", inputRe), err)
	}

	matchedAddrs := make([]IfAddrs, 0, len(ifAddrs))
	excludedAddrs := make([]IfAddrs, 0, len(ifAddrs))
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
func IfByNameExclude(inputRe string, ifAddrs []IfAddrs) ([]IfAddrs, error) {
	_, addrs, err := IfByName(inputRe, ifAddrs)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Unable to compile excludeByName regexp %+q: {{err}}", inputRe), err)
	}

	return addrs, nil
}

// IfByNameInclude includes any interface that matches the given regular
// expression (e.g. a regexp whitelist).
func IfByNameInclude(inputRe string, ifAddrs []IfAddrs) ([]IfAddrs, error) {
	addrs, _, err := IfByName(inputRe, ifAddrs)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Unable to compile includeByName regexp %+q: {{err}}", inputRe), err)
	}

	return addrs, nil
}

// IfByRFC returns a list of matched and non-matched IfAddrs, or an error if the
// regexp fails to compile, that contain the relevant RFC-specified traits.  The
// most common RFC is RFC1918.
func IfByRFC(inputRFC uint, ifAddrs []IfAddrs) (matched, remainder []IfAddrs, err error) {
	matchedIfAddrs := make([]IfAddrs, 0, len(ifAddrs))
	remainingIfAddrs := make([]IfAddrs, 0, len(ifAddrs))

	rfcNets, ok := rfcNetMap[inputRFC]
	if !ok {
		return nil, nil, fmt.Errorf("unsupported RFC %d", inputRFC)
	}

	for _, ifAddr := range ifAddrs {
		matchingAddrs := make([]SockAddr, 0, len(ifAddr.Addrs))
		remainingAddrs := make([]SockAddr, 0, len(ifAddr.Addrs))
		for _, sa := range ifAddr.Addrs {
			for _, rfcNet := range rfcNets {
				if rfcNet.Contains(sa) {
					matchingAddrs = append(matchingAddrs, sa)
				} else {
					remainingAddrs = append(remainingAddrs, sa)
				}
			}
		}

		if len(matchingAddrs) > 0 {
			matchedIfAddrs = append(matchedIfAddrs, IfAddrs{Addrs: matchingAddrs, Interface: ifAddr.Interface})
		}

		if len(remainingAddrs) > 0 {
			remainingIfAddrs = append(remainingIfAddrs, IfAddrs{Addrs: remainingAddrs, Interface: ifAddr.Interface})
		}
	}

	return matchedIfAddrs, remainingIfAddrs, nil
}

// IfByRFCExclude excludes any interface that matches the given regular
// expression (e.g. a regexp blacklist).
func IfByRFCExclude(inputRFC uint, ifAddrs []IfAddrs) ([]IfAddrs, error) {
	_, addrs, err := IfByRFC(inputRFC, ifAddrs)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

// IfByRFCInclude includes any interface that matches the given regular
// expression (e.g. a regexp whitelist).
func IfByRFCInclude(inputRFC uint, ifAddrs []IfAddrs) ([]IfAddrs, error) {
	addrs, _, err := IfByRFC(inputRFC, ifAddrs)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

// IfByType returns a list of matching and non-matching IfAddrs that match the
// specified type.  For instance:
//
// includeByType "^(IPv4|IPv6)$"
//
// will include any IfAddrs that contain at least one IPv4 or IPv6 address.  Any
// addresses on those interfaces that don't match will be omitted from the
// results.
func IfByType(inputRe string, ifAddrs []IfAddrs) (matched, remainder []IfAddrs, err error) {
	re, err := regexp.Compile(inputRe)
	if err != nil {
		return nil, nil, errwrap.Wrapf(fmt.Sprintf("Unable to compile includeByType regexp %+q: {{err}}", inputRe), err)
	}

	matchingIfAddrs := make([]IfAddrs, 0, len(ifAddrs))
	remainingIfAddrs := make([]IfAddrs, 0, len(ifAddrs))
	for _, ifAddr := range ifAddrs {
		matchingAddrs := make([]SockAddr, 0, len(ifAddr.Addrs))
		remainingAddrs := make([]SockAddr, 0, len(ifAddr.Addrs))
		for _, sockAddr := range ifAddr.Addrs {
			if re.MatchString(sockAddr.Type().String()) {
				matchingAddrs = append(matchingAddrs, sockAddr)
			} else {
				remainingAddrs = append(remainingAddrs, sockAddr)
			}
		}

		if len(matchingAddrs) > 0 {
			matchingIfAddrs = append(matchingIfAddrs, IfAddrs{Addrs: matchingAddrs, Interface: ifAddr.Interface})
		}

		if len(remainingAddrs) > 0 {
			remainingIfAddrs = append(remainingIfAddrs, IfAddrs{Addrs: remainingAddrs, Interface: ifAddr.Interface})
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
func IfByTypeExclude(inputRe string, ifAddrs []IfAddrs) ([]IfAddrs, error) {
	_, addrs, err := IfByType(inputRe, ifAddrs)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Unable to compile excludeByType regexp %+q: {{err}}", inputRe), err)
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
func IfByTypeInclude(inputRe string, ifAddrs []IfAddrs) ([]IfAddrs, error) {
	addrs, _, err := IfByType(inputRe, ifAddrs)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Unable to compile includeByType regexp %+q: {{err}}", inputRe), err)
	}

	return addrs, nil
}

// IfReturnAttrAddrs returns all of the matching SockAddr addresses found in the
// input as a flattened array.
func IfReturnAttrAddrs(inputIfAddrs []IfAddrs) SockAddrs {
	numAddrs := 0
	for _, addr := range inputIfAddrs {
		numAddrs += len(addr.Addrs)
	}

	addrs := make([]SockAddr, 0, numAddrs)
	for _, addr := range inputIfAddrs {
		addrs = append(addrs, addr.Addrs...)
	}

	return addrs
}

// IfReturnAttrNames returns all of the interface names in a flattened array.
func IfReturnAttrNames(inputIfAddrs []IfAddrs) []string {
	ifNames := make([]string, 0, len(inputIfAddrs))
	for _, ifAddr := range inputIfAddrs {
		ifNames = append(ifNames, ifAddr.Name)
	}

	return ifNames
}
