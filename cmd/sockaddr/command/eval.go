package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-sockaddr/template"
	"github.com/mitchellh/cli"
)

type EvalCommand struct {
	Ui cli.Ui
}

func (c *EvalCommand) Help() string {
	helpText := `
Usage: sockaddr eval [template ...]

  Parse the sockaddr template and evaluates the output.
`
	return strings.TrimSpace(helpText)
}

func (c *EvalCommand) Run(args []string) int {
	if len(args) == 0 {
		c.Ui.Error(fmt.Sprintf("%s", c.Help()))
		return 1
	}

	for i, in := range args {
		out, err := template.Parse(in)
		if err != nil {
			return 1
		}
		c.Ui.Output(fmt.Sprintf("[%d] in: %q\n[%d] out: %q\n", i, in, i, out))
		if i != len(args)-1 {
			c.Ui.Output(fmt.Sprintf("---\n"))
		}
	}
	return 0
}

func (c *EvalCommand) Synopsis() string {
	return "Evaluates a sockaddr template"
}
