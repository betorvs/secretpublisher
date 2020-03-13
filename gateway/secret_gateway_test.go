package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateHeaderSignature(t *testing.T) {
	longString := string("2aeccc9c03b36fea59ebec69")
	bodyString := string("body")
	timestamp := "1580475458"
	test := createHeaderSignature(timestamp, bodyString, longString)
	assert.Contains(t, test, "v0=dd8c5752a")
}
