package sockaddr

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IPv4Addr implements a convenience wrapper around the union of Go's
// built-in net.IP and net.IPNet types.  In UNIX-speak, IPv4Addr implements
// `sockaddr` when the the address family is set to AF_INET
// (i.e. `sockaddr_in`).
type IPv4Addr struct {
	IPAddr
	address net.IP
	network *net.IPNet
	IPPort  uint16
}

// Address returns the address as a net.IP (should be presized for IPv4).
func (ipv4 IPv4Addr) Address() *net.IP {
	return &ipv4.address
}

// BroadcastAddress is an IPv4Addr-only method that returns the broadcast
// address of the network (IPv6 only supports multicast).
func (ipv4 IPv4Addr) BroadcastAddress() IPv4Addr {
	mask := binary.BigEndian.Uint32(ipv4.network.Mask)
	ipuint := ipv4.ToUint32()
	return newFromUint32(ipv4, ipuint&mask|^mask)
}

// Equal returns true if a Sockaddr is equal to the receiving IPv4Addr.
func (ipv4a IPv4Addr) Equal(sa Sockaddr) bool {
	if sa.Type() != TypeIPv4 {
		return false
	}
	ipv4b, ok := sa.(IPv4Addr)
	if !ok {
		return false
	}

	// Now that the type conversion checks are complete, verify the data
	// is equal.
	if !ipv4a.address.Equal(ipv4b.address) {
		return false
	}

	if ipv4a.network.String() != ipv4b.network.String() {
		return false
	}

	if ipv4a.IPPort != ipv4b.IPPort {
		return false
	}

	return true
}

// FirstUsableAddress returns the first address following the network prefix.
// The first usable address in a network is normally the gateway and should
// not be used except by devices forwarding packets between two
// administratively distinct networks (i.e. a router).  This function does
// not discriminate against first usable vs "first address that should be
// used."  For example, IPv4Addr"192.168.1.10/24" would
func (ipv4 IPv4Addr) FirstUsableAddress() IPAddr {
	netPrefix := ipv4NetworkPrefix(ipv4)
	netPrefixUint := netPrefix.ToUint32()
	// If /32, return the address itself. If /31 assume a point-to-point
	// link and return the lower address.
	if ipv4.Maskbits() < 31 {
		netPrefixUint += 1
	}
	return newFromUint32(ipv4, netPrefixUint)
}

// ipv4NetworkPrefix is a pure, private helper function that calculates and
// returns the network prefix as an IPv4Addr.
func ipv4NetworkPrefix(ipv4 IPv4Addr) IPv4Addr {
	mask := binary.BigEndian.Uint32(ipv4.network.Mask)
	ipuint := ipv4.ToUint32()
	return newFromUint32(ipv4, ipuint&mask)
}

// LastUsableAddress returns the last address before the broadcast address in
// a given network.
func (ipv4 IPv4Addr) LastUsableAddress() IPAddr {
	netBroadcast := ipv4.BroadcastAddress()
	netBroadcastUint := netBroadcast.ToUint32()
	// If /32, return the address itself. If /31 assume a point-to-point
	// link and return the upper address.
	if ipv4.Maskbits() < 31 {
		netBroadcastUint -= 1
	}
	return newFromUint32(ipv4, netBroadcastUint)
}

// func (ipv4 *IPv4Addr) ListenArgs() (net, largs string) {
// 	return "tcp4", fmt.Sprintf("%s:%d", ipv4.address.String(), ipv4.IPPort)
// }

// Maskbits returns the number of network mask bits in a given IPv4Addr.  For
// example, the Maskbits() of "192.168.1.1/24" would return 24.
func (ipv4 IPv4Addr) Maskbits() int {
	return ipMaskbits(ipv4.network)
}

// Network returns a pointer to the net.IPNet within IPv4Addr receiver.
func (ipv4 IPv4Addr) Network() *net.IPNet {
	return ipv4.network
}

// NetworkPrefix returns the network prefix or network address for a given
// network.
func (ipv4 IPv4Addr) NetworkPrefix() IPAddr {
	mask := binary.BigEndian.Uint32(ipv4.network.Mask)
	ipuint := ipv4.ToUint32()
	return newFromUint32(ipv4, ipuint&mask)
}

// newIPv4FromNetIp creates an IPv4Addr from a net.IP byte array.
func newIPv4FromNetIp(ip *net.IP, net *net.IPNet) (ret IPv4Addr) {
	ret.address = *ip
	ret.network = net
	return ret
}

// newFromUint32 is a pure private helper function that copies ipv4 into a
// new IPv4Addr struct using ipuint as the address.
func newFromUint32(ipv4 IPv4Addr, ipuint uint32) (ret IPv4Addr) {
	ret = IPv4Addr{address: ipv4.address, network: ipv4.network, IPPort: ipv4.IPPort}
	ret.address = make(net.IP, IPv4len)
	binary.BigEndian.PutUint32(ret.address, ipuint)
	return ret
}

// Constructs a new IPv4Addr from string
func NewIPv4AddrFromString(s string) (ret IPv4Addr, err error) {
	// Before passing string to ParseCIDR(), test to see if there's a '/'
	// character in the last searchLen characters.  If not, assume a bare
	// IP and append "/32".
	const searchTerm = "/32"
	const searchLen = 3 // len(searchTerm)
	if len(s) > searchLen && strings.IndexByte(s[len(s)-searchLen:], '/') == -1 {
		s = s + searchTerm
	}

	addr, network, err := net.ParseCIDR(s)
	if err != nil {
		return IPv4Addr{}, err
	}

	// NOTE: ParseCIDR() can return true if the address was an IPv6
	// address, however we may have appended /32 to an IPv6.  For the
	// sake of cleanliness, return an error and let the caller handle
	// attempting to parse an IPv6 address.
	if addr.To4() == nil {
		return IPv4Addr{}, fmt.Errorf("Unable to convert %s to an IPv4 address", s)
	}

	return newIPv4FromNetIp(&addr, network), nil
}

// Port returns the configured port for an IPv4Addr
func (ipv4 IPv4Addr) Port() uint16 {
	return ipv4.IPPort
}

// SetPort is a setter method to set an IPv4Addr's port number
func (ipv4 IPv4Addr) SetPort(p uint16) {
	ipv4.IPPort = p
}

// ToBinString returns a string with the IPv4Addr address represented as a
// sequence of '0' and '1' characters.  This method is useful for debugging
// or by operators who want to inspect an address.
func (ipv4 IPv4Addr) ToBinString() string {
	ipv4addr := ipv4.ToUint32()
	return fmt.Sprintf("%032s", strconv.FormatUint(uint64(ipv4addr), 2))
}

// ToHexString returns a string with the IPv4Addr address represented as a
// sequence of hex characters.  This method is useful for debugging or by
// operators who want to inspect an address.
func (ipv4 IPv4Addr) ToHexString() string {
	ipnet := ipv4.ToUint32()
	return fmt.Sprintf("%08s", strconv.FormatUint(uint64(ipnet), 16))
}

// ToUint32 convers an IPv4Addr to a network ordered uint32
func (ipv4 *IPv4Addr) ToUint32() uint32 {
	packedIpAddr := ipv4.address.To4()
	return binary.BigEndian.Uint32(packedIpAddr)

	// if p4 := packedIpAddr.address.To4(); len(p4) == IPv4len {
	// 	ipint = uint32(p4[0])<<24 | uint32(p4[1])<<16 | uint32(p4[2])<<8 | uint32(p4[3])
	// 	return ipint
	// }
}

// Type is used as a type switch and returns TypeIPv4
func (IPv4Addr) Type() SockaddrType {
	return TypeIPv4
}
