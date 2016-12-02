package sockaddr_test

import (
	"net"
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

func TestGetPublicIP(t *testing.T) {
	ip, err := sockaddr.GetPublicIP()
	if err != nil {
		t.Fatalf("unable to get a public IP")
	}

	if ip == "" {
		t.Skip("it's hard to test this reliably")
	}
}

func TestGetInterfaceIP(t *testing.T) {
	ip, err := sockaddr.GetInterfaceIP(`^.*[\d]$`)
	if err != nil {
		t.Fatalf("regexp failed: %v", err)
	}

	if ip == "" {
		t.Skip("it's hard to test this reliably")
	}
}

func TestIfAddrAttr(t *testing.T) {
	tests := []struct {
		name     string
		ifAddr   sockaddr.IfAddr
		attr     string
		expected string
	}{
		{
			name: "name",
			ifAddr: sockaddr.IfAddr{
				Interface: net.Interface{
					Name: "abc0",
				},
			},
			attr:     "name",
			expected: "abc0",
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d must have a name", i)
		}

		result, err := sockaddr.IfAttr(test.attr, sockaddr.IfAddrs{test.ifAddr})
		if err != nil {
			t.Errorf("failed to get attr %q from %v", test.name, test.ifAddr)
		}

		if result != test.expected {
			t.Errorf("unexpected result")
		}
	}

	// Test an empty array
	result, err := sockaddr.IfAttr("name", sockaddr.IfAddrs{})
	if err != nil {
		t.Error(`failed to get attr "name" from an empty array`)
	}

	if result != "" {
		t.Errorf("unexpected result")
	}
}
