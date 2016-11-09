package sockaddr_test

import (
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

// TestGetIfAddrs runs through the motions of calling sockaddr.GetIfAddrs(), but
// doesn't do much in the way of testing beyond verifying that `lo0` has a
// loopback address present.
func TestGetIfAddrs(t *testing.T) {
	ifAddrs, err := sockaddr.GetIfSockAddrs()
	if err != nil {
		t.Fatalf("Unable to proceed: %v", err)
	}
	if len(ifAddrs) == 0 {
		t.Skip()
	}

	var loInt *sockaddr.IfAddr
	for _, ifAddr := range ifAddrs {
		if ifAddr.Name == "lo0" {
			loInt = &ifAddr
			break
		}
	}

	if loInt == nil {
		t.Fatalf("Expected to find an lo0 interface, didn't find any")
	}

	haveIPv4, foundIPv4lo := false, false
	haveIPv6, foundIPv6lo := false, false
	switch loInt.SockAddr.(type) {
	case sockaddr.IPv4Addr:
		haveIPv4 = true

		// Make the semi-brittle assumption that if we have
		// IPv4, we also have an address at 127.0.0.1 available
		// to us.
		if loInt.SockAddr.String() == "127.0.0.1/8" {
			foundIPv4lo = true
		}
	case sockaddr.IPv6Addr:
		haveIPv6 = true
		if loInt.String() == "100::" {
			foundIPv6lo = true
		}
	default:
		t.Fatalf("Unsupported type %v for address %v", loInt.Type(), loInt)
	}

	// While not wise, it's entirely possible a host doesn't have IPv4
	// enabled.
	if haveIPv4 && !foundIPv4lo {
		t.Fatalf("Had an IPv4 w/o an expected IPv4 loopback addresses")
	}

	// While prudent to run without, a sane environment may still contain an
	// IPv6 loopback address.
	if haveIPv6 && !foundIPv6lo {
		t.Fatalf("Had an IPv6 w/o an expected IPv6 loopback addresses")
	}
}

// TestGetDefaultIfName tests to make sure a default interface name is always
// returned from getDefaultIfName().
func TestGetDefaultInterface(t *testing.T) {
	// TODO: only run this test if we have Up interfaces

	_, err := sockaddr.GetDefaultInterfaces()
	if err != nil {
		t.Fatal(err)
	}
}
