package sockaddr

import "net"

// IfAddr is a combined SockAddr and Interface
type IfAddr struct {
	SockAddr
	net.Interface
}
