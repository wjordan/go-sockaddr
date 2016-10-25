# `sockaddr(1)`

`sockaddr` is a CLI utility that wraps and exposes `go-sockaddr` functionality
from the command line.

```text
% sockaddr -h
usage: sockaddr [--version] [--help] <command> [<args>]

Available commands are:
    dump       Parses IP addresses
    eval       Evaluates a sockaddr template
    version    Prints the sockaddr version
```

## `sockaddr dump`

`sockaddr dump` prints out various attributes of its argument.  By default it
prints out all available information:

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
sockaddr dump -H -o host,address,port -o mask_bits 127.0.0.3:8600
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
from the command line.  Here are a few impractical examples to get you started:

```text
$ sockaddr eval '{{. | includeByIfName "lo0" | includeByType "IPv6" | ifAddrs | sortByAddr | joinAddrs " "}}'
[0] in: "{{. | includeByIfName \"lo0\" | includeByType \"IPv6\" | ifAddrs | sortByAddr | joinAddrs \" \"}}"
[0] out: "100:: fe80::1/64"
$ sockaddr eval '{{. | includeByRFC 1918 | ifNames | print | len | lt 2}}'
[0] in: "{{. | includeByRFC 1918 | ifNames | print | len | lt 2}}"
[0] out: "true"
$ sockaddr eval '{{with $ifSet := includeByIfName "lo0" . }}{{ range includeByType "IPv6" $ifSet | ifAddrs | sortByAddr | reverseAddrs}}{{ . }} {{end}}{{end}}'
[0] in: "{{with $ifSet := includeByIfName \"lo0\" . }}{{ range includeByType \"IPv6\" $ifSet | ifAddrs | sortByAddr | reverseAddrs}}{{ . }} {{end}}{{end}}"
[0] out: "fe80::1/64 100:: "
$ sockaddr eval '{{. | includeByIfName "lo0" | includeByType "IPv6" | ifAddrs | sortByAddr | joinAddrs " "}}'
[0] in: "{{. | includeByIfName \"lo0\" | includeByType \"IPv6\" | ifAddrs | sortByAddr | joinAddrs \" \"}}"
[0] out: "100:: fe80::1/64"
```

## `sockaddr version`

The lowly version stub.

```text
$ sockaddr version
sockaddr 0.1.0-dev
```
