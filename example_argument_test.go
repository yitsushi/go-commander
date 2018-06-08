package commander_test

import (
	"errors"
	"strconv"
	"strings"

	commander "github.com/yitsushi/go-commander"
)

type MyCustomType struct {
	ID   uint64
	Name string
}

func ExampleRegisterArgumentType_simple() {
	commander.RegisterArgumentType("Int32", func(value string) (interface{}, error) {
		return strconv.ParseInt(value, 10, 32)
	})
}

func ExampleRegisterArgumentType_customStruct() {
	commander.RegisterArgumentType("MyType", func(value string) (interface{}, error) {
		values := strings.Split(value, ":")

		if len(values) < 2 {
			return &MyCustomType{}, errors.New("Invalid format! MyType => 'ID:Name'")
		}

		id, err := strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			return &MyCustomType{}, errors.New("Invalid format! MyType => 'ID:Name'")
		}

		return &MyCustomType{
				ID:   id,
				Name: values[1],
			},
			nil
	})
}
