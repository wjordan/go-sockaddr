package template

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/hashicorp/errwrap"
	sockaddr "github.com/hashicorp/go-sockaddr"
)

var (
	// SourceFuncs is a map of all top-level functions that generate
	// sockaddr data types.
	SourceFuncs template.FuncMap

	// SortFuncs is a map of all functions used in sorting
	SortFuncs template.FuncMap

	// FilterFuncs is a map of all functions used in sorting
	FilterFuncs template.FuncMap

	// HelperFuncs is a map of all functions used in sorting
	HelperFuncs template.FuncMap
)

func init() {
	SourceFuncs = template.FuncMap{
		// Generates a set of IfAddr inputs for the rest of the template
		// pipeline.  `GetAllInterfaces` is the default input and
		// original "dot" in the pipeline.
		//
		"GetAllInterfaces": sockaddr.GetAllInterfaces,

		// Return an IfAddr that is attached to the default route.
		"GetDefaultInterfaces": sockaddr.GetDefaultInterfaces,

		// Return a single IfAddr that represents a Private RFC 6890 IP
		// address that is attached to the default route.
		"GetPrivateInterfaces": sockaddr.GetPrivateInterfaces,

		// Return a single IfAddr that represents a Public RFC 6890 IP
		// address that is attached to the default route.
		"GetPublicInterfaces": sockaddr.GetPublicInterfaces,
	}

	SortFuncs = template.FuncMap{
		"sort": sockaddr.SortIfBy,
	}

	FilterFuncs = template.FuncMap{
		"exclude": sockaddr.ExcludeIfs,
		"include": sockaddr.IncludeIfs,
	}

	HelperFuncs = template.FuncMap{
		// Misc functions that operate on IfAddrs inputs
		"join":     sockaddr.JoinIfAddrs,
		"limit":    sockaddr.LimitIfAddrs,
		"offset":   sockaddr.OffsetIfAddrs,
		"reverse":  sockaddr.ReverseIfAddrs,
		"uniqueBy": sockaddr.UniqueIfAddrsBy,

		// Return a Private RFC 6890 IP address string that is attached
		// to the default route.
		"GetPrivateIP": sockaddr.GetPrivateIP,

		// Return a Public RFC 6890 IP address string that is attached
		// to the default route.
		"GetPublicIP": sockaddr.GetPublicIP,
	}
}

// Parse parses input as template input using the addresses available on the
// host, then returns the string output if there are no errors.
func Parse(input string) (string, error) {
	addrs, err := sockaddr.GetAllInterfaces()
	if err != nil {
		return "", errwrap.Wrapf("unable to query interface addresses: {{err}}", err)
	}

	return ParseIfAddrs(input, addrs)
}

// ParseIfAddrs parses input as template input using the IfAddrs inputs, then
// returns the string output if there are no errors.
func ParseIfAddrs(input string, ifAddrs sockaddr.IfAddrs) (string, error) {
	return ParseIfAddrsTemplate(input, ifAddrs, template.New("sockaddr.Parse"))
}

// ParseIfAddrsTemplate parses input as template input using the IfAddrs inputs,
// then returns the string output if there are no errors.
func ParseIfAddrsTemplate(input string, ifAddrs sockaddr.IfAddrs, tmplIn *template.Template) (string, error) {
	// Create a template, add the function map, and parse the text.
	tmpl, err := tmplIn.Option("missingkey=error").
		Funcs(SourceFuncs).
		Funcs(SortFuncs).
		Funcs(FilterFuncs).
		Funcs(HelperFuncs).
		Parse(input)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("unable to parse template %+q: {{err}}", input), err)
	}

	var outWriter bytes.Buffer
	err = tmpl.Execute(&outWriter, ifAddrs)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("unable to execute sockaddr input %+q: {{err}}", input), err)
	}

	return outWriter.String(), nil
}

// SortByAddrs takes an array of SockAddrs and orders them by address.
// SockAddrs that are not comparable will be at the end of the list, however
// their order is non-deterministic.
func SortByAddrs(inputIfAddrs sockaddr.IfAddrs) sockaddr.IfAddrs {
	sortedIfAddrs := append(sockaddr.IfAddrs(nil), inputIfAddrs...)
	sockaddr.OrderedIfAddrBy(sockaddr.AscIfAddress).Sort(sortedIfAddrs)
	return sortedIfAddrs
}
