package commander

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type argumentTypeFunction func(string) (interface{}, error)

var argumentTypeList map[string]argumentTypeFunction

// RegisterArgumentType registers a new argument type
func RegisterArgumentType(name string, f argumentTypeFunction) {
	argumentTypeList[name] = f
}

// Argument represents a single argument
type Argument struct {
	Name          string
	Type          string
	OriginalValue string
	Value         interface{}
	Error         error
	FailOnError   bool
}

// SetValue saves the original value to the argument.
// Returns with an error if conversion failed
func (a *Argument) SetValue(original string) error {
	a.OriginalValue = original
	a.Value, a.Error = argumentTypeList[a.Type](a.OriginalValue)

	return a.Error
}

func init() {
	argumentTypeList = map[string]argumentTypeFunction{}

	RegisterArgumentType("String", func(value string) (interface{}, error) {
		return value, nil
	})

	RegisterArgumentType("Int64", func(value string) (interface{}, error) {
		return strconv.ParseInt(value, 10, 64)
	})

	RegisterArgumentType("Uint64", func(value string) (interface{}, error) {
		return strconv.ParseUint(value, 10, 64)
	})

	RegisterArgumentType("StringArray[]", func(value string) (interface{}, error) {
		arr := strings.Split(value, ",")

		return arr, nil
	})

	RegisterArgumentType("FilePath", func(value string) (interface{}, error) {
		_, err := os.Stat(value)

		if os.IsNotExist(err) {
			log.Println(err)
			value = ""
		}
		return value, err
	})
}
