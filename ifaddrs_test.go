package sockaddr_test

import (
	"net"
	"reflect"
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

// NOTE: A number of these code paths are exercised in template/ and
// cmd/sockaddr/.
//
// TODO(sean@): Add better coverage for filtering functions (e.g. ExcludeBy*,
// IncludeBy*).

// TestGetIfAddrs runs through the motions of calling sockaddr.GetIfAddrs(), but
// doesn't do much in the way of testing beyond verifying that `lo0` has a
// loopback address present.
func TestGetIfAddrs(t *testing.T) {
	ifAddrs, err := sockaddr.GetAllInterfaces()
	if err != nil {
		t.Fatalf("Unable to proceed: %v", err)
	}
	if len(ifAddrs) == 0 {
		t.Skip()
	}

	var loInt *sockaddr.IfAddr
	for _, ifAddr := range ifAddrs {
		val := sockaddr.IfAddrAttr(ifAddr, "name")
		if val == "" {
			t.Fatalf("name failed")
		} else if val == "lo0" {
			loInt = &ifAddr
			break
		}
	}

	if val := sockaddr.IfAddrAttr(*loInt, "flags"); val != "up|loopback|multicast" {
		t.Fatalf("expected different flags from lo0: %q", val)
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
		if loInt.String() == "::1" {
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
	ifAddrs, err := sockaddr.GetDefaultInterfaces()
	if err != nil {
		switch {
		case len(ifAddrs) == 0:
			t.Fatal(err)
		case ifAddrs[0].Flags&net.FlagUp == 0:
			// if the first IfAddr isn't up, skip.
			t.Skip(err)
		default:
			t.Fatal(err)
		}
	}
}

func TestGetPrivateIP(t *testing.T) {
	ip, err := sockaddr.GetPrivateIP()
	if err != nil {
		t.Fatalf("private IP failed: %v", err)
	}

	if len(ip) == 0 {
		t.Fatalf("no private IP found")
	}
}

func TestIfAddrAttrs(t *testing.T) {
	attrs := sockaddr.IfAddrAttrs()
	if len(attrs) != 2 {
		t.Fatalf("wrong number of attrs")
	}
}

func TestGetAllInterfaces(t *testing.T) {
	ifAddrs, err := sockaddr.GetAllInterfaces()
	if err != nil {
		t.Fatalf("unable to gather interfaces: %v", err)
	}

	initialLen := len(ifAddrs)
	if initialLen == 0 {
		t.Fatalf("no interfaces available")
	}

	ifAddrs, err = sockaddr.SortIfBy("name,type,port,size,address", ifAddrs)
	if err != nil {
		t.Fatalf("unable to initially sort address")
	}

	ascSorted, err := sockaddr.SortIfBy("name,type,port,size,address", ifAddrs)
	if err != nil {
		t.Fatalf("unable to asc sort address")
	}

	descSorted, err := sockaddr.SortIfBy("name,type,port,size,-address", ascSorted)
	if err != nil {
		t.Fatalf("unable to desc sort address")
	}

	if initialLen != len(ascSorted) && len(ascSorted) != len(descSorted) {
		t.Fatalf("wrong len")
	}

	for i := initialLen - 1; i >= 0; i-- {
		if !reflect.DeepEqual(descSorted[i], ifAddrs[i]) {
			t.Errorf("wrong sort order: %d %v %v", i, descSorted[i], ifAddrs[i])
		}
	}
}

func TestGetDefaultInterfaces(t *testing.T) {
	ifAddrs, err := sockaddr.GetDefaultInterfaces()
	if err != nil {
		t.Fatalf("unable to gather default interfaces: %v", err)
	}

	if len(ifAddrs) == 0 {
		t.Fatalf("no default interfaces available")
	}
}

func TestGetPrivateInterfaces(t *testing.T) {
	ifAddrs, err := sockaddr.GetPrivateInterfaces()
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if len(ifAddrs) == 0 {
		t.Skip("no public IPs found")
	}

	if len(ifAddrs[0].String()) == 0 {
		t.Fatalf("no string representation of private IP found")
	}
}

func TestGetPublicInterfaces(t *testing.T) {
	ifAddrs, err := sockaddr.GetPublicInterfaces()
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if len(ifAddrs) == 0 {
		t.Skip("no public IPs found")
	}
}

func TestNewIPAddr(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
		pass   bool
	}{
		{
			name:   "ipv4",
			input:  "1.2.3.4",
			output: "1.2.3.4",
			pass:   true,
		},
		{
			name:   "ipv6",
			input:  "::1",
			output: "::1",
			pass:   true,
		},
		{
			name:   "invalid",
			input:  "255.255.255.256",
			output: "",
			pass:   false,
		},
	}

	for _, test := range tests {
		ip, err := sockaddr.NewIPAddr(test.input)
		switch {
		case err == nil && test.pass,
			err != nil && !test.pass:

		default:
			t.Errorf("expected %s's success to be %t", test.input, test.pass)
		}

		if !test.pass {
			continue
		}

		ipStr := ip.String()
		if ipStr != test.output {
			t.Errorf("Expected %q to match %q", test.input, test.output, ipStr)
		}

	}
}

func TestIPAttrs(t *testing.T) {
	const expectedIPAttrs = 11
	ipAttrs := sockaddr.IPAttrs()
	if len(ipAttrs) != expectedIPAttrs {
		t.Fatalf("wrong number of args")
	}
}
