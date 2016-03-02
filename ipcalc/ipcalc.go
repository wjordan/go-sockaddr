package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-netaddr"
	"github.com/mitchellh/cli"
)

type dumpCommand struct{}

func (c *dumpCommand) Help() string {
	helpText := `
Usage: ipcalc dump [address ...]

  Parse address(es) and dumps various output.
`
	return strings.TrimSpace(helpText)
}

func (c *dumpCommand) Run(args []string) int {
	for _, arg := range args {
		na, err := netaddr.New(arg)
		if err != nil {
			return 1
		}
		dumpNetAddr(na)
	}
	return 0
}

func (c *dumpCommand) Synopsis() string {
	return "Parses IP addresses"
}

func dumpNetAddr(na *netaddr.NetAddr) {
	fmt.Printf("NetAddr.Address: %s\n", na.Address)
	fmt.Printf("NetAddr.Network: %s\n", na.Network)

	fmt.Printf("NetAddr.Network.Network(): %s\n", na.Network.Network())
	fmt.Printf("NetAddr.Network.String(): %s\n", na.Network.String())

	fmt.Printf("NetAddr.Address.IsGlobalUnicast(): %v\n", na.Address.IsGlobalUnicast())
	fmt.Printf("NetAddr.Address.IsInterfaceLocalMulticast(): %v\n", na.Address.IsInterfaceLocalMulticast())
	fmt.Printf("NetAddr.Address.IsLinkLocalMulticast(): %v\n", na.Address.IsLinkLocalMulticast())
	fmt.Printf("NetAddr.Address.IsLinkLocalUnicast(): %v\n", na.Address.IsLinkLocalUnicast())
	fmt.Printf("NetAddr.Address.IsLoopback(): %v\n", na.Address.IsLoopback())
	fmt.Printf("NetAddr.Address.IsMulticast(): %v\n", na.Address.IsMulticast())
	fmt.Printf("NetAddr.Address.IsUnspecified(): %v\n", na.Address.IsUnspecified())
	fmt.Printf("NetAddr.Address.String(): %v\n", na.Address.String())
	fmt.Printf("NetAddr.Address.To16(): %s\n", na.Address.To16())
	fmt.Printf("NetAddr.Address.To4(): %s\n", na.Address.To4())
	ipuint, ok := na.ToUint32()
	if !ok {
		panic("Unable to uint32")
	}

	fmt.Printf("NetAddr.NetworkAddress(): %s\n", na.NetworkAddress())
	fmt.Printf("NetAddr.BroadcastAddress(): %s\n", na.BroadcastAddress())
	fmt.Printf("NetAddr.Maskbits(): %d\n", na.Maskbits())
	fmt.Printf("NetAddr.ToBinString(): %s\n", na.ToBinString())
	fmt.Printf("NetAddr.ToHexString(): %s\n", na.ToHexString())
	fmt.Printf("NetAddr.ToUint32(): %d\n", ipuint)
}

func main() {
	c := cli.NewCLI("ipcalc", "0.1.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"dump": func() (cli.Command, error) {
			return &dumpCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
