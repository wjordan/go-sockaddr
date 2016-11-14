# go-sockaddr

## `sockaddr` Library
Socket address convenience functions for Go.  `go-sockaddr` is a convenience
library that makes doing the right thing with IP addresses easy.  `go-sockaddr`
is loosely modeled after the UNIX `sockaddr_t` and creates a union of the family
of `sockaddr_t` types (see below for an ascii diagram).  Library documentation
is available
at
[https://godoc.org/github.com/hashicorp/go-sockaddr](https://godoc.org/github.com/hashicorp/go-sockaddr).
The primary intent of the library was to make it possible to define heuristics
for selecting IP addresses at process initialization time.  See the docs, tests,
and CLI utility for details and hints as to how to use this library.

## `sockaddr` CLI

Given the possible complexity of the `sockaddr` library, there is a CLI utility
that accompanies the library, also
called
[`sockaddr`](https://github.com/hashicorp/go-sockaddr/tree/master/cmd/sockaddr).
The
[`sockaddr`](https://github.com/hashicorp/go-sockaddr/tree/master/cmd/sockaddr)
utility exposes nearly all of the functionailty of the library and can be used
either as an administrative tool or testing tool.  To install
the
[`sockaddr`](https://github.com/hashicorp/go-sockaddr/tree/master/cmd/sockaddr),
run:

```text
$ go install github.com/hashicorp/go-sockaddr/cmd/sockaddr
```

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
