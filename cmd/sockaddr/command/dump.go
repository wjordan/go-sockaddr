package command

import (
	"strings"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/mitchellh/cli"
)

type DumpCommand struct {
	Ui cli.Ui
}

func (c *DumpCommand) Help() string {
	helpText := `
Usage: ipcalc dump [address ...]

  Parse address(es) and dumps various output.
`
	return strings.TrimSpace(helpText)
}

func (c *DumpCommand) Run(args []string) int {
	for _, arg := range args {
		sa, err := sockaddr.NewSockAddr(arg)
		if err != nil {
			return 1
		}
		dumpSockAddr(&sa)
	}
	return 0
}

func (c *DumpCommand) Synopsis() string {
	return "Parses IP addresses"
}

func dumpSockAddr(na *sockaddr.SockAddr) {
	// fmt.Printf("SockAddr.Address: %s\n", na.Address)
	// fmt.Printf("SockAddr.Network: %s\n", na.Network)

	// fmt.Printf("SockAddr.Network.Network(): %s\n", na.Network.Network())
	// fmt.Printf("SockAddr.Network.String(): %s\n", na.Network.String())

	// fmt.Printf("SockAddr.Address.IsGlobalUnicast(): %v\n", na.Address.IsGlobalUnicast())
	// fmt.Printf("SockAddr.Address.IsInterfaceLocalMulticast(): %v\n", na.Address.IsInterfaceLocalMulticast())
	// fmt.Printf("SockAddr.Address.IsLinkLocalMulticast(): %v\n", na.Address.IsLinkLocalMulticast())
	// fmt.Printf("SockAddr.Address.IsLinkLocalUnicast(): %v\n", na.Address.IsLinkLocalUnicast())
	// fmt.Printf("SockAddr.Address.IsLoopback(): %v\n", na.Address.IsLoopback())
	// fmt.Printf("SockAddr.Address.IsMulticast(): %v\n", na.Address.IsMulticast())
	// fmt.Printf("SockAddr.Address.IsUnspecified(): %v\n", na.Address.IsUnspecified())
	// fmt.Printf("SockAddr.Address.String(): %v\n", na.Address.String())
	// fmt.Printf("SockAddr.Address.To16(): %s\n", na.Address.To16())
	// fmt.Printf("SockAddr.Address.To4(): %s\n", na.Address.To4())
	// ipuint, ok := na.ToUint32()
	// if !ok {
	// 	panic("Unable to uint32")
	// }

	// fmt.Printf("SockAddr.NetworkAddress(): %s\n", na.NetworkAddress())
	// fmt.Printf("SockAddr.BroadcastAddress(): %s\n", na.BroadcastAddress())
	// fmt.Printf("SockAddr.Maskbits(): %d\n", na.Maskbits())
	// fmt.Printf("SockAddr.ToBinString(): %s\n", na.ToBinString())
	// fmt.Printf("SockAddr.ToHexString(): %s\n", na.ToHexString())
	// fmt.Printf("SockAddr.ToUint32(): %d\n", ipuint)
}
