package netaddr

type NetAddrs []*NetAddr

func (s NetAddrs) Len() int      { return len(s) }
func (s NetAddrs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort by smallest network/most specific to largest network
type BySpecificMaskLen struct{ NetAddrs }

func (s BySpecificMaskLen) Less(i, j int) bool {
	return s.NetAddrs[i].Maskbits() > s.NetAddrs[j].Maskbits()
}

// Sort by largest to smallest network size
type ByBroadMaskLen struct{ NetAddrs }

func (s ByBroadMaskLen) Less(i, j int) bool {
	return s.NetAddrs[i].Maskbits() < s.NetAddrs[j].Maskbits()
}

// Sort by Network address
type ByNetwork struct{ NetAddrs }

func (s ByNetwork) Less(i, j int) bool {
	iIpUint, _ := s.NetAddrs[i].ToUint32()
	jIpUint, _ := s.NetAddrs[i].ToUint32()

	return iIpUint < jIpUint
}
