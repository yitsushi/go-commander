package commander

import (
	"flag"
	"fmt"
	"strings"
)

// CommandHelper is a helper struct
// CommandHandler.Execute will get this as an argument
// and you can access extra functions, farsed flags with this
type CommandHelper struct {
	// If -d is defined
	DebugMode bool
	// If -v is defined
	VerboseMode bool
	// Boolean opts
	Flags map[string]bool
	// Other opts passed
	Opts map[string]string
	// Non-flag arguments
	Args []string
}

// Log is a logger function for debug messages
// it prints a message if DebugeMode is true
func (c *CommandHelper) Log(message string) {
	if c.DebugMode {
		fmt.Printf("[Debug] %s\n", message)
	}
}

// Parse is a helper method that parses all passed arguments
// flags, opts and arguments
func (c *CommandHelper) Parse() {
	c.Flags = map[string]bool{}
	c.Opts = map[string]string{}

	if len(flag.Args()) < 1 {
		return
	}

	arguments := flag.Args()[1:]
	for _, arg := range arguments {
		if len(arg) > 1 && arg[0:2] == "--" {
			parts := strings.SplitN(arg[2:], "=", 2)
			if len(parts) > 1 {
				// has exact value
				c.Opts[parts[0]] = parts[1]
			} else {
				c.Flags[parts[0]] = true
			}
			continue
		}

		if arg[0] == '-' {
			for _, o := range []byte(arg[1:]) {
				c.Flags[string(o)] = true
			}
			continue
		}

		c.Args = append(c.Args, arg)
	}

	if c.Flags["d"] {
		c.DebugMode = true
	}

	if c.Flags["v"] {
		c.VerboseMode = true
	}
}
