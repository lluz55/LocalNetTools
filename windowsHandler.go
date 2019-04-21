// +build windows

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

// Add key to windowsRegistry
func _addToStartup() {

	startWithOsFlagged = true

	// // Windows
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}

	defer k.Close()

	currentAppName := os.Args[0]

	if err = k.SetStringValue(appName, fmt.Sprintf("%s %s", currentAppName, startedByOS)); err != nil {
		panic(err)
	}
}

// Remove key from windowsRegistry
func _removeFromStartup() {
	// Windows
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}

	defer k.Close()

	if err = k.DeleteValue(appName); err != nil {
		panic(err)
	}
}

func _startedWithOS() {
	// Windows
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		panic(err)
	}

	defer k.Close()

	_, _, err = k.GetStringValue(appName)
	if err == nil {
		startWithOsFlagged = true

	}
}

