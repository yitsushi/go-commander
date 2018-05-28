package commander

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var executeCalled bool

// Start: Commands
// Simple command
type MyCommand struct {
}

func (c *MyCommand) Execute(opts *CommandHelper) {
	executeCalled = opts.VerboseMode

	if opts.Flag("fail-me") {
		panic("PANIC!!! PANIC!!! PANIC!!! Calm down, please!")
	}
}

// SubCommand system
type MySubCommand struct {
}

func (c *MySubCommand) Execute(opts *CommandHelper) {
	executeCalled = opts.VerboseMode
}

type MyMainCommand struct {
}

func (c *MyMainCommand) Execute(opts *CommandHelper) {
	registry := NewCommandRegistry()
	registry.Depth = 1
	registry.Register(func(appName string) *CommandWrapper {
		return &CommandWrapper{
			Handler: &MySubCommand{},
			Help: &CommandDescriptor{
				Name:             "my-subcommand",
				ShortDescription: "This is my own SubCommand",
				Arguments:        "",
			},
		}
	})
	registry.Execute()
}

// End: Commands

var mockOutput string

func mockPrintf(format string, n ...interface{}) (int, error) {
	mockOutput += fmt.Sprintf(format, n...)
	return 0, nil
}

func mockEverything() {
	OSExtExecutable = func() (string, error) {
		return "/some/random/path/my-executable", nil
	}

	mockOutput = ""
	FmtPrintf = mockPrintf
}

func TestCommandRegistry_executableName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{want: "my-executable"},
	}

	mockEverything()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommandRegistry()
			if got := c.executableName(); got != tt.want {
				t.Errorf("CommandRegistry.executableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandRegistry(t *testing.T) {
	tests := []struct {
		name     string
		cliArgs  []string
		commands []NewCommandFunc
		test     func(*CommandRegistry, string) string
	}{
		{
			name:     "No command, print help",
			cliArgs:  []string{},
			commands: []NewCommandFunc{},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				expected := "help [command]   Display this help or a command specific help\n"
				if output != expected {
					return fmt.Sprintf("output(%s), want(%s)", output, expected)
				}
				return
			},
		},
		{
			name:     "Help command, print help",
			cliArgs:  []string{"help"},
			commands: []NewCommandFunc{},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				expected := "help [command]   Display this help or a command specific help\n"
				if output != expected {
					return fmt.Sprintf("output(%s), want(%s)", output, expected)
				}
				return
			},
		},
		{
			name:     "No command, invalid command",
			cliArgs:  []string{"invalid-call"},
			commands: []NewCommandFunc{},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				value := "Command not found: invalid-call"
				if !strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) not found in output(%s)", value, output)
				}
				return
			},
		},
		{
			name:    "Register one command, but no command called",
			cliArgs: []string{},
			commands: []NewCommandFunc{
				func(appName string) *CommandWrapper {
					return &CommandWrapper{
						Handler: &MyCommand{},
						Help: &CommandDescriptor{
							Name:             "my-command",
							ShortDescription: "This is my own command",
							LongDescription: `This is a very long
description about this command.`,
							Arguments: "<filename> [optional-argument]",
							Examples: []string{
								"test.txt",
								"test.txt copy",
								"test.txt move",
							},
						},
					}
				},
			},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				value := "my-command <filename> [optional-argument]"
				if !strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) not found in output(%s)", value, output)
				}
				return
			},
		},
		{
			name:    "Register one command, call help for command",
			cliArgs: []string{"help", "my-command"},
			commands: []NewCommandFunc{
				func(appName string) *CommandWrapper {
					return &CommandWrapper{
						Handler: &MyCommand{},
						Help: &CommandDescriptor{
							Name:             "my-command",
							ShortDescription: "This is my own command",
							LongDescription: `This is a very long
description about this command.`,
							Arguments: "<filename> [optional-argument]",
							Examples: []string{
								"test.txt",
								"test.txt copy",
								"test.txt move",
							},
						},
					}
				},
			},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				values := []string{
					"Usage: my-executable my-command <filename> [optional-argument]",
					"This is a very long\ndescription about this command.",
					"my-executable my-command test.txt",
					"my-executable my-command test.txt copy",
					"my-executable my-command test.txt move",
				}
				for _, value := range values {
					if !strings.Contains(output, value) {
						return fmt.Sprintf("value(%s) not found in output(%s)", value, output)
					}
				}
				return
			},
		},
		{
			name:    "Register one command, call command",
			cliArgs: []string{"my-command", "-v"},
			commands: []NewCommandFunc{
				func(appName string) *CommandWrapper {
					return &CommandWrapper{
						Handler: &MyCommand{},
						Help: &CommandDescriptor{
							Name:             "my-command",
							ShortDescription: "This is my own command",
							LongDescription: `This is a very long
description about this command.`,
							Arguments: "<filename> [optional-argument]",
							Examples: []string{
								"test.txt",
								"test.txt copy",
								"test.txt move",
							},
						},
					}
				},
			},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				if !executeCalled {
					return "Command should be called with VerboseMode"
				}
				return
			},
		},
		{
			name:    "Main and SubCommand, no arg",
			cliArgs: []string{},
			commands: []NewCommandFunc{
				func(appName string) *CommandWrapper {
					return &CommandWrapper{
						Handler: &MyCommand{},
						Help: &CommandDescriptor{
							Name:             "my-command",
							ShortDescription: "This is my own MainCommand",
							Arguments:        "<subcommand>",
						},
					}
				},
			},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				var value string

				value = "my-command <subcommand>"
				if !strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) not found in output(%s)", value, output)
				}

				value = "my-subcommand"
				if strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) found in output(%s)", value, output)
				}
				return
			},
		},
		{
			name:    "Main and SubCommand, help my-command",
			cliArgs: []string{"help", "my-command"},
			commands: []NewCommandFunc{
				func(appName string) *CommandWrapper {
					return &CommandWrapper{
						Handler: &MyMainCommand{},
						Help: &CommandDescriptor{
							Name:             "my-command",
							ShortDescription: "This is my own MainCommand",
							Arguments:        "<subcommand>",
						},
					}
				},
			},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				var value string

				value = "Usage: my-executable my-command <subcommand>"
				if !strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) not found in output(%s)", value, output)
				}

				value = "my-subcommand"
				if strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) found in output(%s)", value, output)
				}
				return
			},
		},
		{
			name:    "Main and SubCommand, my-command without arg",
			cliArgs: []string{"my-command"},
			commands: []NewCommandFunc{
				func(appName string) *CommandWrapper {
					return &CommandWrapper{
						Handler: &MyMainCommand{},
						Help: &CommandDescriptor{
							Name:             "my-command",
							ShortDescription: "This is my own MainCommand",
							Arguments:        "<subcommand>",
						},
					}
				},
			},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				var value string

				value = "my-subcommand"
				if !strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) not found in output(%s)", value, output)
				}

				value = "my-command"
				if strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) found in output(%s)", value, output)
				}
				return
			},
		},
		{
			name:    "Main and SubCommand, my-command help",
			cliArgs: []string{"my-command", "help"},
			commands: []NewCommandFunc{
				func(appName string) *CommandWrapper {
					return &CommandWrapper{
						Handler: &MyMainCommand{},
						Help: &CommandDescriptor{
							Name:             "my-command",
							ShortDescription: "This is my own MainCommand",
							Arguments:        "<subcommand>",
						},
					}
				},
			},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				var value string

				value = "my-subcommand"
				if !strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) not found in output(%s)", value, output)
				}

				value = "my-command"
				if strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) found in output(%s)", value, output)
				}
				return
			},
		},
		{
			name:    "Main and SubCommand, my-command help my-subcommand",
			cliArgs: []string{"my-command", "help", "my-subcommand"},
			commands: []NewCommandFunc{
				func(appName string) *CommandWrapper {
					return &CommandWrapper{
						Handler: &MyMainCommand{},
						Help: &CommandDescriptor{
							Name:             "my-command",
							ShortDescription: "This is my own MainCommand",
							Arguments:        "<subcommand>",
						},
					}
				},
			},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				var value string

				value = "Usage: my-executable my-command my-subcommand"
				if !strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) not found in output(%s)", value, output)
				}
				return
			},
		},
		{
			name:    "Register one command, and fail it",
			cliArgs: []string{"my-command", "-v", "--fail-me"},
			commands: []NewCommandFunc{
				func(appName string) *CommandWrapper {
					return &CommandWrapper{
						Handler: &MyCommand{},
						Help: &CommandDescriptor{
							Name:             "my-command",
							ShortDescription: "This is my own command",
						},
					}
				},
			},
			test: func(r *CommandRegistry, output string) (errMsg string) {
				value := "[E] PANIC!!! PANIC!!! PANIC!!! Calm down, please!"
				if !strings.Contains(output, value) {
					return fmt.Sprintf("value(%s) not found in output(%s)", value, output)
				}

				if !executeCalled {
					return "Command should be called with VerboseMode"
				}

				return
			},
		},
	}
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Pre-boot
			os.Args = append([]string{"/some/random/path/my-executable"}, tt.cliArgs...)
			executeCalled = false

			// Boot
			c := NewCommandRegistry()
			for _, command := range tt.commands {
				c.Register(command)
			}

			mockOutput = ""
			c.Execute()

			errMsg := tt.test(c, mockOutput)
			if errMsg != "" {
				t.Error(errMsg)
			}
		})
	}
	os.Args = oldArgs
}
