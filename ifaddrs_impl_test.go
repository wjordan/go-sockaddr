package sockaddr

import "testing"

func Test_parseBSDDefaultIfName(t *testing.T) {
	testCases := []struct {
		name     string
		routeOut string
		want     string
	}{
		{
			name: "macOS Sierra 10.12 - Common",
			routeOut: `   route to: default
destination: default
       mask: default
    gateway: 10.23.9.1
  interface: en0
      flags: <UP,GATEWAY,DONE,STATIC,PRCLONING>
 recvpipe  sendpipe  ssthresh  rtt,msec    rttvar  hopcount      mtu     expire
       0         0         0         0         0         0      1500         0 
`,
			want: "en0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseBSDDefaultIfName([]byte(tc.routeOut))
			if err != nil {
				t.Fatalf("unable to parse default interface from route output: %v", err)
			}

			if got != tc.want {
				t.Errorf("got %s; want %s", got, tc.want)
			}
		})
	}
}
