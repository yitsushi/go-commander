package commander

// CommandDescriptor describes a command for Help calls
type CommandDescriptor struct {
	// Required! name of the command
	Name string
	// Optional: argument list as a string
	// Basic convention: <required_argument> [optional_argument]
	Arguments string
	// Optional: Short description is used in general help
	ShortDescription string
	// Optional: Long description is used in command specific help
	LongDescription string
	// Optional: Examples array is used in command specific help
	Examples []string
}
