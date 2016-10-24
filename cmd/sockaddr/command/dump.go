package command

import (
	"fmt"
	"math/big"
	"strings"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
)

type DumpCommand struct {
	Ui cli.Ui
}

func (c *DumpCommand) Help() string {
	helpText := `
Usage: sockaddr dump [address ...]

  Parse address(es) and dumps various output.
`
	return strings.TrimSpace(helpText)
}

func (c *DumpCommand) Run(args []string) int {
	if len(args) == 0 {
		c.Ui.Error(fmt.Sprintf("%s", c.Help()))
		return 1
	}

	for _, arg := range args {
		sa, err := sockaddr.NewSockAddr(arg)
		if err != nil {
			return 1
		}
		dumpSockAddr(c.Ui, sa)
	}
	return 0
}

func (c *DumpCommand) Synopsis() string {
	return "Parses IP addresses"
}

func dumpSockAddr(ui cli.Ui, sa sockaddr.SockAddr) {
	const numAttrs = 17
	output := make([]string, 0, numAttrs)

	output = append(output, "Attribute | Value")
	output = append(output, fmt.Sprintf("type | %s", sa.Type()))
	output = append(output, fmt.Sprintf("string | %s", sa.String()))

	// Attributes for all IP types (both IPv4 and IPv6)
	if sa.Type()&sockaddr.TypeIP != 0 {
		ip := *sockaddr.ToIPAddr(sa)
		output = append(output, fmt.Sprintf("host | %s", ip.Host()))
		output = append(output, fmt.Sprintf("port | %d", ip.IPPort()))
		output = append(output, fmt.Sprintf("network address | %s", ip.NetIP()))
		output = append(output, fmt.Sprintf("network mask | %s", ip.NetIPMask()))
		output = append(output, fmt.Sprintf("network | %s", ip.Network()))
		output = append(output, fmt.Sprintf("mask bits | %d", ip.Maskbits()))
		output = append(output, fmt.Sprintf("binary | %s", ip.AddressBinString()))
		output = append(output, fmt.Sprintf("hex | %s", ip.AddressHexString()))
		output = append(output, fmt.Sprintf("first usable | %s", ip.FirstUsable()))
		output = append(output, fmt.Sprintf("last usable | %s", ip.LastUsable()))

		{
			octets := ip.Octets()
			octetStrs := make([]string, 0, len(octets))
			for _, octet := range octets {
				octetStrs = append(octetStrs, fmt.Sprintf("%d", octet))
			}
			output = append(output, fmt.Sprintf("octets | %s", strings.Join(octetStrs, " ")))
		}
	}

	if sa.Type() == sockaddr.TypeIPv4 {
		ipv4 := *sockaddr.ToIPv4Addr(sa)
		output = append(output, fmt.Sprintf("broadcast | %s", ipv4.Broadcast()))
		output = append(output, fmt.Sprintf("uint32 | %d", uint32(ipv4.Address)))
	}

	if sa.Type() == sockaddr.TypeIPv6 {
		ipv6 := *sockaddr.ToIPv6Addr(sa)
		{
			b := big.Int(*ipv6.Address)
			output = append(output, fmt.Sprintf("uint128 | %s", b.Text(10)))
		}
	}

	if sa.Type() == sockaddr.TypeUnix {
		us := *sockaddr.ToUnixSock(sa)
		output = append(output, fmt.Sprintf("path | %s", us.Path()))
	}

	// Developer-focused arguments
	{
		arg1, arg2 := sa.DialPacketArgs()
		output = append(output, fmt.Sprintf("dial packet args | %+q %+q", arg1, arg2))
	}
	{
		arg1, arg2 := sa.DialStreamArgs()
		output = append(output, fmt.Sprintf("dial stream args | %+q %+q", arg1, arg2))
	}
	{
		arg1, arg2 := sa.ListenPacketArgs()
		output = append(output, fmt.Sprintf("listen packet args | %+q %+q", arg1, arg2))
	}
	{
		arg1, arg2 := sa.ListenStreamArgs()
		output = append(output, fmt.Sprintf("listen stream args | %+q %+q", arg1, arg2))
	}

	result := columnize.SimpleFormat(output)
	ui.Output(result)

	// fmt.Printf("SockAddr.Address.IsGlobalUnicast(): %v\n", na.Address.IsGlobalUnicast())
	// fmt.Printf("SockAddr.Address.IsInterfaceLocalMulticast(): %v\n", na.Address.IsInterfaceLocalMulticast())
	// fmt.Printf("SockAddr.Address.IsLinkLocalMulticast(): %v\n", na.Address.IsLinkLocalMulticast())
	// fmt.Printf("SockAddr.Address.IsLinkLocalUnicast(): %v\n", na.Address.IsLinkLocalUnicast())
	// fmt.Printf("SockAddr.Address.IsLoopback(): %v\n", na.Address.IsLoopback())
	// fmt.Printf("SockAddr.Address.IsMulticast(): %v\n", na.Address.IsMulticast())
	// fmt.Printf("SockAddr.Address.IsUnspecified(): %v\n", na.Address.IsUnspecified())
	// fmt.Printf("SockAddr.Address.To16(): %s\n", na.Address.To16())
	// fmt.Printf("SockAddr.Address.To4(): %s\n", na.Address.To4())
	// ipuint, ok := na.ToUint32()
	// if !ok {
	// 	panic("Unable to uint32")
	// }

	// fmt.Printf("SockAddr.ToUint32(): %d\n", ipuint)
}
