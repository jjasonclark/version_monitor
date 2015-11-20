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
	m := NewMonitor(slackURL)

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
