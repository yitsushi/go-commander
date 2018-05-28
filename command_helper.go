package commander

import (
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
		FmtPrintf("[Debug] %s\n", message)
	}
}

// Arg return with an item from Flags based on the given index
// emtpy string if not exists
func (c *CommandHelper) Arg(index int) string {
	if len(c.Args) > index {
		return c.Args[index]
	}

	return ""
}

// Flag return with an item from Flags based on the given key
// false if not exists
func (c *CommandHelper) Flag(key string) bool {
	if value, ok := c.Flags[key]; ok {
		return value
	}

	return false
}

// Opt return with an item from Opts based on the given key
// empty string if not exists
func (c *CommandHelper) Opt(key string) string {
	if value, ok := c.Opts[key]; ok {
		return value
	}

	return ""
}

// Parse is a helper method that parses all passed arguments
// flags, opts and arguments
func (c *CommandHelper) Parse(flag []string) {
	c.Flags = map[string]bool{}
	c.Opts = map[string]string{}

	if len(flag) < 2 {
		return
	}

	arguments := flag[1:]
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
