package commander

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
