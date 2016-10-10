package commander

// CommandInterface defined a command.
// If a struct implements all the required function,
// it is acceptable as a Command for CommandRegistry
type CommandInterface interface {
	// Execute function will be executed when the command is called
	Execute()
	// Description is used to display the "one-liner" help message in genera help
	Description() string
	// ArgumentDescription is used in general help
	ArgumentDescription() string
	// Help output would be a human readable long command specific help message
	Help() string
	// Examples is an array of possible calls
	Examples() []string
}
