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
		// pipeline.  `GetIfSockAddrs` is the default input and original
		// "dot" in the pipeline.
		//
		"GetIfSockAddrs": sockaddr.GetIfSockAddrs,
	}

	SortFuncs = template.FuncMap{
		// *sortBy* functions sort their IfAddrs
		//
		"sortByAddr": sockaddr.SortByAddr,
		"sortByPort": sockaddr.SortByPort,
		"sortByType": sockaddr.SortByType,
	}

	FilterFuncs = template.FuncMap{
		// The exclude* and include* functions filter IfAddrs
		//
		// *ByIfName filters by Interface.Name
		"excludeByIfName": sockaddr.IfByNameExclude,
		"includeByIfName": sockaddr.IfByNameInclude,

		// *ByType filters by address types
		"excludeByType": sockaddr.IfByTypeExclude,
		"includeByType": sockaddr.IfByTypeInclude,

		// *ByRFC filters by RFC status
		"excludeByRFC": sockaddr.IfByRFCExclude,
		"includeByRFC": sockaddr.IfByRFCInclude,

		// Extracts a set of attributes from IfAddrs
		//
		"ifAddrs": sockaddr.IfReturnAttrAddrs,
		"ifNames": sockaddr.IfReturnAttrNames,
	}

	HelperFuncs = template.FuncMap{
		// Misc functions that operate on []SockAddr inputs
		"joinAddrs":    sockaddr.JoinAddrs,
		"limitAddrs":   sockaddr.LimitAddrs,
		"reverseAddrs": sockaddr.ReverseAddrs,
	}
}

// Parse parses input as template input using the addresses available on the
// host, then returns the string output if there are no errors.
func Parse(input string) (string, error) {
	addrs, err := sockaddr.GetIfSockAddrs()
	if err != nil {
		return "", errwrap.Wrapf("unable to query interface addresses: {{err}}", err)
	}

	return ParseIfAddrs(input, addrs)
}

// ParseIfAddrs parses input as template input using the IfAddrs inputs, then
// returns the string output if there are no errors.
func ParseIfAddrs(input string, ifAddrs []sockaddr.IfAddrs) (string, error) {
	return ParseIfAddrsTemplate(input, ifAddrs, template.New("sockaddr.Parse"))
}

// ParseIfAddrsTemplate parses input as template input using the IfAddrs inputs,
// then returns the string output if there are no errors.
func ParseIfAddrsTemplate(input string, ifAddrs []sockaddr.IfAddrs, tmplIn *template.Template) (string, error) {
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
