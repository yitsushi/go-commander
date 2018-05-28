[![Go Report Card](https://goreportcard.com/badge/github.com/yitsushi/go-commander)](https://goreportcard.com/report/github.com/yitsushi/go-commander)
[![Build Status](https://travis-ci.org/yitsushi/go-commander.svg?branch=master)](https://travis-ci.org/yitsushi/go-commander)

This is a simple Go library to manage commands for your CLI tool.
Easy to use and now you can focus on Business Logic instead of building
the command routing.

### What this library does for you?

Manage your separated commands. How? Generates a general help and command
specific helps for your commands. If your command fails somewhere
(`panic` for example), commander will display the error message and
the command specific help to guide your user.

### Install

```shell
$ go get github.com/yitsushi/go-commander
```

### Sample output _(from [totp-cli](https://github.com/yitsushi/totp-cli))_

```shell
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

Every single command has to implement `CommandHandler`.
Check [this project](https://github.com/yitsushi/totp-cli) for examples.

```go
package main

// Import the package
import "github.com/yitsushi/go-commander"

// Your Command
type YourCommand struct {
}

// Executed only on command call
func (c *YourCommand) Execute(opts *commander.CommandHelper) {
  // Command Action
}

func NewYourCommand(appName string) *commander.CommandWrapper {
  return &commander.CommandWrapper{
    Handler: &YourCommand{},
    Help: &commander.CommandDescriptor{
      Name:             "your-command",
      ShortDescription: "This is my own command",
      LongDescription:  `This is a very long
description about this command.`,
      Arguments:        "<filename> [optional-argument]",
      Examples:         []string {
        "test.txt",
        "test.txt copy",
        "test.txt move",
      },
    },
  }
}

// Main Section
func main() {
	registry := commander.NewCommandRegistry()

	registry.Register(NewYourCommand)

	registry.Execute()
}
```

Now you have a CLI tool with two commands: `help` and `your-command`.

```bash
❯ go build mytool.go

❯ ./mytool
your-command <filename> [optional-argument]   This is my own command
help [command]                                Display this help or a command specific help

❯ ./mytool help your-command
Usage: mytool your-command <filename> [optional-argument]

This is a very long
description about this command.

Examples:
  mytool your-command test.txt
  mytool your-command test.txt copy
  mytool your-command test.txt move
```

#### How to use subcommand pattern?

When you create your main command, just create a new `CommandRegistry` inside
the `Execute` function like you did in your `main()` and change `Depth`.

```go
import subcommand "github.com/yitsushi/mypackage/command/something"

func (c *Something) Execute(opts *commander.CommandHelper) {
	registry := commander.NewCommandRegistry()
	registry.Depth = 1
	registry.Register(subcommand.NewSomethingMySubCommand)
	registry.Execute()
}
```

### PreValidation

If you want to write a general pre-validation for your command
or just simply keep your validation logic separated:

```go
// Or you can define inline if you want
func MyValidator(c *commander.CommandHelper) {
  if c.Arg(0) == "" {
    panic("File?")
  }

  info, err := os.Stat(c.Arg(0))
  if err != nil {
    panic("File not found")
  }

  if !info.Mode().IsRegular() {
    panic("It's not a regular file")
  }

  if c.Arg(1) != "" {
    if c.Arg(1) != "copy" && c.Arg(1) != "move" {
      panic("Invalid operation")
    }
  }
}

func NewYourCommand(appName string) *commander.CommandWrapper {
  return &commander.CommandWrapper{
    Handler: &YourCommand{},
    Validator: MyValidator
    Help: &commander.CommandDescriptor{
      Name:             "your-command",
      ShortDescription: "This is my own command",
      LongDescription:  `This is a very long
description about this command.`,
      Arguments:        "<filename> [optional-argument]",
      Examples:         []string {
        "test.txt",
        "test.txt copy",
        "test.txt move",
      },
    },
  }
}
```
