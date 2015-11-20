package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

var errFetchFail = errors.New("Failed to fetch SHA")

func fetchVersionSha(source string) (string, error) {
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
	return string(body[:n]), nil
}

type checkResult struct {
	Check
	Err    error
	Result string
}

func (left checkResult) sameAs(right checkResult) bool {
	return strings.Compare(left.Result, right.Result) == 0
}

type Check struct {
	Name   string
	Source string
	Verify string
}

func (c Check) check(results chan<- checkResult) {
	sha, err := fetchVersionSha(c.Source)
	results <- checkResult{c, err, sha}
}

func queueAllChecks(checks []Check, resultsChannel chan<- checkResult) {
	for _, c := range checks {
		go c.check(resultsChannel)
	}
}
