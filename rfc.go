package sockaddr

import "log"

// VisitAllRFCs iterates over all known RFCs and calls the visitor
func VisitAllRFCs(fn func(rfcNum uint, sockaddrs SockAddrs)) {
	for rfcNum, sas := range rfcNetMap {
		fn(rfcNum, sas)
	}
}

// mustIPv4Addr is a helper method that must return an IPv4Addr or panic on
// invalid input..
func mustIPv4Addr(addr string) IPv4Addr {
	ipv4, err := NewIPv4Addr(addr)
	if err != nil {
		log.Fatalf("Unable to create an IPv4Addr from %+q: %v", addr, err)
	}
	return ipv4
}

// mustIPv6Addr is a helper method that must return an IPv6Addr or panic on invalid
// input.
func mustIPv6Addr(addr string) IPv6Addr {
	ipv6, err := NewIPv6Addr(addr)
	if err != nil {
		log.Fatalf("Unable to create an IPv6Addr from %+q: %v", addr, err)
	}
	return ipv6
}

// KnownRFCs returns an initial set of known RFCs
func KnownRFCs() map[uint]SockAddrs {
	return map[uint]SockAddrs{
		1918: SockAddrs{
			mustIPv4Addr("10.0.0.0/8"),
			mustIPv4Addr("192.168.0.0/16"),
			mustIPv4Addr("172.16.0.0/12"),
		},
		4193: SockAddrs{
			mustIPv6Addr("fd00::/8"),
		},
		5735: SockAddrs{
			mustIPv4Addr("192.0.2.0/24"),    // TEST-NET-1
			mustIPv4Addr("198.51.100.0/24"), // TEST-NET-2
			mustIPv4Addr("203.0.113.0/24"),  // TEST-NET-3
			mustIPv4Addr("198.18.0.0/15"),   // Benchmarks
		},
		6598: SockAddrs{
			mustIPv4Addr("100.64.0.0/10"),
		},
		6890: SockAddrs{
			mustIPv4Addr("0.0.0.0/8"),          // [RFC1122] §3.2.1.3 "This host on this network".
			mustIPv4Addr("10.0.0.0/8"),         // [RFC1918]          Private-Use.
			mustIPv4Addr("100.64.0.0/10"),      // [RFC6598]          Shared Address Space.
			mustIPv4Addr("127.0.0.0/8"),        // [RFC1122] §3.2.1.3 Loopback.
			mustIPv4Addr("169.254.0.0/16"),     // [RFC3927]          Link-Local.
			mustIPv4Addr("172.16.0.0/12"),      // [RFC1918]          Private-Use.
			mustIPv4Addr("192.0.0.0/24"),       // [RFC6890] §2.1     IETF Protocol Assignments.
			mustIPv4Addr("192.0.0.0/29"),       // [RFC6333]          DS-Lite.
			mustIPv4Addr("192.0.2.0/24"),       // [RFC5737]          Documentation (TEST-NET-1).
			mustIPv4Addr("192.88.99.0/24"),     // [RFC3068]          6to4 Relay Anycast.
			mustIPv4Addr("198.18.0.0/15"),      // [RFC2544]          Benchmarking.
			mustIPv4Addr("192.168.0.0/16"),     // [RFC1918]          Private-Use.
			mustIPv4Addr("198.51.100.0/24"),    // [RFC5737]          Documentation (TEST-NET-2).
			mustIPv4Addr("203.0.113.0/24"),     // [RFC5737]          Documentation (TEST-NET-3).
			mustIPv4Addr("240.0.0.0/4"),        // [RFC1112] §4       Reserved.
			mustIPv4Addr("255.255.255.255/32"), // [RFC0919] §7       Limited Broadcast.

			mustIPv6Addr("::1/128"),       // [RFC4291] Loopback Address
			mustIPv6Addr("::/128"),        // [RFC4291] Unspecified Address
			mustIPv6Addr("64:ff9b::/96"),  // [RFC6052] IPv4-IPv6 Translat.
			mustIPv6Addr("::ffff:0:0/96"), // [RFC4291] IPv4-mapped Address
			mustIPv6Addr("100::/64"),      // [RFC6666] Discard-Only Address Block
			mustIPv6Addr("2001::/23"),     // [RFC2928] IETF Protocol Assignments
			mustIPv6Addr("2001::/32"),     // [RFC4380] TEREDO
			mustIPv6Addr("2001:2::/48"),   // [RFC5180] Benchmarking
			mustIPv6Addr("2001:db8::/32"), // [RFC3849] ocumentation
			mustIPv6Addr("2001:10::/28"),  // [RFC4843] ORCHID
			mustIPv6Addr("2002::/16"),     // [RFC3056] 6to4
			mustIPv6Addr("fc00::/7"),      // [RFC4193] Unique-Local
			mustIPv6Addr("fe80::/10"),     // [RFC4291] Linked-Scoped Unicast
		},
	}
}
