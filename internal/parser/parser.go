package parser

import (
	"fmt"
	"strings"
)

func Parse(input string) ([]string, error) {
	split := strings.Split(strings.TrimSpace(input), "\r\n")

	if len(split) < 3 || !strings.HasPrefix(split[0], "*") || len(split[0]) < 2 {
		return nil, fmt.Errorf("invalid first element of RESP array: %q", GetRaw(input))
	}

	args := []string{}

	for i := 2; i < len(split); i++ {
		if !strings.HasPrefix(split[i], "$") {
			args = append(args, split[i])
		}
	}

	return args, nil
}

func GetCommand(input []string) string {
	switch strings.ToLower(input[0]) {
	case "ping":
		return "ping"
	case "echo":
		if len(input) == 2 {
			return "echo"
		}
	case "set":
		if len(input) == 3 {
			return "set"
		}
		if len(input) == 5 && input[3] == "px" {
			return "set_expiry"
		}
	case "get":
		if len(input) == 2 {
			return "get"
		}
	case "config":
		if len(input) != 3 {
			return ""
		}
		if strings.ToLower(input[1]) == "get" {
			return "config_get"
		}
	}
	return ""
}

func ToSimpleString(msg string) string {
	return "+" + msg + "\r\n"
}

func ToBulkString(msg string) string {
	return fmt.Sprint("$", len(msg), "\r\n", msg, "\r\n")
}

func GetRaw(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, "\r", "\\r"), "\n", "\\n")
}

func ToRESPArray(values []string) string {
	result := "*" + fmt.Sprintf("%d", len(values)) + "\r\n"
	for _, val := range values {
		result += ToBulkString(val)
	}
	return result
}
