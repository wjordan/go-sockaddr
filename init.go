package sockaddr

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
}
