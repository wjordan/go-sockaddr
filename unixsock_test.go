package sockaddr_test

import (
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

func TestUnixSock_impl_SockAddr(t *testing.T) {
	tests := []struct {
		name             string
		input            sockaddr.UnixSock
		dialPacketArgs   []string
		dialStreamArgs   []string
		listenPacketArgs []string
		listenStreamArgs []string
	}{
		{
			name:             "simple",
			input:            sockaddr.MustUnixSock("/tmp/foo"),
			dialPacketArgs:   []string{"unixgram", "/tmp/foo"},
			dialStreamArgs:   []string{"unixgram", "/tmp/foo"},
			listenPacketArgs: []string{"unixgram", "/tmp/foo"},
			listenStreamArgs: []string{"unixgram", "/tmp/foo"},
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		arg1, arg2 := test.input.DialPacketArgs()
		if arg1 != test.dialPacketArgs[0] && arg2 != test.dialPacketArgs[1] {
			t.Fatalf("%s: %q %q", test.name, arg1, arg2)
		}

		arg1, arg2 = test.input.DialStreamArgs()
		if arg1 != test.dialStreamArgs[0] && arg2 != test.dialStreamArgs[1] {
			t.Fatalf("%s: %q %q", test.name, arg1, arg2)
		}

		arg1, arg2 = test.input.ListenPacketArgs()
		if arg1 != test.listenPacketArgs[0] && arg2 != test.listenPacketArgs[1] {
			t.Fatalf("%s: %q %q", test.name, arg1, arg2)
		}

		arg1, arg2 = test.input.ListenStreamArgs()
		if arg1 != test.listenStreamArgs[0] && arg2 != test.listenStreamArgs[1] {
			t.Fatalf("%s: %q %q", test.name, arg1, arg2)
		}
	}
}

func TestUnixSockAttrs(t *testing.T) {
	const expectedNumAttrs = 1
	usa := sockaddr.UnixSockAttrs()
	if len(usa) != expectedNumAttrs {
		t.Fatalf("wrong number of UnixSockAttrs: %d vs %d", len(usa), expectedNumAttrs)
	}
}
