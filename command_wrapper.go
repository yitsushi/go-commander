package commander

import (
	"fmt"

	"github.com/kardianos/osext"
)

// CommandWrapper is a general wrapper for a command
// CommandRegistry will know what to do this a struct like this
type CommandWrapper struct {
	// Help contains all information about the command
	Help *CommandDescriptor
	// Handler will be called when the user calls that specific command
	Handler CommandHandler
}

// NewCommandFunc is the expected type for CommandRegistry.Register
type NewCommandFunc func(appName string) *CommandWrapper

// OSExtExecutable returns current executable path
var OSExtExecutable = osext.Executable

// FmtPrintf is fmt.Printf
var FmtPrintf = fmt.Printf
