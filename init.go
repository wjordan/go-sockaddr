package sockaddr

var rfc1918Networks []IPv4Addr

func init() {
	RFC1918Nets := []string{
		"10.0.0.0/8",
		"192.168.0.0/16",
		"172.16.0.0/12",
	}
	rfc1918Networks = make([]IPv4Addr, 0, len(RFC1918Nets))
	for _, network := range RFC1918Nets {
		sa, err := NewIPv4Addr(network)
		if err != nil {
			panic("Invalid RFC1918 network")
		}
		rfc1918Networks = append(rfc1918Networks, sa)
	}
}
