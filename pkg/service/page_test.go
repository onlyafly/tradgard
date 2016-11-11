package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	var input, actual, expected string

	input = "{Calvin} is a good man with beautiful {hair}."
	actual = transformWikiLinks(input)
	expected = "[Calvin](Calvin) is a good man with beautiful [hair](hair)."
	assert.Equal(t, expected, actual)

	input = "{Calvin is a good man} with beautiful {hair}."
	actual = transformWikiLinks(input)
	expected = "[Calvin is a good man](Calvin+is+a+good+man) with beautiful [hair](hair)."
	assert.Equal(t, expected, actual)
}
