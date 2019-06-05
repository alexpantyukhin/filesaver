package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestgetKeyboardFromNames_NamesLessThanButtonsPerRow_OneButtonRow(t *testing.T) {
	markup := getKeyboardFromNames([]string{"1", "2"}, 3)

	assert.Equal(t, 1, len(markup.InlineKeyboard))
}

func TestgetKeyboardFromNames_NamesEqualButtonsPerRow_OneButtonRow(t *testing.T) {
	markup := getKeyboardFromNames([]string{"1", "2", "3"}, 3)

	assert.Equal(t, 1, len(markup.InlineKeyboard))
}


func TestgetKeyboardFromNames_NamesGreaterThanButtonsPerRow_ManyButtonRows(t *testing.T) {
	markup := getKeyboardFromNames([]string{"1", "2", "3", "4"}, 3)

	assert.Equal(t, 2, len(markup.InlineKeyboard))
}

