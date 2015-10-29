package main

import "fmt"

var AppName string   // Set from build
var Version string   // Set from build
var BuildTime string // Set from build

func main() {
	fmt.Printf("%s\nVersion: %s %s\n\n", AppName, Version, BuildTime)
}
