package sockaddr

var rfcNetMap map[uint]SockAddrs

func init() {
	rfcNetMap = KnownRFCs()
}
