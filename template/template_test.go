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
			name:   "basic includeByIfName",
			input:  `{{GetIfSockAddrs | includeByIfName "lo0" | printf "%v"}}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast} 100:: {1 16384 lo0  up|loopback|multicast} fe80::1/64 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			// NOTE(sean@): This test, as its written now, will only
			// pass on macOS.
			name:   "GetDefaultInterface",
			input:  `{{GetDefaultInterface | includeByType "IPv4" | limit 1 | join "name" " " }}`,
			output: `en0`,
		},
		{
			name:   "includeByIfName regexp",
			input:  `{{GetIfSockAddrs | includeByIfName "^(en|lo)0$" | excludeByIfName "^en0$" | sortByType | sortByAddr | join "address" " " }}`,
			output: `127.0.0.1/8 100:: fe80::1/64`,
		},
		{
			name:   "excludeByIfName",
			input:  `{{. | includeByIfName "^(en|lo)0$" | excludeByIfName "^en0$" | sortByType | sortByAddr | join "address" " " }}`,
			output: `127.0.0.1/8 100:: fe80::1/64`,
		},
		{
			name:   `"dot" pipeline, IPv4 type`,
			input:  `{{. | includeByType "IPv4" | includeByIfName "^lo0$" | sortByType | sortByAddr }}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   "includeByType IPv6",
			input:  `{{. | includeByType "IPv6" | includeByIfName "^lo0$" | sortByAddr }}`,
			output: `[100:: {1 16384 lo0  up|loopback|multicast} fe80::1/64 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   "better regexp example for IP types",
			input:  `{{. | includeByType "^IPv[46]$" | includeByIfName "^lo0$" | sortByType | sortByAddr }}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast} 100:: {1 16384 lo0  up|loopback|multicast} fe80::1/64 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   "ifAddrs1",
			input:  `{{. | includeByType "^IPv4$" | includeByIfName "^lo0$"}}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   "ifAddrs2",
			input:  `{{. | includeByType "^IPv(4|6)$" | includeByIfName "^lo0$" | sortByType | sortByAddr }}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast} 100:: {1 16384 lo0  up|loopback|multicast} fe80::1/64 {1 16384 lo0  up|loopback|multicast}]`,
		},
		// {
		// 	name:   "ifNames",
		// 	input:  `{{. | includeByType "^IPv(4|6)$" | includeByIfName "^lo0$" | ifNames }}`,
		// 	output: `[lo0]`,
		// },
		{
			name:   `range "dot" example`,
			input:  `{{range . | includeByType "^IPv(4|6)$" | includeByIfName "^lo0$"}}{{.Name}} {{.SockAddr}} {{end}}`,
			output: `lo0 127.0.0.1/8 lo0 100:: lo0 fe80::1/64 `,
		},
		{
			name:   "excludeByType",
			input:  `{{. | excludeByType "^IPv(4)$" | includeByIfName "^lo0$" | sortByAddr | uniqueBy "name" | join "name" " "}} {{range . | excludeByType "^IPv(4)$" | includeByIfName "^lo0$"}}{{.SockAddr}} {{end}}`,
			output: `lo0 100:: fe80::1/64 `,
		},
		{
			name:   "with variable pipeline",
			input:  `{{with $ifSet := includeByType "^IPv(4)$" . | includeByIfName "^lo0$"}}{{range $ifSet }}{{.Name}} {{end}}{{range $ifSet}}{{.SockAddr}} {{end}}{{end}}`,
			output: `lo0 127.0.0.1/8 `,
		},
		// {
		// 	input:  `{{with $ifSet := excludeByRFC 1918 . | includeByIfName "^lo0$" | includeByType "^IPv4$"}}{{range $ifSet}}{{.}} {{end}}{{range $ifSet}}{{.}} {{end}}{{end}}`,
		// 	output: `lo0 127.0.0.1/8 `,
		// },
		{
			// NOTE(sean@): Difficult to reliably test includeByRFC.
			// In this case, we ass-u-me that the host running the
			// test has at least one RFC1918 address on their host
			// and that its length is greater than len("[]").
			name:   "includeByRFC",
			input:  `{{(. | includeByRFC 1918 | limit 1 | join "name" " ")}}`,
			output: `en0`,
		},
		{
			name:   "test for non-empty array",
			input:  `{{. | includeByType "^IPv4$" | includeByRFC 1918 | print | len | eq (len "[]")}}`,
			output: `false`,
		},
		{
			// NOTE(sean@): There are no non-IPv4 RFC1918 addresses.
			name:   "non-IPv4 RFC1918",
			input:  `{{. | excludeByType "^IPv4$" | includeByRFC 1918 | print | len | eq (len "[]")}}`,
			output: `true`,
		},
		{
			// NOTE(sean@): There are no RFC6598 addresses on most testing hosts so this should be empty.
			name:   "rfc6598",
			input:  `{{. | includeByType "^IPv4$" | includeByRFC 6598 | print | len | eq (len "[]")}}`,
			output: `true`,
		},
		{
			name:   "invalid RFC",
			input:  `{{. | includeByType "^IPv4$" | includeByRFC 99999999999 | print | len | eq (len "[]")}}`,
			output: `true`,
			fail:   true,
		},
		{
			name:   "sortByAddr",
			input:  `{{with $ifSet := includeByIfName "lo0" . }}{{ range includeByType "IPv4" $ifSet | sortByAddr}}{{ .SockAddr }} {{end}}{{ range includeByType "IPv6" $ifSet | sortByAddr}}{{ .SockAddr }} {{end}}{{end}}`,
			output: `127.0.0.1/8 100:: fe80::1/64 `,
		},
		{
			name:   "sortByAddr with reverse",
			input:  `{{with $ifSet := includeByIfName "lo0" . }}{{ range includeByType "IPv6" $ifSet | sortByAddr | reverse}}{{ .SockAddr }} {{end}}{{end}}`,
			output: `fe80::1/64 100:: `,
		},
		{
			name:   "sortByPort with reverse",
			input:  `{{with $ifSet := includeByIfName "lo0" . }}{{ range includeByType "IPv6" $ifSet | sortByAddr | reverse}}{{ .SockAddr }} {{end}}{{end}}`,
			output: `fe80::1/64 100:: `,
		},
		{
			name:   "lo0 limit 1",
			input:  `{{. | includeByIfName "lo0" | includeByType "IPv6" | sortByAddr | limit 1 | len}}`,
			output: `1`,
		},
		{
			name:   "join address",
			input:  `{{. | includeByIfName "lo0" | includeByType "IPv6" | sortByAddr | join "address" " " }}`,
			output: `100:: fe80::1/64`,
		},
		{
			name:   "join name",
			input:  `{{. | includeByIfName "lo0" | includeByType "IPv6" | sortByAddr | join "name" " " }}`,
			output: `lo0 lo0`,
		},
		// {
		// 	name:   "lo0 flags up and limit 1",
		// 	input:  `{{. | includeByIfName "lo0" | includeByFlag "Up" | limit 1}}`,
		// 	output: `[100::]`,
		// },
		{
			// NOTE(sean@): This is the HashiCorp default in 2016.
			// Indented for effect.  Using "true" as the output
			// instead of printing the correct $rfc*Addrs values.
			name: "HashiCorpDefault2016",
			input: `
{{- with $addrs := GetIfSockAddrs | includeByType "^IP(v[46])?$" -}}
  {{- $rfc1918Addrs := $addrs | includeByRFC 1918 | sortByAddr | limit 1 | join "address" " " -}}
  {{- $rfc6598Addrs := $addrs | includeByRFC 6598 | sortByAddr | limit 1 | join "address" " " -}}

  {{- if ($rfc1918Addrs | len) gt 0 -}}
    {{- print "true" -}}{{/* print $rfc1918Addrs*/ -}}
  {{- else if ($rfc6598Addrs | len) gt 0 -}}
    {{- print "true" -}}{{/* print $rfc6598Addrs*/ -}}
  {{- end -}}
{{- end -}}`,
			output: `true`,
		},
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
