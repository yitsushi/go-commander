package commander

import (
	"fmt"

	"github.com/kardianos/osext"
)

// OSExtExecutable returns current executable path
var OSExtExecutable = osext.Executable

// FmtPrintf is fmt.Printf
var FmtPrintf = fmt.Printf

// NewCommandFunc is the expected type for CommandRegistry.Register
type NewCommandFunc func(appName string) *CommandWrapper

// ValidatorFunc can pre-validate the command and it's arguments
// Just throw a panic if something is wrong
type ValidatorFunc func(opts *CommandHelper)

// CommandWrapper is a general wrapper for a command
// CommandRegistry will know what to do this a struct like this
type CommandWrapper struct {
	// Help contains all information about the command
	Help *CommandDescriptor
	// Handler will be called when the user calls that specific command
	Handler CommandHandler
	// Validator will be executed before Execute on the Handler
	Validator ValidatorFunc
}
