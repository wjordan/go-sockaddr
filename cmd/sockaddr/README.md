# `sockaddr(1)`

`sockaddr` is a CLI utility that wraps and exposes `go-sockaddr` functionality
from the command line.

```text
% sockaddr -h
usage: sockaddr [--version] [--help] <command> [<args>]

Available commands are:
    dump       Parses IP addresses
    eval       Evaluates a sockaddr template
    rfc        Test to see if an IP is part of a known RFC
    version    Prints the sockaddr version
```

## `sockaddr dump`

```text
Usage: sockaddr dump [options] address [...]

  Parse address(es) and dumps various output.

Options:

  -H  Machine readable output
  -o  Name of an attribute to pass through
```

### `sockaddr dump` example output

By default it prints out all available information unless the `-o` flag is
specified.

```text
% sockaddr dump 127.0.0.2/8
Attribute     Value
type          IPv4
string        127.0.0.2/8
host          127.0.0.2
address       127.0.0.2
port          0
netmask       ff000000
network       127.0.0.0/8
mask_bits     8
binary        01111111000000000000000000000010
hex           7f000002
first_usable  127.0.0.1
last_usable   127.255.255.254
octets        127 0 0 2
broadcast     127.255.255.255
uint32        2130706434
DialPacket    "udp4" ""
DialStream    "tcp4" ""
ListenPacket  "udp4" ""
ListenStream  "tcp4" ""
$ sockaddr dump -H -o host,address,port -o mask_bits 127.0.0.3:8600
host	127.0.0.3:8600
address	127.0.0.3
port	8600
mask_bits	32
$ sockaddr dump -o type,address,hex,network '[2001:db8::3/32]'
Attribute  Value
type       IPv6
address    2001:db8::3
network    2001:db8::/32
hex        20010db8000000000000000000000003
$ sockaddr dump /tmp/example.sock
Attribute     Value
type          UNIX
string        "/tmp/example.sock"
path          /tmp/example.sock
DialPacket    "unixgram" "/tmp/example.sock"
DialStream    "unix" "/tmp/example.sock"
ListenPacket  "unixgram" "/tmp/example.sock"
ListenStream  "unix" "/tmp/example.sock"
```

## `sockaddr eval`

The `sockaddr` library has the potential to be very complex, which is why the
`sockaddr` command supports an `eval` subcommand in order to test configurations
from the command line.  If the argument passed to `eval` is a dash (`-`), then
`sockaddr eval` will read from stdin.

```
Usage: sockaddr eval [options] [template ...]

  Parse the sockaddr template and evaluates the output.

Options:

  -d  Debug output
  -n  Suppress newlines between args
```

Here are a few impractical examples to get you started:

```text
$ sockaddr eval '{{. | includeByIfName "lo0" | includeByType "IPv6" | ifAddrs | sortByAddr | join "address" " "}}'
[0] in: "{{. | includeByIfName \"lo0\" | includeByType \"IPv6\" | ifAddrs | sortByAddr | join "address" \" \"}}"
[0] out: "100:: fe80::1/64"
$ sockaddr eval '{{. | includeByRFC 1918 | ifNames | print | len | lt 2}}'
[0] in: "{{. | includeByRFC 1918 | ifNames | print | len | lt 2}}"
[0] out: "true"
$ sockaddr eval '{{with $ifSet := includeByIfName "lo0" . }}{{ range includeByType "IPv6" $ifSet | ifAddrs | sortByAddr | reverse}}{{ . }} {{end}}{{end}}'
[0] in: "{{with $ifSet := includeByIfName \"lo0\" . }}{{ range includeByType \"IPv6\" $ifSet | ifAddrs | sortByAddr | reverse}}{{ . }} {{end}}{{end}}"
[0] out: "fe80::1/64 100:: "
$ sockaddr eval '{{. | includeByIfName "lo0" | includeByType "IPv6" | ifAddrs | sortByAddr | join "address" " "}}'
[0] in: "{{. | includeByIfName \"lo0\" | includeByType \"IPv6\" | ifAddrs | sortByAddr | join "address" \" \"}}"
[0] out: "100:: fe80::1/64"
$ cat <<'EOF' | sockaddr eval -
{{. | includeByIfName "lo0" | includeByType "IPv6" | ifAddrs | sortByAddr | join "address" " "}}
EOF
[0] in: "{{. | includeByIfName \"lo0\" | includeByType \"IPv6\" | ifAddrs | sortByAddr | join "address" \" \"}}"
[0] out: "100:: fe80::1/64"
```

## `sockaddr rfc`

> Tests a given IP address to see if it is part of a known RFC.  If the IP
> address belongs to a known RFC, return exit code 0 and print the status.  If
> the IP does not belong to an RFC, return 1.  If the RFC is not known, return
> 2.

```text
$ sockaddr rfc 1918 192.168.1.10
192.168.1.10 is part of RFC 1918
$ sockaddr rfc 6890 '[::1]'
100:: is part of RFC 6890
$ sockaddr rfc list
1918
4193
5735
6598
6890
```

## `sockaddr version`

The lowly version stub.

```text
$ sockaddr version
sockaddr 0.1.0-dev
```
