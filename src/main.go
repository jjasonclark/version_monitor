package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kardianos/osext"
	"github.com/spf13/viper"
)

func versionDisplay() string {
	return fmt.Sprintf("%s\nVersion: %s SHA %s\n", AppName, Version, BuildSHA)
}

var errFetchFail = errors.New("Failed to fetch SHA")

func fetchVersionSha(source, verify string) (string, error) {
	resp, err := http.Get(source)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errFetchFail
	}
	var body [40]byte
	n, err := resp.Body.Read(body[0:])
	if err != nil && err != io.EOF {
		return "", err
	}
	if n < 1 {
		return "", errFetchFail
	}
	return fmt.Sprintf(verify, string(body[:n])), nil
}

type check struct {
	Name   string
	Source string
	Verify string
}

type checkResult struct {
	check
	Err    error
	Result string
}

func (c *check) Check() *checkResult {
	result := &checkResult{*c, nil, ""}
	result.Result, result.Err = fetchVersionSha(result.Source, result.Verify)
	return result
}

func loadChecks() ([]check, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if sourcePWD, err := osext.ExecutableFolder(); err == nil {
		viper.AddConfigPath(sourcePWD)
	}
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var checks []check
	if err := viper.UnmarshalKey("checks", &checks); err != nil {
		return nil, err
	}
	return checks, nil
}

func main() {
	checks, err := loadChecks()
	if err != nil {
		log.Fatalf("Fatal error config file: %v \n", err)
	}

	fmt.Println(versionDisplay())

	for _, c := range checks {
		blah := c.Check()
		if blah.Err != nil {
			log.Fatal(blah.Err)
		}
		fmt.Printf("%s: %s\n", blah.Name, blah.Result)
	}
}
