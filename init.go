package sockaddr

var rfc1918Networks []IPv4Addr
var rfc6598Networks []IPv4Addr

func init() {
	{ // initialize rfc1918
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

	{ // initialize rfc6598
		RFC6598Nets := []string{
			"100.64.0.0/10",
		}
		rfc6598Networks = make([]IPv4Addr, 0, len(RFC6598Nets))
		for _, network := range RFC6598Nets {
			sa, err := NewIPv4Addr(network)
			if err != nil {
				panic("Invalid RFC6598 network")
			}
			rfc6598Networks = append(rfc6598Networks, sa)
		}
	}
}
