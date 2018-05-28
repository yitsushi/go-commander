package commander

import (
	"fmt"
	"testing"
)

func TestCommandHelper_Opt(t *testing.T) {
	type fields struct {
		DebugMode   bool
		VerboseMode bool
		Flags       map[string]bool
		Opts        map[string]string
		Args        []string
	}
	tests := []struct {
		name   string
		fields fields
		key    string
		want   string
	}{
		{
			name: "Key found",
			fields: fields{
				Opts: map[string]string{"file": "something.txt"},
			},
			key:  "file",
			want: "something.txt",
		},
		{
			name: "Key not found",
			fields: fields{
				Opts: map[string]string{"file": "something.txt"},
			},
			key:  "files",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandHelper{
				DebugMode:   tt.fields.DebugMode,
				VerboseMode: tt.fields.VerboseMode,
				Flags:       tt.fields.Flags,
				Opts:        tt.fields.Opts,
				Args:        tt.fields.Args,
			}
			if got := c.Opt(tt.key); got != tt.want {
				t.Errorf("CommandHelper.Opt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandHelper_Flag(t *testing.T) {
	type fields struct {
		DebugMode   bool
		VerboseMode bool
		Flags       map[string]bool
		Opts        map[string]string
		Args        []string
	}
	tests := []struct {
		name   string
		fields fields
		key    string
		want   bool
	}{
		{
			name: "Key found",
			fields: fields{
				Flags: map[string]bool{"force": true},
			},
			key:  "force",
			want: true,
		},
		{
			name: "Key not found",
			fields: fields{
				Opts: map[string]string{},
			},
			key:  "force",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandHelper{
				DebugMode:   tt.fields.DebugMode,
				VerboseMode: tt.fields.VerboseMode,
				Flags:       tt.fields.Flags,
				Opts:        tt.fields.Opts,
				Args:        tt.fields.Args,
			}
			if got := c.Flag(tt.key); got != tt.want {
				t.Errorf("CommandHelper.Flag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandHelper_Arg(t *testing.T) {
	type fields struct {
		DebugMode   bool
		VerboseMode bool
		Flags       map[string]bool
		Opts        map[string]string
		Args        []string
	}
	tests := []struct {
		name   string
		fields fields
		index  int
		want   string
	}{
		{
			name: "First item",
			fields: fields{
				Args: []string{"first"},
			},
			index: 0,
			want:  "first",
		},
		{
			name: "Out of index",
			fields: fields{
				Args: []string{"first"},
			},
			index: 1,
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandHelper{
				DebugMode:   tt.fields.DebugMode,
				VerboseMode: tt.fields.VerboseMode,
				Flags:       tt.fields.Flags,
				Opts:        tt.fields.Opts,
				Args:        tt.fields.Args,
			}
			if got := c.Arg(tt.index); got != tt.want {
				t.Errorf("CommandHelper.Arg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandHelper_Parse(t *testing.T) {
	tests := []struct {
		name string
		flag []string
		test func(*CommandHelper) string
	}{
		{
			name: "no error without args",
			flag: []string{},
			test: func(c *CommandHelper) (errMsg string) {
				return
			},
		},
		{
			name: "simple argument",
			flag: []string{"command", "my_param"},
			test: func(c *CommandHelper) (errMsg string) {
				value := "my_param"
				if c.Arg(0) != value {
					return fmt.Sprintf(
						"Argument not found. Want(%s) : Got(%s)",
						value,
						c.Arg(0),
					)
				}
				return
			},
		},
		{
			name: "double dash",
			flag: []string{"command", "--file=something.txt"},
			test: func(c *CommandHelper) (errMsg string) {
				value := "something.txt"
				if c.Opt("file") != value {
					return fmt.Sprintf(
						"Option not found. Want(%s) : Got(%s)",
						value,
						c.Opt("file"),
					)
				}
				return
			},
		},
		{
			name: "double dash flag",
			flag: []string{"command", "--force"},
			test: func(c *CommandHelper) (errMsg string) {
				if c.Opt("force") != "" {
					return fmt.Sprintf(
						"Option not found. Want(%s) : Got(%s)",
						"",
						c.Opt("force"),
					)
				}
				if !c.Flag("force") {
					return "Force Flag is false, but we expect true."
				}
				return
			},
		},
		{
			name: "single dash flag",
			flag: []string{"command", "-f"},
			test: func(c *CommandHelper) (errMsg string) {
				if c.Opt("f") != "" {
					return fmt.Sprintf(
						"Option not found. Want(%s) : Got(%s)",
						"",
						c.Opt("f"),
					)
				}
				if !c.Flag("f") {
					return "'f' Flag is false, but we expect true."
				}
				return
			},
		},
		{
			name: "all together",
			flag: []string{"command", "-f", "--file=something.txt", "simple_arg"},
			test: func(c *CommandHelper) (errMsg string) {
				value := "simple_arg"
				if c.Arg(0) != value {
					return fmt.Sprintf(
						"Argument not found. Want(%s) : Got(%s)",
						value,
						c.Arg(0),
					)
				}

				if c.Opt("f") != "" {
					return fmt.Sprintf(
						"Option not found. Want(%s) : Got(%s)",
						"",
						c.Opt("f"),
					)
				}
				if !c.Flag("f") {
					return "'f' Flag is false, but we expect true."
				}

				value = "something.txt"
				if c.Opt("file") != value {
					return fmt.Sprintf(
						"Option not found. Want(%s) : Got(%s)",
						"",
						c.Opt("file"),
					)
				}
				if c.Flag("file") {
					return "'f' Flag is true, but we expect false."
				}
				return
			},
		},
		{
			name: "enable debug",
			flag: []string{"command", "-d"},
			test: func(c *CommandHelper) (errMsg string) {
				if c.Opt("d") != "" {
					return fmt.Sprintf(
						"Option not found. Want(%s) : Got(%s)",
						"",
						c.Opt("d"),
					)
				}
				if !c.Flag("d") {
					return "'d' Flag is false, but we expect true."
				}
				if !c.DebugMode {
					return "'c.DebugMode' is false, but we expect true."
				}
				return
			},
		},
		{
			name: "enable verbose",
			flag: []string{"command", "-v"},
			test: func(c *CommandHelper) (errMsg string) {
				if c.Opt("v") != "" {
					return fmt.Sprintf(
						"Option not found. Want(%s) : Got(%s)",
						"",
						c.Opt("v"),
					)
				}
				if !c.Flag("v") {
					return "'v' Flag is false, but we expect true."
				}
				if !c.VerboseMode {
					return "'c.VerboseMode' is false, but we expect true."
				}
				return
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandHelper{}
			c.Parse(tt.flag)
			errMsg := tt.test(c)
			if errMsg != "" {
				t.Errorf("CommandHelper.Parse() => %s", errMsg)
				t.Errorf("%v", c)
			}
		})
	}
}
