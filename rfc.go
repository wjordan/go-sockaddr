package sockaddr

import "log"

var rfcNetMap map[uint]SockAddrs

func init() {
	rfcNetMap = KnownRFCs()
}

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

// KnownRFCs returns an initial set of known RFCs.
//
// NOTE (sean@): As this list evolves over time, please submit patches to keep
// this list current.  If something isn't right, inquire, as it may just be a
// bug on my part.  Some of the inclusions were based on my judgement as to what
// would be a useful value (e.g. RFC3330).
//
// Useful resources:
//
// * https://www.iana.org/assignments/ipv6-address-space/ipv6-address-space.xhtml
// * https://www.iana.org/assignments/ipv6-unicast-address-assignments/ipv6-unicast-address-assignments.xhtml
// * https://www.iana.org/assignments/ipv6-address-space/ipv6-address-space.xhtml
func KnownRFCs() map[uint]SockAddrs {
	// NOTE(sean@): One could argue - decisively, I might add - that this
	// list belongs a once.Do() and is only instantiated when needed.
	// Further, these multiple SockAddrs per RFC scream RADIX tree, but
	// `ENOTIME`.  Patches welcome.
	return map[uint]SockAddrs{
		1112: SockAddrs{
			// [RFC1112] Host Extensions for IP Multicasting
			mustIPv4Addr("224.0.0.0/4"),
		},
		1918: SockAddrs{
			// [RFC1918] Address Allocation for Private Internets
			mustIPv4Addr("10.0.0.0/8"),
			mustIPv4Addr("192.168.0.0/16"),
			mustIPv4Addr("172.16.0.0/12"),
		},
		2544: SockAddrs{
			// [RFC2544] Benchmarking Methodology for Network
			// Interconnect Devices
			mustIPv4Addr("198.18.0.0/15"),
		},
		2765: SockAddrs{
			// [RFC2765] Stateless IP/ICMP Translation Algorithm
			// (SIIT) (obsoleted by RFCs 6145, which itself was
			// later obsoleted by 7915).

			// §2.1 Addresses
			mustIPv6Addr("0:0:0:0:0:ffff:0:0/96"),
		},
		2928: SockAddrs{
			// [RFC2928] Initial IPv6 Sub-TLA ID Assignments
			mustIPv6Addr("2001::/16"), // Superblock
			//mustIPv6Addr("2001:0000::/23"), // IANA
			//mustIPv6Addr("2001:0200::/23"), // APNIC
			//mustIPv6Addr("2001:0400::/23"), // ARIN
			//mustIPv6Addr("2001:0600::/23"), // RIPE NCC
			//mustIPv6Addr("2001:0800::/23"), // (future assignment)
			// ...
			//mustIPv6Addr("2001:FE00::/23"), // (future assignment)
		},
		3056: SockAddrs{ // 6to4 address
			// [RFC3056] Connection of IPv6 Domains via IPv4 Clouds

			// §2 IPv6 Prefix Allocation
			mustIPv6Addr("2002::/16"),
		},
		3068: SockAddrs{
			// [RFC3068] An Anycast Prefix for 6to4 Relay Routers
			// (obsolete by RFC7526)

			// §6to4 Relay anycast address
			mustIPv4Addr("192.88.99.0/24"),

			// §2.5 6to4 IPv6 relay anycast address
			//
			// NOTE: /120 == 128-(32-24)
			mustIPv6Addr("2002:c058:6301::/120"),
		},
		3171: SockAddrs{
			// [RFC3171] IANA Guidelines for IPv4 Multicast Address Assignments
			mustIPv4Addr("224.0.0.0/4"),
		},
		3330: SockAddrs{
			// [RFC3330] Special-Use IPv4 Addresses

			// Addresses in this block refer to source hosts on
			// "this" network.  Address 0.0.0.0/32 may be used as a
			// source address for this host on this network; other
			// addresses within 0.0.0.0/8 may be used to refer to
			// specified hosts on this network [RFC1700, page 4].
			mustIPv4Addr("0.0.0.0/8"),

			// 10.0.0.0/8 - This block is set aside for use in
			// private networks.  Its intended use is documented in
			// [RFC1918].  Addresses within this block should not
			// appear on the public Internet.
			mustIPv4Addr("10.0.0.0/8"),

			// 14.0.0.0/8 - This block is set aside for assignments
			// to the international system of Public Data Networks
			// [RFC1700, page 181]. The registry of assignments
			// within this block can be accessed from the "Public
			// Data Network Numbers" link on the web page at
			// http://www.iana.org/numbers.html.  Addresses within
			// this block are assigned to users and should be
			// treated as such.

			// 24.0.0.0/8 - This block was allocated in early 1996
			// for use in provisioning IP service over cable
			// television systems.  Although the IANA initially was
			// involved in making assignments to cable operators,
			// this responsibility was transferred to American
			// Registry for Internet Numbers (ARIN) in May 2001.
			// Addresses within this block are assigned in the
			// normal manner and should be treated as such.

			// 39.0.0.0/8 - This block was used in the "Class A
			// Subnet Experiment" that commenced in May 1995, as
			// documented in [RFC1797].  The experiment has been
			// completed and this block has been returned to the
			// pool of addresses reserved for future allocation or
			// assignment.  This block therefore no longer has a
			// special use and is subject to allocation to a
			// Regional Internet Registry for assignment in the
			// normal manner.

			// 127.0.0.0/8 - This block is assigned for use as the Internet host
			// loopback address.  A datagram sent by a higher level protocol to an
			// address anywhere within this block should loop back inside the host.
			// This is ordinarily implemented using only 127.0.0.1/32 for loopback,
			// but no addresses within this block should ever appear on any network
			// anywhere [RFC1700, page 5].
			mustIPv4Addr("127.0.0.0/8"),

			// 128.0.0.0/16 - This block, corresponding to the
			// numerically lowest of the former Class B addresses,
			// was initially and is still reserved by the IANA.
			// Given the present classless nature of the IP address
			// space, the basis for the reservation no longer
			// applies and addresses in this block are subject to
			// future allocation to a Regional Internet Registry for
			// assignment in the normal manner.

			// 169.254.0.0/16 - This is the "link local" block.  It
			// is allocated for communication between hosts on a
			// single link.  Hosts obtain these addresses by
			// auto-configuration, such as when a DHCP server may
			// not be found.
			mustIPv4Addr("169.254.0.0/16"),

			// 172.16.0.0/12 - This block is set aside for use in
			// private networks.  Its intended use is documented in
			// [RFC1918].  Addresses within this block should not
			// appear on the public Internet.
			mustIPv4Addr("172.16.0.0/12"),

			// 191.255.0.0/16 - This block, corresponding to the numerically highest
			// to the former Class B addresses, was initially and is still reserved
			// by the IANA.  Given the present classless nature of the IP address
			// space, the basis for the reservation no longer applies and addresses
			// in this block are subject to future allocation to a Regional Internet
			// Registry for assignment in the normal manner.

			// 192.0.0.0/24 - This block, corresponding to the
			// numerically lowest of the former Class C addresses,
			// was initially and is still reserved by the IANA.
			// Given the present classless nature of the IP address
			// space, the basis for the reservation no longer
			// applies and addresses in this block are subject to
			// future allocation to a Regional Internet Registry for
			// assignment in the normal manner.

			// 192.0.2.0/24 - This block is assigned as "TEST-NET" for use in
			// documentation and example code.  It is often used in conjunction with
			// domain names example.com or example.net in vendor and protocol
			// documentation.  Addresses within this block should not appear on the
			// public Internet.
			mustIPv4Addr("192.0.2.0/24"),

			// 192.88.99.0/24 - This block is allocated for use as 6to4 relay
			// anycast addresses, according to [RFC3068].
			mustIPv4Addr("192.88.99.0/24"),

			// 192.168.0.0/16 - This block is set aside for use in private networks.
			// Its intended use is documented in [RFC1918].  Addresses within this
			// block should not appear on the public Internet.
			mustIPv4Addr("192.168.0.0/16"),

			// 198.18.0.0/15 - This block has been allocated for use
			// in benchmark tests of network interconnect devices.
			// Its use is documented in [RFC2544].
			mustIPv4Addr("198.18.0.0/15"),

			// 223.255.255.0/24 - This block, corresponding to the
			// numerically highest of the former Class C addresses,
			// was initially and is still reserved by the IANA.
			// Given the present classless nature of the IP address
			// space, the basis for the reservation no longer
			// applies and addresses in this block are subject to
			// future allocation to a Regional Internet Registry for
			// assignment in the normal manner.

			// 224.0.0.0/4 - This block, formerly known as the Class
			// D address space, is allocated for use in IPv4
			// multicast address assignments.  The IANA guidelines
			// for assignments from this space are described in
			// [RFC3171].
			mustIPv4Addr("224.0.0.0/4"),

			// 240.0.0.0/4 - This block, formerly known as the Class E address
			// space, is reserved.  The "limited broadcast" destination address
			// 255.255.255.255 should never be forwarded outside the (sub-)net of
			// the source.  The remainder of this space is reserved
			// for future use.  [RFC1700, page 4]
			mustIPv4Addr("240.0.0.0/4"),
		},
		4038: SockAddrs{
			// RFC4038] Application Aspects of IPv6 Transition

			// §4.2. IPv6 Applications in a Dual-Stack Node
			mustIPv6Addr("0:0:0:0:0:ffff::/96"),
		},
		4193: SockAddrs{
			// [RFC4193] Unique Local IPv6 Unicast Addresses
			mustIPv6Addr("fc00::/7"),
		},
		4291: SockAddrs{
			// [RFC4291] IP Version 6 Addressing Architecture

			// §2.5.2 The Unspecified Address
			mustIPv6Addr("::/128"),

			// §2.5.3 The Loopback Address
			mustIPv6Addr("::1/128"),

			// §2.5.5.1.  IPv4-Compatible IPv6 Address
			mustIPv6Addr("0:0:0:0:0:0::/96"),

			// §2.5.5.2.  IPv4-Mapped IPv6 Address
			mustIPv6Addr("0:0:0:0:0:ffff::/96"),

			// §2.5.6 Link-Local IPv6 Unicast Addresses
			mustIPv6Addr("fe80::/10"),

			// §2.5.7 Site-Local IPv6 Unicast Addresses
			// (depreciated)
			mustIPv6Addr("fec0::/10"),

			// §2.7 Multicast Addresses
			mustIPv6Addr("ff00::/8"),

			// IPv6 Multicast Information.
			//
			// In the following "table" below, `ff0x` is replaced
			// with the following values depending on the scope of
			// the query:
			//
			// IPv6 Multicast Scopes:
			// * ff00/9 // reserved
			// * ff01/9 // interface-local
			// * ff02/9 // link-local
			// * ff03/9 // realm-local
			// * ff04/9 // admin-local
			// * ff05/9 // site-local
			// * ff08/9 // organization-local
			// * ff0e/9 // global
			// * ff0f/9 // reserved
			//
			// IPv6 Multicast Addresses:
			// * ff0x::2 // All routers
			// * ff02::5 // OSPFIGP
			// * ff02::6 // OSPFIGP Designated Routers
			// * ff02::9 // RIP Routers
			// * ff02::a // EIGRP Routers
			// * ff02::d // All PIM Routers
			// * ff02::1a // All RPL Routers
			// * ff0x::fb // mDNSv6
			// * ff0x::101 // All Network Time Protocol (NTP) servers
			// * ff02::1:1 // Link Name
			// * ff02::1:2 // All-dhcp-agents
			// * ff02::1:3 // Link-local Multicast Name Resolution
			// * ff05::1:3 // All-dhcp-servers
			// * ff02::1:ff00:0/104 // Solicited-node multicast address.
			// * ff02::2:ff00:0/104 // Node Information Queries
		},
		4773: SockAddrs{
			// [RFC4773] Administration of the IANA Special Purpose IPv6 Address Block
			mustIPv6Addr("2001:0000::/23"), // IANA
		},
		5180: SockAddrs{
			// [RFC5180] IPv6 Benchmarking Methodology for Network Interconnect Devices
			mustIPv6Addr("2001:2::/48"),
		},
		5735: SockAddrs{
			// [RFC5735] Special Use IPv4 Addresses
			mustIPv4Addr("192.0.2.0/24"),    // TEST-NET-1
			mustIPv4Addr("198.51.100.0/24"), // TEST-NET-2
			mustIPv4Addr("203.0.113.0/24"),  // TEST-NET-3
			mustIPv4Addr("198.18.0.0/15"),   // Benchmarks
		},
		6598: SockAddrs{
			// [RFC6598] IANA-Reserved IPv4 Prefix for Shared Address Space
			mustIPv4Addr("100.64.0.0/10"),
		},
		6666: SockAddrs{
			// [RFC6666] A Discard Prefix for IPv6
			mustIPv6Addr("0100::/64"),
		},
		6890: SockAddrs{
			// [RFC6890] Special-Purpose IP Address Registries
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
			mustIPv6Addr("2001::/16"),     // [RFC2928] IETF Protocol Assignments
			mustIPv6Addr("2002::/16"),     // [RFC3056] 6to4
			mustIPv6Addr("fc00::/7"),      // [RFC4193] Unique-Local
			mustIPv6Addr("fe80::/10"),     // [RFC4291] Linked-Scoped Unicast
		},
	}
}
