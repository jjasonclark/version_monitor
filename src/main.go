package main

import "fmt"

func main() {
	fmt.Printf("%s\nVersion: %s SHA %s\n\n", AppName, Version, BuildSHA)
}
