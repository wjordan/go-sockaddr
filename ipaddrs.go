package sockaddr

type IPAddrs []IPAddr

func (s IPAddrs) Len() int      { return len(s) }
func (s IPAddrs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort IPAddrs by smallest network/most specific to largest network
type SortIPAddrsBySpecificMaskLen struct{ IPAddrs }

func (s SortIPAddrsBySpecificMaskLen) Less(i, j int) bool {
	return s.IPAddrs[i].Maskbits() > s.IPAddrs[j].Maskbits()
}

// Sort by largest to smallest network size
type SortIPAddrsByBroadMaskLen struct{ IPAddrs }

func (s SortIPAddrsByBroadMaskLen) Less(i, j int) bool {
	return s.IPAddrs[i].Maskbits() < s.IPAddrs[j].Maskbits()
}
