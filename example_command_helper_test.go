package commander_test

import (
	"fmt"
	"log"

	commander "github.com/yitsushi/go-commander"
)

// (c *MyCommand) Execute(opts *commander.CommandHelper)
var opts *commander.CommandHelper

func init() {
	opts = &commander.CommandHelper{}
}

func ExampleCommandHelper_TypedOpt() {
	opts.AttachArgumentList([]*commander.Argument{
		&commander.Argument{
			Name: "list",
			Type: "StringArray[]",
		},
	})
	opts.Parse([]string{"my-command", "--list=one,two,three"})

	// list is a StringArray[]
	if opts.ErrorForTypedOpt("list") == nil {
		log.Println(opts.TypedOpt("list"))
		myList := opts.TypedOpt("list").([]string)
		if len(myList) > 0 {
			fmt.Printf("My list: %v\n", myList)
		}
	}

	// Never defined, always shoud be an empty string
	if opts.TypedOpt("no-key").(string) != "" {
		panic("Something went wrong!")
	}

	// Output: My list: [one two three]
}

func ExampleCommandHelper_Arg() {
	opts.Parse([]string{"my-command", "plain-argument"})

	fmt.Println(opts.Arg(0))
	// Output: plain-argument
}

func ExampleCommandHelper_Flag() {
	opts.Parse([]string{"my-command", "plain-argument", "-l", "--no-color"})

	if opts.Flag("l") {
		fmt.Println("-l is defined")
	}

	if opts.Flag("no-color") {
		fmt.Println("Color mode is disabled")
	}

	// Output: -l is defined
	// Color mode is disabled
}
