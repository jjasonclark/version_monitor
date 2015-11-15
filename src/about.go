package main

import "fmt"

// AppName of application
var AppName = "Version_monitor"

// BuildSHA is the Git SHA of last commit used in build
var BuildSHA = "" // Updated from build

// Version of application
var Version = "0.0.1-pre1"

func versionDisplay() string {
	return fmt.Sprintf("%s\nVersion: %s SHA %s\n", AppName, Version, BuildSHA)
}
