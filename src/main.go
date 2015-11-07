package main

import (
	"fmt"
	"log"
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

func processResults(results <-chan checkResult, quit <-chan bool) {
	for {
		select {
		case r := <-results:
			if r.Err != nil {
				fmt.Printf("%s: Error: %s\n", r.Name, r.Err)
			} else {
				fmt.Printf("%s: %s\n", r.Name, r.Result)
			}
		case <-quit:
			return
		}
	}
}

func pollVersions(checks []Check, interval time.Duration, quitChannel <-chan bool) {
	fmt.Printf("Polling versions every %s\n", interval)

	resultsChannel := make(chan checkResult, 10)
	go processResults(resultsChannel, quitChannel)
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

func main() {
	fmt.Println(versionDisplay())

	if err := loadConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s\n", err)
	}

	quitChannel := make(chan bool)
	listenForKillSignal(quitChannel)
	pollVersions(getChecks(), getInterval(), quitChannel)
}
