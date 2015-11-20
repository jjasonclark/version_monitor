package main

import (
	"fmt"
	"net/http"
	"net/url"
)

type monitor struct {
	past     map[string]checkResult
	slackURL string
}

// NewMonitor creates a monitor
func NewMonitor(slackURL string) monitor {
	return monitor{make(map[string]checkResult), slackURL}
}

func (m *monitor) postSlack(msg string) error {

	data := url.Values{}
	data.Set("payload", msg)
	r, err := http.PostForm(m.slackURL, data)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (m *monitor) compareVersions(last, r checkResult) {
	if last.sameAs(r) {
		fmt.Printf("%s: Same version as last time %s\n", r.Name, fmt.Sprintf(r.Verify, r.Result))
	} else {
		m.formatMessage(r)
	}
}

func (m *monitor) formatMessage(r checkResult) {
	msg, err := r.output()
	if err != nil {
		m.reportToSlack(msg)
	} else {
		fmt.Printf("Error creating Slack message: %s\n", err)
	}
}

func (m *monitor) reportToSlack(msg string) {
	fmt.Printf("Slack message: %s\n", msg)
	if e := m.postSlack(msg); e != nil {
		fmt.Printf("Failed to post Slack message: %s\n", e)
	}
}

func (m *monitor) processResults(results <-chan checkResult, quit <-chan bool) {
	for {
		select {
		case <-quit:
			return
		case r := <-results:
			if r.Err != nil {
				fmt.Printf("%s: Error: %s\n", r.Name, r.Err)
			} else {
				last, ok := m.past[r.Name]
				m.past[r.Name] = r
				if ok {
					m.compareVersions(last, r)
				} else {
					fmt.Printf("%s: Initial version %s\n", r.Name, fmt.Sprintf(r.Verify, r.Result))
				}
			}
		}
	}
}
