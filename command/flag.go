package command

import (
	"errors"
	"fmt"
)

// isFlag checks whether the given argument is a valid flag or not
func isFlag(arg string) bool {
	if _, err := parseFlag(arg); err != nil {
		return false
	}

	return true
}

// parseFlag parses a flags name. A flag can be in form of --name=value,
// -name=value, -n=value, or --name, -name=, etc...  If it's a correct flag,
// the name is returned. If not an empty string and an error message is
// returned
func parseFlag(arg string) (string, error) {
	if arg == "" {
		return "", errors.New("argument is empty")
	}

	if len(arg) == 1 {
		return "", errors.New("argument is too short")
	}

	if arg[0] != '-' {
		return "", errors.New("argument doesn't start with dash")
	}

	numMinuses := 1

	if arg[1] == '-' {
		numMinuses++
		if len(arg) == 2 {
			return "", errors.New("argument is too short")
		}
	}

	name := arg[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return "", fmt.Errorf("bad flag syntax: %s", arg)
	}

	return name, nil
}

// parseValue parses the value from the given flag. A flag name can be in
// form of name=value, n=value, n=, n.
func parseValue(flag string) (name, value string) {
	for i, r := range flag {
		if r == '=' {
			value = flag[i+1:]
			name = flag[0:i]
		}
	}

	// special case of "n"
	if name == "" {
		name = flag
	}

	return
}

// parseName parses the given flagName from the args slice and returns the
// value passed to the flag. An example: args: ["--provider", "aws", "--foo"],
// flagName: "provider" will return "aws".
func parseFlagValue(flagName string, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("argument is empty")
	}

	for i, arg := range args {
		flag, err := parseFlag(arg)
		if err != nil {
			fmt.Println("err")
			continue
		}

		name, value := parseValue(flag)
		if name != flagName && name[0] != flagName[0] {
			continue
		}

		if value != "" {
			return value, nil // found
		}

		// no value found, check out the next argument. at least two args must
		// be present
		if len(args) > 1 {
			// value must be next argument
			if !isFlag(args[i+1]) {
				return args[i+1], nil
			}
		}
	}

	return "", errors.New("couldn't find")
}
