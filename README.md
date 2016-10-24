# go-sockaddr
Socket convenience functions for Go

If you're familiar with UNIX's `sockaddr` struct's, the following diagram
mapping the C `sockaddr` (top) to `go-sockaddr` structs (bottom) and
interfaces will be helpful:

```
+-------------------------------------------------------+
|                                                       |
|                        sockaddr                       |
|                        SockAddr                       |
|                                                       |
| +--------------+ +----------------------------------+ |
| | sockaddr_un  | |                                  | |
| | SockAddrUnix | |           sockaddr_in{,6}        | |
| +--------------+ |                IPAddr            | |
|                  |                                  | |
|                  | +-------------+ +--------------+ | |
|                  | | sockaddr_in | | sockaddr_in6 | | |
|                  | |   IPv4Addr  | |   IPv6Addr   | | |
|                  | +-------------+ +--------------+ | |
|                  |                                  | |
|                  +----------------------------------+ |
|                                                       |
+-------------------------------------------------------+
```

## Sorting

`go-sockaddr` was designed to permit sorting of heterogeneous `SockAddr`
addresses in different ways.  For example, it may be desirable to sort a
collection of IPv4 and IPv6 addresses by the size of the network.  This
allows a consumer to select the most specific IP address first (i.e. a /32
over a /120).

SortPolicy
* Type (i.e. Unix, IPv4, IPv6)
* LocalFirst (RFC1918, or IPv6 Site-Local)
* NetworkSize (/32 first, followed by /120).

To incorporate:

https://tools.ietf.org/html/rfc3849
https://www.iana.org/assignments/ipv6-unicast-address-assignments/ipv6-unicast-address-assignments.xhtml
https://www.iana.org/assignments/ipv6-address-space/ipv6-address-space.xhtml
https://tools.ietf.org/html/rfc4291#section-2.5.2

