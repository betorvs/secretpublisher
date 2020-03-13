package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveQuotes(t *testing.T) {
	value := "\"value"
	test := RemoveQuotes(value)
	assert.Contains(t, test, "value")
}

func TestErrorHandler(t *testing.T) {
	err := fmt.Errorf("test")
	test := ErrorHandler(err)
	testString := fmt.Sprintf("%v", test)
	assert.Contains(t, testString, "[ERROR]")
}
