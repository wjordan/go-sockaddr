package sockaddr

var rfc1918Networks []IPv4Addr
var rfc4193Networks []IPv6Addr
var rfc5735Networks []IPv4Addr
var rfc6598Networks []IPv4Addr
var rfc6890Networks []IPAddr
var isRFCMap map[uint]func(SockAddr) bool
var rfcNetMap map[uint][]SockAddr

func _IPv4(addr IPv4Addr, err error) IPv4Addr {
	return addr
}

func _IPv6(addr IPv6Addr, err error) IPv6Addr {
	return addr
}

func init() {
	rfcNetMap = map[uint][]SockAddr{
		1918: []SockAddr{
			_IPv4(NewIPv4Addr("10.0.0.0/8")),
			_IPv4(NewIPv4Addr("192.168.0.0/16")),
			_IPv4(NewIPv4Addr("172.16.0.0/12")),
		},
		4193: []SockAddr{
			_IPv6(NewIPv6Addr("fd00::/8")),
		},
		5735: []SockAddr{
			_IPv4(NewIPv4Addr("192.0.2.0/24")),    // TEST-NET-1
			_IPv4(NewIPv4Addr("198.51.100.0/24")), // TEST-NET-2
			_IPv4(NewIPv4Addr("203.0.113.0/24")),  // TEST-NET-3
			_IPv4(NewIPv4Addr("198.18.0.0/15")),   // Benchmarks
		},
		6598: []SockAddr{
			_IPv4(NewIPv4Addr("100.64.0.0/10")),
		},
		6890: []SockAddr{
			_IPv4(NewIPv4Addr("0.0.0.0/8")),          // [RFC1122] §3.2.1.3 "This host on this network".
			_IPv4(NewIPv4Addr("10.0.0.0/8")),         // [RFC1918]          Private-Use.
			_IPv4(NewIPv4Addr("100.64.0.0/10")),      // [RFC6598]          Shared Address Space.
			_IPv4(NewIPv4Addr("127.0.0.0/8")),        // [RFC1122] §3.2.1.3 Loopback.
			_IPv4(NewIPv4Addr("169.254.0.0/16")),     // [RFC3927]          Link-Local.
			_IPv4(NewIPv4Addr("172.16.0.0/12")),      // [RFC1918]          Private-Use.
			_IPv4(NewIPv4Addr("192.0.0.0/24")),       // [RFC6890] §2.1     IETF Protocol Assignments.
			_IPv4(NewIPv4Addr("192.0.0.0/29")),       // [RFC6333]          DS-Lite.
			_IPv4(NewIPv4Addr("192.0.2.0/24")),       // [RFC5737]          Documentation (TEST-NET-1).
			_IPv4(NewIPv4Addr("192.88.99.0/24")),     // [RFC3068]          6to4 Relay Anycast.
			_IPv4(NewIPv4Addr("198.18.0.0/15")),      // [RFC2544]          Benchmarking.
			_IPv4(NewIPv4Addr("192.168.0.0/16")),     // [RFC1918]          Private-Use.
			_IPv4(NewIPv4Addr("198.51.100.0/24")),    // [RFC5737]          Documentation (TEST-NET-2).
			_IPv4(NewIPv4Addr("203.0.113.0/24")),     // [RFC5737]          Documentation (TEST-NET-3).
			_IPv4(NewIPv4Addr("240.0.0.0/4")),        // [RFC1112] §4       Reserved.
			_IPv4(NewIPv4Addr("255.255.255.255/32")), // [RFC0919] §7       Limited Broadcast.

			_IPv6(NewIPv6Addr("::1/128")),       // [RFC4291] Loopback Address
			_IPv6(NewIPv6Addr("::/128")),        // [RFC4291] Unspecified Address
			_IPv6(NewIPv6Addr("64:ff9b::/96")),  // [RFC6052] IPv4-IPv6 Translat.
			_IPv6(NewIPv6Addr("::ffff:0:0/96")), // [RFC4291] IPv4-mapped Address
			_IPv6(NewIPv6Addr("100::/64")),      // [RFC6666] Discard-Only Address Block
			_IPv6(NewIPv6Addr("2001::/23")),     // [RFC2928] IETF Protocol Assignments
			_IPv6(NewIPv6Addr("2001::/32")),     // [RFC4380] TEREDO
			_IPv6(NewIPv6Addr("2001:2::/48")),   // [RFC5180] Benchmarking
			_IPv6(NewIPv6Addr("2001:db8::/32")), // [RFC3849] ocumentation
			_IPv6(NewIPv6Addr("2001:10::/28")),  // [RFC4843] ORCHID
			_IPv6(NewIPv6Addr("2002::/16")),     // [RFC3056] 6to4
			_IPv6(NewIPv6Addr("fc00::/7")),      // [RFC4193] Unique-Local
			_IPv6(NewIPv6Addr("fe80::/10")),     // [RFC4291] Linked-Scoped Unicast
		},
	}

	{ // initialize RFC1918
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

	{ // initialize RFC4193
		RFC4193Nets := []string{
			"fd00::/8",
		}
		rfc4193Networks = make([]IPv6Addr, 0, len(RFC4193Nets))
		for _, network := range RFC4193Nets {
			sa, err := NewIPv6Addr(network)
			if err != nil {
				panic("Invalid RFC4193 network")
			}
			rfc4193Networks = append(rfc4193Networks, sa)
		}
	}

	{ // initialize RFC5735
		RFC5735Nets := []string{
			"192.0.2.0/24",    // TEST-NET-1
			"198.51.100.0/24", // TEST-NET-2
			"203.0.113.0/24",  // TEST-NET-3
			"198.18.0.0/15",   // Benchmarks
		}
		rfc5735Networks = make([]IPv4Addr, 0, len(RFC5735Nets))
		for _, network := range RFC5735Nets {
			sa, err := NewIPv4Addr(network)
			if err != nil {
				panic("Invalid RFC5735 network")
			}
			rfc5735Networks = append(rfc5735Networks, sa)
		}
	}

	{ // initialize RFC6598
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

	{ // initialize RFC6890
		RFC6890Nets := []string{
			"0.0.0.0/8",          // [RFC1122] §3.2.1.3 "This host on this network".
			"10.0.0.0/8",         // [RFC1918]          Private-Use.
			"100.64.0.0/10",      // [RFC6598]          Shared Address Space.
			"127.0.0.0/8",        // [RFC1122] §3.2.1.3 Loopback.
			"169.254.0.0/16",     // [RFC3927]          Link-Local.
			"172.16.0.0/12",      // [RFC1918]          Private-Use.
			"192.0.0.0/24",       // [RFC6890] §2.1     IETF Protocol Assignments.
			"192.0.0.0/29",       // [RFC6333]          DS-Lite.
			"192.0.2.0/24",       // [RFC5737]          Documentation (TEST-NET-1).
			"192.88.99.0/24",     // [RFC3068]          6to4 Relay Anycast.
			"198.18.0.0/15",      // [RFC2544]          Benchmarking.
			"192.168.0.0/16",     // [RFC1918]          Private-Use.
			"198.51.100.0/24",    // [RFC5737]          Documentation (TEST-NET-2).
			"203.0.113.0/24",     // [RFC5737]          Documentation (TEST-NET-3).
			"240.0.0.0/4",        // [RFC1112] §4       Reserved.
			"255.255.255.255/32", // [RFC0919] §7       Limited Broadcast.

			"::1/128",       // [RFC4291] Loopback Address
			"::/128",        // [RFC4291] Unspecified Address
			"64:ff9b::/96",  // [RFC6052] IPv4-IPv6 Translat.
			"::ffff:0:0/96", // [RFC4291] IPv4-mapped Address
			"100::/64",      // [RFC6666] Discard-Only Address Block
			"2001::/23",     // [RFC2928] IETF Protocol Assignments
			"2001::/32",     // [RFC4380] TEREDO
			"2001:2::/48",   // [RFC5180] Benchmarking
			"2001:db8::/32", // [RFC3849] ocumentation
			"2001:10::/28",  // [RFC4843] ORCHID
			"2002::/16",     // [RFC3056] 6to4
			"fc00::/7",      // [RFC4193] Unique-Local
			"fe80::/10",     // [RFC4291] Linked-Scoped Unicast
		}
		rfc6890Networks = make([]IPAddr, 0, len(RFC6890Nets))
		for _, network := range RFC6890Nets {
			sa, err := NewIPAddr(network)
			if err != nil {
				panic("Invalid RFC6890 network")
			}
			rfc6890Networks = append(rfc6890Networks, sa)
		}
	}
}
