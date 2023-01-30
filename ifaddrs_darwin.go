package sockaddr

import (
	"errors"
)

// GetAllInterfaces iterates over all available network interfaces and finds all
// available IP addresses on each interface and converts them to
// sockaddr.IPAddrs, and returning the result as an array of IfAddr.
func GetAllInterfaces() (IfAddrs, error) {
	return nil, errors.New("not implemented for darwin")
}
