package command

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-sockaddr/template"
	"github.com/mitchellh/cli"
)

type EvalCommand struct {
	Ui cli.Ui
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
	// noop, no flags to parse for this command
}

// Run executes this command.
func (c *EvalCommand) Run(args []string) int {
	if len(args) == 0 {
		c.Ui.Error(c.Help())
		return 1
	}

	for i, in := range args {
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
		}

		out, err := template.Parse(in)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("ERROR[%d] in: %q\n[%d] msg: %v\n", i, in, i, err))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("[%d] in: %q\n[%d] out: %q\n", i, in, i, out))
		if i != len(args)-1 {
			c.Ui.Output(fmt.Sprintf("---\n"))
		}
	}
	return 0
}

// Synopsis returns a terse description used when listing sub-commands.
func (c *EvalCommand) Synopsis() string {
	return `Evaluates a sockaddr template`
}

// Usage is the one-line usage description
func (c *EvalCommand) Usage() string {
	return `sockaddr eval [template ...]`
}

// VisitAllFlags forwards the visitor function to the FlagSet
func (c *EvalCommand) VisitAllFlags(func(*flag.Flag)) {
	// noop, no flags to parse for this command
}
