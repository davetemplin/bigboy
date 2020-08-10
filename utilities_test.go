package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileExists(t *testing.T) {
	expected := true
	actual := fileExists("README.md")

	assert.Equal(t, expected, actual, "Returns true")
}

func TestFileExistsAcceptsSubDirectories(t *testing.T) {
	expected := true
	actual := fileExists(".vscode/launch.json")

	assert.Equal(t, expected, actual, "Accepts Subdirectories")
}

func TestFileExistsNotExist(t *testing.T) {
	expected := false
	actual := fileExists("NOT_EXIST")

	assert.Equal(t, expected, actual, "Returns false")
}

func TestFileExistsIgnoresDirectories(t *testing.T) {
	expected := false
	actual := fileExists(".vscode")

	assert.Equal(t, expected, actual, "Ignores Directories")
}

func TestFileExistsThrowsError(t *testing.T) {
	notDir := "LICENSE/NOT_EXIST"
	expected := true

	var actual bool
	defer func() {
		if err := recover(); err != nil {
			actual = true
		}

		assert.Equal(t, expected, actual, "Throws Error")
	}()

	fileExists(notDir)
}
