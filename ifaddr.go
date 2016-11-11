package sockaddr

import "net"

// IfAddr is a union of a SockAddr and a net.Interface.
type IfAddr struct {
	SockAddr
	net.Interface
}
