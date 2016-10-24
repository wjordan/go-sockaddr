package sockaddr

import "fmt"

type UnixSock struct {
	SockAddr
	path string
}
type UnixSocks []*UnixSock

// NewUnixSock creates an UnixSock from a string path.  String can be in the
// form of either URI-based string (e.g. `file:///etc/passwd`), an absolute
// path (e.g. `/etc/passwd`), or a relative path (e.g. `./foo`).
func NewUnixSock(s string) (ret UnixSock, err error) {
	ret.path = s
	return ret, nil
}

// DialPacketArgs returns the arguments required to be passed to net.DialUnix()
// with the `unixgram` network type.
func (us UnixSock) DialPacketArgs() (network, dialArgs string) {
	return "unixgram", us.path
}

// DialStreamArgs returns the arguments required to be passed to net.DialUnix()
// with the `unix` network type.
func (us UnixSock) DialStreamArgs() (network, dialArgs string) {
	return "unix", us.path
}

// ListenPacketArgs returns the arguments required to be passed to
// net.ListenUnixgram() with the `unixgram` network type.
func (us UnixSock) ListenPacketArgs() (network, dialArgs string) {
	return "unixgram", us.path
}

// ListenStreamArgs returns the arguments required to be passed to
// net.ListenUnix() with the `unix` network type.
func (us UnixSock) ListenStreamArgs() (network, dialArgs string) {
	return "unix", us.path
}

// Path returns the given path of the UnixSock
func (us UnixSock) Path() string {
	return us.path
}

// String returns the path of the UnixSock
func (us UnixSock) String() string {
	return fmt.Sprintf("%+q", us.path)
}

// Type is used as a type switch and returns TypeUnix
func (UnixSock) Type() SockAddrType {
	return TypeUnix
}
