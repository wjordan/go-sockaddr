package command

import (
	"flag"
	"fmt"

	"github.com/hashicorp/errwrap"
	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
)

type DumpCommand struct {
	Ui cli.Ui

	// attrNames is a list of attribute names to include in the output
	attrNames []string

	// flags is a list of options belonging to this command
	flags *flag.FlagSet

	// machineMode changes the output format to be machine friendly
	// (i.e. tab-separated values).
	machineMode bool

	// valueOnly changes the output format to include only values
	valueOnly bool

	// ipOnly parses the input as an IP address (either IPv4 or IPv6)
	ipOnly bool

	// v4Only parses the input exclusively as an IPv4 address
	v4Only bool

	// v6Only parses the input exclusively as an IPv6 address
	v6Only bool

	// unixOnly parses the input exclusively as a UNIX Socket
	unixOnly bool
}

// Description is the long-form command help.
func (c *DumpCommand) Description() string {
	return `Parse address(es) and dumps various output.`
}

// Help returns the full help output expected by `sockaddr -h cmd`
func (c *DumpCommand) Help() string {
	return MakeHelp(c)
}

// InitOpts is responsible for setup of this command's configuration via the
// command line.  InitOpts() does not parse the arguments (see parseOpts()).
func (c *DumpCommand) InitOpts() {
	c.flags = flag.NewFlagSet("dump", flag.ContinueOnError)
	c.flags.Usage = func() { c.Ui.Output(c.Help()) }
	c.flags.BoolVar(&c.machineMode, "H", false, "Machine readable output")
	c.flags.BoolVar(&c.valueOnly, "n", false, "Show only the value")
	c.flags.BoolVar(&c.v4Only, "4", false, "Parse the address as IPv4 only")
	c.flags.BoolVar(&c.v6Only, "6", false, "Parse the address as IPv6 only")
	c.flags.BoolVar(&c.ipOnly, "i", false, "Parse the address as IP address (either IPv4 or IPv6)")
	c.flags.BoolVar(&c.unixOnly, "u", false, "Parse the address as a UNIX Socket only")
	c.flags.Var((*MultiArg)(&c.attrNames), "o", "Name of an attribute to pass through")
}

// Run executes this command.
func (c *DumpCommand) Run(args []string) int {
	if len(args) == 0 {
		c.Ui.Error(c.Help())
		return 1
	}

	c.InitOpts()
	addrs, err := c.parseOpts(args)
	if err != nil {
		if !errwrap.Contains(err, "flag: help requested") {
			c.Ui.Error(fmt.Sprintf("error parsing args: %v", err))
		}
		return 1
	}
	for _, addr := range addrs {
		var sa sockaddr.SockAddr
		var err error
		switch {
		case c.v4Only:
			sa, err = sockaddr.NewIPv4Addr(addr)
		case c.v6Only:
			sa, err = sockaddr.NewIPv6Addr(addr)
		case c.unixOnly:
			sa, err = sockaddr.NewUnixSock(addr)
		case c.ipOnly:
			sa, err = sockaddr.NewIPAddr(addr)
		default:
			sa, err = sockaddr.NewSockAddr(addr)
		}
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Unable to parse %+q: %v", addr, err))
			return 1
		}
		c.dumpSockAddr(sa)
	}
	return 0
}

// Synopsis returns a terse description used when listing sub-commands.
func (c *DumpCommand) Synopsis() string {
	return `Parses IP addresses`
}

// Usage is the one-line usage description
func (c *DumpCommand) Usage() string {
	return `sockaddr dump [options] address [...]`
}

// VisitAllFlags forwards the visitor function to the FlagSet
func (c *DumpCommand) VisitAllFlags(fn func(*flag.Flag)) {
	c.flags.VisitAll(fn)
}

func (c *DumpCommand) dumpSockAddr(sa sockaddr.SockAddr) {
	reservedAttrs := []sockaddr.AttrName{"Attribute"}
	const maxNumAttrs = 32

	output := make([]string, 0, maxNumAttrs+len(reservedAttrs))
	allowedAttrs := make(map[sockaddr.AttrName]struct{}, len(c.attrNames)+len(reservedAttrs))
	for _, attr := range reservedAttrs {
		allowedAttrs[attr] = struct{}{}
	}
	for _, attr := range c.attrNames {
		allowedAttrs[sockaddr.AttrName(attr)] = struct{}{}
	}

	// allowedAttr returns true if the attribute is allowed to be appended
	// to the output.
	allowedAttr := func(k sockaddr.AttrName) bool {
		if len(allowedAttrs) == len(reservedAttrs) {
			return true
		}

		_, found := allowedAttrs[k]
		return found
	}

	// outFmt is a small helper function to reduce the tedium below.  outFmt
	// returns a new slice and expects the value to already be a string.
	outFmt := func(o []string, k sockaddr.AttrName, v interface{}) []string {
		if !allowedAttr(k) {
			return o
		}
		switch {
		case c.valueOnly:
			return append(o, fmt.Sprintf("%s", v))
		case !c.valueOnly && c.machineMode:
			return append(o, fmt.Sprintf("%s\t%s", k, v))
		case !c.valueOnly && !c.machineMode:
			fallthrough
		default:
			return append(o, fmt.Sprintf("%s | %s", k, v))
		}
	}

	if !c.machineMode {
		output = outFmt(output, "Attribute", "Value")
	}

	// Attributes for all SockAddr types
	for _, attr := range sockaddr.SockAddrAttrs() {
		output = outFmt(output, attr, sockaddr.SockAddrAttr(sa, attr))
	}

	// Attributes for all IP types (both IPv4 and IPv6)
	if sa.Type()&sockaddr.TypeIP != 0 {
		ip := *sockaddr.ToIPAddr(sa)
		for _, attr := range sockaddr.IPAttrs() {
			output = outFmt(output, attr, sockaddr.IPAddrAttr(ip, attr))
		}
	}

	if sa.Type() == sockaddr.TypeIPv4 {
		ipv4 := *sockaddr.ToIPv4Addr(sa)
		for _, attr := range sockaddr.IPv4Attrs() {
			output = outFmt(output, attr, sockaddr.IPv4AddrAttr(ipv4, attr))
		}
	}

	if sa.Type() == sockaddr.TypeIPv6 {
		ipv6 := *sockaddr.ToIPv6Addr(sa)
		for _, attr := range sockaddr.IPv6Attrs() {
			output = outFmt(output, attr, sockaddr.IPv6AddrAttr(ipv6, attr))
		}
	}

	if sa.Type() == sockaddr.TypeUnix {
		us := *sockaddr.ToUnixSock(sa)
		for _, attr := range sockaddr.UnixSockAttrs() {
			output = outFmt(output, attr, sockaddr.UnixSockAttr(us, attr))
		}
	}

	// Developer-focused arguments
	{
		arg1, arg2 := sa.DialPacketArgs()
		output = outFmt(output, "DialPacket", fmt.Sprintf("%+q %+q", arg1, arg2))
	}
	{
		arg1, arg2 := sa.DialStreamArgs()
		output = outFmt(output, "DialStream", fmt.Sprintf("%+q %+q", arg1, arg2))
	}
	{
		arg1, arg2 := sa.ListenPacketArgs()
		output = outFmt(output, "ListenPacket", fmt.Sprintf("%+q %+q", arg1, arg2))
	}
	{
		arg1, arg2 := sa.ListenStreamArgs()
		output = outFmt(output, "ListenStream", fmt.Sprintf("%+q %+q", arg1, arg2))
	}

	result := columnize.SimpleFormat(output)
	c.Ui.Output(result)

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

// parseOpts is responsible for parsing the options set in InitOpts().  Returns
// a list of non-parsed flags.
func (c *DumpCommand) parseOpts(args []string) ([]string, error) {
	if err := c.flags.Parse(args); err != nil {
		return nil, err
	}

	conflictingOptsCount := 0
	if c.v4Only {
		conflictingOptsCount++
	}
	if c.v6Only {
		conflictingOptsCount++
	}
	if c.unixOnly {
		conflictingOptsCount++
	}
	if c.ipOnly {
		conflictingOptsCount++
	}
	if conflictingOptsCount > 1 {
		return nil, fmt.Errorf("Conflicting options specified, only one parsing mode may be specified at a time")
	}

	return c.flags.Args(), nil
}
