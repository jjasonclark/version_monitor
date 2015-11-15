package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func listenForKillSignal(quitChannel chan<- bool) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	go func() {
		sig := <-signalChannel
		fmt.Printf("\nExiting from signal: %s\n", sig)
		quitChannel <- true
	}()
}

type monitor struct {
	past     map[string]checkResult
	slackURL string
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

func (m *monitor) processResults(results <-chan checkResult, quit <-chan bool) {
	for {
		select {
		case r := <-results:
			if r.Err != nil {
				fmt.Printf("%s: Error: %s\n", r.Name, r.Err)
			} else {
				last, ok := m.past[r.Name]
				m.past[r.Name] = r
				if ok {
					if last.sameAs(r) {
						fmt.Printf("%s: Same version as last time %s\n", r.Name, fmt.Sprintf(r.Verify, r.Result))
					} else {
						msg, err := r.output()
						if err != nil {
							fmt.Printf("Slack message: %s\n", msg)
							if e := m.postSlack(msg); e != nil {
								fmt.Printf("Failed to post Slack message: %s\n", e)
							}
						} else {
							fmt.Printf("Error creating Slack message: %s\n", err)
						}
					}
				} else {
					fmt.Printf("%s: Initial version %s\n", r.Name, fmt.Sprintf(r.Verify, r.Result))
				}
			}
		case <-quit:
			return
		}
	}
}

func main() {
	fmt.Println(versionDisplay())

	if err := loadConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s\n", err)
	}

	quitChannel := make(chan bool)
	listenForKillSignal(quitChannel)

	checks := getChecks()
	interval := getInterval()
	slackURL := getSlackURL()
	m := monitor{make(map[string]checkResult), slackURL}

	fmt.Printf("Polling versions every %s\n", interval)

	resultsChannel := make(chan checkResult, 10)
	go m.processResults(resultsChannel, quitChannel)
	for {
		queueAllChecks(checks, resultsChannel)
		select {
		case <-time.After(interval):
			fmt.Println()
		case <-quitChannel:
			return
		}
	}
}
