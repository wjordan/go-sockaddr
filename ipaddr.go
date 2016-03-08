package sockaddr

import (
	"fmt"
	"net"
)

// IPAddr is a generic IP address interface for IPv4 and IPv6 addresses,
// networks, and endpoints.
type IPAddr interface {
	Sockaddr
	Address() *net.IP
	FirstUsableAddress() IPAddr
	LastUsableAddress() IPAddr
	Maskbits() int
	Network() *net.IPNet
	NetworkPrefix() IPAddr
	Port() uint16
	SetPort(uint16)
	ToBinString() string
	ToHexString() string
}

// Maskbits returns the number of bits in the netmask of an IPAddr.  The
// domain of valid values for IPv4 are 0-32.  The domain for IPv6 is 0-128.
func ipMaskbits(ipn *net.IPNet) int {
	if ipn.Mask == nil {
		panic("nil network mask")
	}
	maskOnes, _ := ipn.Mask.Size()
	return maskOnes
}

// ipString returns a CIDR string using the unmasked address
func ipString(ipa IPAddr) string {
	var mask int = -1
	switch ipa.Type() {
	case TypeIPv4:
		mask = ipMaskbits(ipa.Network())
	case TypeIPv6:
		mask = ipMaskbits(ipa.Network())
	default:
		panic("Unknown type")
	}

	return fmt.Sprintf("%s/%d", ipa.Address().String(), mask)
}
