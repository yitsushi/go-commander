[![Go Report
Card](https://goreportcard.com/badge/github.com/Yitsushi/go-commander)](https://goreportcard.com/report/github.com/Yitsushi/go-commander)

This is a simple Go library to manage commands for your CLI tool.
Easy to use and now you can focus on Business Logic instead of building
the command routing.

### What this library does for you?

Manage your separated commands. How? Generates a general help and command
specific helps for your commands. If your command fails somewhere
(`panic` for example), commander will display the error message and
the command specific help to guide your user.

### Install

```
$ go get https://github.com/Yitsushi/go-commander
```

### Sample output _(from [totp-cli](https://github.com/Yitsushi/totp-cli))_

```
$ totp-cli help

change-password                   Change password
update                            Check and update totp-cli itself
version                           Print current version of this application
generate <namespace>.<account>    Generate a specific OTP
add-token [namespace] [account]   Add new token
list [namespace]                  List all available namespaces or accounts under a namespace
delete <namespace>[.account]      Delete an account or a whole namespace
help [command]                    Display this help or a command specific help
```

### Usage

Every single command has to implement `CommandInterface`.
Check [this project](https://github.com/Yitsushi/totp-cli) for examples.

```
package main

// Import the package
import "github.com/Yitsushi/go-commander"

// Your Command
type YourCommand struct {
}

// Executed only on command call
func (c *YourCommand) Execute() {
  // Command Action
}

// Argument list, only for help messages
// If you don't have any arguments, just return an empty string
// Basic convention: <registered_argument> [optional_argument]
func (c *YourCommand) ArgumentDescription() string {
  return "[name]"
}

// Help message for your command, only for help messages
// General Help
func (c *YourCommand) Description() string {
  return "This is my first command"
}

// Help message, long format.
// Command specific help
func (c *YourCommand) Help() string {
  return "This is a useless command, but at least I have one command"
}

// Examples are generated in this form:
//   loop through return array:
//     print: appname commandname item_in_this_array
func (c *YourCommand) Examples() []string {
  return []string{"", "test"}
}

// Main Section
func main() {
	registry := commander.NewCommandRegistry()

	registry.Register("your-command", &YourCommand{})

	registry.Execute()
}
```

Now you have a CLI tool with two commands: `help` and `your-command`.

```
❯ go build mytool.go

❯ ./mytool
your-command [name]   This is my first command
help [command]        Display this help or a command specific help

❯ ./mytool help your-command
Usage: mytool your-command [name]

This is a useless command, but at least I have one command

Examples:
  mytool your-command
  mytool your-command test
```
