package command

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/go-sockaddr/template"
	"github.com/mitchellh/cli"
)

type EvalCommand struct {
	Ui cli.Ui

	// debugOutput emits framed output vs raw output.
	debugOutput bool

	// flags is a list of options belonging to this command
	flags *flag.FlagSet

	// suppressNewline changes whether or not there's a newline between each
	// arg passed to the eval subcommand.
	suppressNewline bool
}

// Description is the long-form command help.
func (c *EvalCommand) Description() string {
	return `Parse the sockaddr template and evaluates the output.`
}

// Help returns the full help output expected by `sockaddr -h cmd`
func (c *EvalCommand) Help() string {
	return MakeHelp(c)
}

// InitOpts is responsible for setup of this command's configuration via the
// command line.  InitOpts() does not parse the arguments (see parseOpts()).
func (c *EvalCommand) InitOpts() {
	c.flags = flag.NewFlagSet("eval", flag.ContinueOnError)
	c.flags.Usage = func() { c.Ui.Output(c.Help()) }
	c.flags.BoolVar(&c.debugOutput, "d", false, "Debug output")
	c.flags.BoolVar(&c.suppressNewline, "n", false, "Suppress newlines between args")
}

// Run executes this command.
func (c *EvalCommand) Run(args []string) int {
	if len(args) == 0 {
		c.Ui.Error(c.Help())
		return 1
	}

	c.InitOpts()
	tmpls := c.parseOpts(args)
	inputs, output := []string{}, []string{}
	for i, in := range tmpls {
		if in == "-" {
			var f io.Reader = os.Stdin
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, f); err != nil {
				c.Ui.Error(fmt.Sprintf("[ERROR]: Error reading from stdin: %v", err))
				return 1
			}
			in = buf.String()
			if len(in) == 0 {
				return 0
			}
			inputs = append(inputs, in)
		}

		out, err := template.Parse(in)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("ERROR[%d] in: %q\n[%d] msg: %v\n", i, in, i, err))
			return 1
		}
		output = append(output, out)
	}

	if c.debugOutput {
		for i, out := range output {
			c.Ui.Output(fmt.Sprintf("[%d] in: %q\n[%d] out: %q\n", i, inputs[i], i, out))
			if i != len(output)-1 {
				if c.debugOutput {
					c.Ui.Output(fmt.Sprintf("---\n"))
				}
			}
		}
	} else {
		sep := "\n"
		if c.suppressNewline {
			sep = ""
		}
		c.Ui.Output(strings.Join(output, sep))
	}

	return 0
}

// Synopsis returns a terse description used when listing sub-commands.
func (c *EvalCommand) Synopsis() string {
	return `Evaluates a sockaddr template`
}

// Usage is the one-line usage description
func (c *EvalCommand) Usage() string {
	return `sockaddr eval [options] [template ...]`
}

// VisitAllFlags forwards the visitor function to the FlagSet
func (c *EvalCommand) VisitAllFlags(fn func(*flag.Flag)) {
	c.flags.VisitAll(fn)
}

// parseOpts is responsible for parsing the options set in InitOpts().  Returns
// a list of non-parsed flags.
func (c *EvalCommand) parseOpts(args []string) []string {
	if err := c.flags.Parse(args); err != nil {
		return nil
	}

	return c.flags.Args()
}
