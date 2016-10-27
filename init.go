package sockaddr

var rfcNetMap map[uint][]SockAddr

func init() {
	rfcNetMap = KnownRFCs()
}
