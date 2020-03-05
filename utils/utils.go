package utils

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorHandler func
func ErrorHandler(err error) error {
	rerr := fmt.Sprintf("[ERROR]: %s", err)
	return errors.New(rerr)
}

// RemoveQuotes func
func RemoveQuotes(s string) string {
	s = strings.Replace(s, "\"", "", -1)
	s = strings.TrimRight(s, "\n")
	return s
}
