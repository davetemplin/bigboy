package main

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	actual := fileExists("README.md")

	assert.True(t, actual, "Checks if file exists")
}

func TestFileExistsAcceptsSubDirectories(t *testing.T) {
	actual := fileExists(".vscode/launch.json")

	assert.True(t, actual, "Accepts Subdirectories")
}

func TestFileExistsNotExist(t *testing.T) {
	actual := fileExists("NOT_EXIST")

	assert.False(t, actual, "Returns false")
}

func TestFileExistsIgnoresDirectories(t *testing.T) {
	actual := fileExists(".vscode")

	assert.False(t, actual, "Ignores Directories")
}

func TestFileExistsThrowsError(t *testing.T) {
	notDir := "LICENSE/NOT_EXIST"

	var actual bool
	if runtime.GOOS != "windows" {
		defer func() {
			if err := recover(); err != nil {
				actual = true
			}

			assert.True(t, actual, "Throws Error")
		}()
	}

	actualWindows := fileExists(notDir)
	if runtime.GOOS == "windows" {
		assert.False(t, actualWindows, "Doesn't throw an Error")
	}
}
