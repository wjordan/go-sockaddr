package template_test

import (
	"testing"

	socktmpl "github.com/hashicorp/go-sockaddr/template"
)

func TestSockAddr_Parse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
		fail   bool
	}{
		{
			name:   "basic includeByName",
			input:  `{{interfaceAddrs | includeByName "lo0" | printf "%v"}}`,
			output: `[{[100:: 127.0.0.1/8 fe80::1/64] {1 16384 lo0  up|loopback|multicast}}]`,
		},
		{
			name:   "includeByName regexp",
			input:  `{{interfaceAddrs | includeByName "^(en|lo)0$" | excludeByName "^en0$" | printf "%v"}}`,
			output: `[{[100:: 127.0.0.1/8 fe80::1/64] {1 16384 lo0  up|loopback|multicast}}]`,
		},
		{
			name:   "excludeByName",
			input:  `{{. | includeByName "^(en|lo)0$" | excludeByName "^en0$" | printf "%v"}}`,
			output: `[{[100:: 127.0.0.1/8 fe80::1/64] {1 16384 lo0  up|loopback|multicast}}]`,
		},
		{
			name:   `"dot" pipeline, IPv4 type`,
			input:  `{{. | includeByType "IPv4" | includeByName "^lo0$"}}`,
			output: `[{[127.0.0.1/8] {1 16384 lo0  up|loopback|multicast}}]`,
		},
		{
			name:   "includeByType IPv6",
			input:  `{{. | includeByType "IPv6" | includeByName "^lo0$"}}`,
			output: `[{[100:: fe80::1/64] {1 16384 lo0  up|loopback|multicast}}]`,
		},
		{
			name:   "better regexp example for IP types",
			input:  `{{. | includeByType "^IPv[46]$" | includeByName "^lo0$"}}`,
			output: `[{[100:: 127.0.0.1/8 fe80::1/64] {1 16384 lo0  up|loopback|multicast}}]`,
		},
		{
			name:   "ifAddrs1",
			input:  `{{. | includeByType "^IPv4$" | includeByName "^lo0$" | ifAddrs }}`,
			output: `[127.0.0.1/8]`,
		},
		{
			name:   "ifAddrs2",
			input:  `{{. | includeByType "^IPv(4|6)$" | includeByName "^lo0$" | ifAddrs }}`,
			output: `[100:: 127.0.0.1/8 fe80::1/64]`,
		},
		{
			name:   "ifNames",
			input:  `{{. | includeByType "^IPv(4|6)$" | includeByName "^lo0$" | ifNames }}`,
			output: `[lo0]`,
		},
		{
			name:   `range "dot" example`,
			input:  `{{range . | includeByType "^IPv(4|6)$" | includeByName "^lo0$" | ifNames}}{{.}} {{end}}{{range . | includeByType "^IPv(4|6)$" | includeByName "^lo0$" | ifAddrs}}{{.}} {{end}}`,
			output: `lo0 100:: 127.0.0.1/8 fe80::1/64 `,
		},
		{
			name:   "excludeByType",
			input:  `{{range . | excludeByType "^IPv(4)$" | includeByName "^lo0$" | ifNames}}{{.}} {{end}}{{range . | excludeByType "^IPv(4)$" | includeByName "^lo0$" | ifAddrs}}{{.}} {{end}}`,
			output: `lo0 100:: fe80::1/64 `,
		},
		{
			name:   "with variable pipeline",
			input:  `{{with $ifSet := includeByType "^IPv(4)$" . | includeByName "^lo0$"}}{{range $ifSet | ifNames }}{{.}} {{end}}{{range $ifSet | ifAddrs}}{{.}} {{end}}{{end}}`,
			output: `lo0 127.0.0.1/8 `,
		},
		// {
		// 	input:  `{{with $ifSet := excludeByRFC 1918 . | includeByName "^lo0$" | includeByType "^IPv4$"}}{{range $ifSet | ifNames }}{{.}} {{end}}{{range $ifSet | ifAddrs}}{{.}} {{end}}{{end}}`,
		// 	output: `lo0 127.0.0.1/8 `,
		// },
		{
			// NOTE(sean@): Difficult to reliably test includeByRFC.
			// In this case, we ass-u-me that the host running the
			// test has at least one RFC1918 address on their host
			// and that its length is greater than len("[]").
			name:   "includeByRFC",
			input:  `{{. | includeByRFC 1918 | ifNames | print | len | lt 2}}`,
			output: `true`,
		},
		{
			name:   "test for non-empty array",
			input:  `{{. | includeByType "^IPv4$" | includeByRFC 1918 | ifNames | print | len | eq (len "[]")}}`,
			output: `false`,
		},
		{
			// NOTE(sean@): There are no non-IPv4 RFC1918 addresses.
			name:   "non-IPv4 RFC1918",
			input:  `{{. | excludeByType "^IPv4$" | includeByRFC 1918 | ifAddrs | print | len | eq (len "[]")}}`,
			output: `true`,
		},
		{
			// NOTE(sean@): There are no RFC6598 addresses on most testing hosts so this should be empty.
			name:   "rfc6598",
			input:  `{{. | includeByType "^IPv4$" | includeByRFC 6598 | ifAddrs | print | len | eq (len "[]")}}`,
			output: `true`,
		},
		{
			name:   "invalid RFC",
			input:  `{{. | includeByType "^IPv4$" | includeByRFC 99999999999 | ifAddrs | print | len | eq (len "[]")}}`,
			output: `true`,
			fail:   true,
		},
		{
			name:   "sortByAddr",
			input:  `{{with $ifSet := includeByName "lo0" . }}{{ range includeByType "IPv4" $ifSet | ifAddrs | sortByAddr}}{{ . }} {{end}}{{ range includeByType "IPv6" $ifSet | ifAddrs | sortByAddr}}{{ . }} {{end}}{{end}}`,
			output: `127.0.0.1/8 100:: fe80::1/64 `,
		},
		{
			name:   "sortByAddr with reverse",
			input:  `{{with $ifSet := includeByName "lo0" . }}{{ range includeByType "IPv6" $ifSet | ifAddrs | sortByAddr | reverseAddrs}}{{ . }} {{end}}{{end}}`,
			output: `fe80::1/64 100:: `,
		},
		{
			name:   "sortByPort with reverse",
			input:  `{{with $ifSet := includeByName "lo0" . }}{{ range includeByType "IPv6" $ifSet | ifAddrs | sortByAddr | reverseAddrs}}{{ . }} {{end}}{{end}}`,
			output: `fe80::1/64 100:: `,
		},
		{
			name:   "lo0 limit 1",
			input:  `{{. | includeByName "lo0" | includeByType "IPv6" | ifAddrs | sortByAddr | limitAddrs 1}}`,
			output: `[100::]`,
		},
		{
			name:   "joinAddrs",
			input:  `{{. | includeByName "lo0" | includeByType "IPv6" | ifAddrs | sortByAddr | joinAddrs " "}}`,
			output: `100:: fe80::1/64`,
		},
		// {
		// 	// NOTE(sean@): This is the HashiCorp default in 2016.
		// 	name:   "HashiCorpDefault2016",
		// 	input:  `{{range . | includeByType "^IP(v[46])?$" | includeByRFC 1918 | sortByAddr }}{{ . }} {{end}}{{range . | includeByType "^IP(v[46])?$" | includeByRFC 6598 | sortByAddr }}{{ . }} {{end}}`,
		// 	output: `true`,
		// },
	}

	for _, test := range tests {
		test := test // capture range variable
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			out, err := socktmpl.Parse(test.input)
			if err != nil && !test.fail {
				t.Fatalf("%q: bad: %v", test.name, err)
			}

			if out != test.output && !test.fail {
				t.Fatalf("%q: Expected %+q, received %+q\n%+q", test.name, test.output, out, test.input)
			}
		})
	}
}
