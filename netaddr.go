package netaddr

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type netAddrType int

const (
	IPv4len = 4
	IPv6len = 16
)

const (
	netAddrTypeIP netAddrType = iota
	netAddrTypeNet
)

type NetAddr struct {
	AddrType netAddrType
	Address  net.IP
	Network  *net.IPNet
}

func (na *NetAddr) BroadcastAddress() *NetAddr {
	mask := binary.BigEndian.Uint32(na.Network.Mask)
	ipuint, ok := na.ToUint32()
	if !ok {
		return nil
	}
	return newFromUint32(na, ipuint&mask|^mask)
}

func (na *NetAddr) FirstUsableAddress() *NetAddr {
	firstUsable := na.NetworkAddress()
	ipuint, ok := firstUsable.ToUint32()
	if !ok {
		return nil
	}
	// If /32, return the address itself. If /31 assume a P2P link and
	// return the lower address.
	if na.Maskbits() < 31 {
		ipuint += 1
	}
	return newFromUint32(na, ipuint)
}

func (na *NetAddr) LastUsableAddress() *NetAddr {
	firstUsable := na.BroadcastAddress()
	ipuint, ok := firstUsable.ToUint32()
	if !ok {
		return nil
	}
	// If /32, return the address itself. If /31 assume a P2P link and
	// return the upper address.
	if na.Maskbits() < 31 {
		ipuint -= 1
	}
	return newFromUint32(na, ipuint)
}

func (na *NetAddr) Maskbits() int {
	if na == nil {
		panic("nil NetAddr")
	}
	if na.Network == nil {
		panic("nil network")
	}
	if na.Network.Mask == nil {
		panic("nil network mask")
	}
	maskOnes, _ := na.Network.Mask.Size()
	return maskOnes
}

func (na *NetAddr) NetworkAddress() *NetAddr {
	mask := binary.BigEndian.Uint32(na.Network.Mask)
	ipuint, ok := na.ToUint32()
	if !ok {
		return nil
	}
	return newFromUint32(na, ipuint&mask)
}

func New(s string) (na *NetAddr, err error) {
	na = new(NetAddr)
	na.Address, na.Network, err = net.ParseCIDR(s)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "invalid CIDR address: ") {
			return nil, err
		}

		// A bare IPv4 address is canonically referred to as a /32,
		// try again.
		s = s + "/32"
		na.Address, na.Network, err = net.ParseCIDR(s)
		if err != nil {
			return nil, err
		}
	}

	return na, nil
}

// newFromUint32 copies na into a new NetAddr struct using ipuint as the
// address.
func newFromUint32(na *NetAddr, ipuint uint32) (ret *NetAddr) {
	ret = new(NetAddr)
	*ret = *na // POD data copied by default, nice
	ipv4addr := make(net.IP, IPv6len)
	binary.BigEndian.PutUint32(ipv4addr, ipuint)
	// Use IPv4() ctor since it initializes in the last four bytes of
	// net.IP which is sized to fit a v6 address.
	ret.Address = net.IPv4(ipv4addr[0], ipv4addr[1], ipv4addr[2], ipv4addr[3])
	return ret
}

func (na *NetAddr) String() string {
	mask := na.Maskbits()
	return fmt.Sprintf("%s/%d", na.Address.String(), mask)
}

func (na *NetAddr) ToBinString() (s string) {
	ipnet, ok := na.ToUint32()
	if ok {
		s = fmt.Sprintf("%032s", strconv.FormatUint(uint64(ipnet), 2))
	}
	return s
}

func (na *NetAddr) ToHexString() (s string) {
	ipnet, ok := na.ToUint32()
	if ok {
		s = fmt.Sprintf("%08s", strconv.FormatUint(uint64(ipnet), 16))
	}
	return s
}

func (na *NetAddr) ToUint32() (ipint uint32, ok bool) {
	ipv4 := na.Address.To4()
	if len(ipv4) == IPv4len {
		return binary.BigEndian.Uint32(na.Address.To4()), true
	}

	// if p4 := na.Address.To4(); len(p4) == IPv4len {
	// 	ipint = uint32(p4[0])<<24 | uint32(p4[1])<<16 | uint32(p4[2])<<8 | uint32(p4[3])
	// 	return ipint, true
	// }
	return 0, false
}
