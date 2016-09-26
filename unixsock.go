package sockaddr

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

// Type is used as a type switch and returns TypeUnix
func (UnixSock) Type() SockAddrType {
	return TypeUnix
}
