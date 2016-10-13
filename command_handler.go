package commander

// CommandHandler defines a command.
// If a struct implements all the required function,
// it is acceptable as a CommandHandler for CommandRegistry
type CommandHandler interface {
	// Execute function will be executed when the command is called
	// opts can be used for logging, parsing flags like '-v'
	Execute(opts *CommandHelper)
}
