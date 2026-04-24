//go:build windows

package main

import "golang.org/x/sys/windows"

func init() {
	// Set console input and output to UTF-8 (codepage 65001)
	// Required for box-drawing characters and Unicode output.
	windows.SetConsoleCP(65001)
	windows.SetConsoleOutputCP(65001)
}
