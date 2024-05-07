package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestGetUserCrudSelection(t *testing.T) {
	// Test cases
	testCases := []struct {
		Input          string
		ExpectedOutput CrudOption
	}{
		{"1", READ},
		{"2", WRITE},
		{"3", DELETE},
	}

	for _, tc := range testCases {
		// Set up a mock reader with the input
		reader := bufio.NewReader(strings.NewReader(tc.Input + "\n"))

		// Call the function
		result := getUserCrudSelection(reader)

		// Check the result
		if result != tc.ExpectedOutput {
			t.Errorf("Expected %s, but got %s", tc.ExpectedOutput, result)
		}
	}
}

func TestGetUserLanguageSelection(t *testing.T) {
	// Test cases
	testCases := []struct {
		Input          string
		ExpectedOutput string
	}{
		{"1\n", "English"},
		{"2\n", "Spanish"},
		{"3\n", "French"},
	}

	optionsList := []LanguageOption{
		{Number: 1, LanguageName: "English"},
		{Number: 2, LanguageName: "Spanish"},
		{Number: 3, LanguageName: "French"},
	}

	for _, tc := range testCases {
		// Set up a mock reader with the input
		reader := bufio.NewReader(strings.NewReader(tc.Input))

		// Call the function
		result := getUserLanguageSelection(optionsList, reader)

		// Check the result
		if result.LanguageName != tc.ExpectedOutput {
			t.Errorf("Expected %s, but got %s", tc.ExpectedOutput, result.LanguageName)
		}
	}
}
