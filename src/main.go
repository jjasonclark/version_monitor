package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fmt.Printf("%s\nVersion: %s SHA %s\n\n", AppName, Version, BuildSHA)

	resp, err := http.Get("http://www.par8o.com/revision")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("https://github.com/par8o/par8o/commit/%s\n", body)
}
