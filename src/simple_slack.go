package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
	"strings"
)

func main() {
	api := slack.New("xoxb-80733663043-WmG8TgDS3V1c24QmPpaIr9dN")
	logger := log.New(os.Stdout, "simple_slack: ", log.Lshortfile|log.LstdFlags)

	logger.Print("set up API")

	slack.SetLogger(logger)
	debug := os.Getenv("DEBUG")
	if debug == "yes" {
		api.SetDebug(true)
	} else {
		api.SetDebug(false)
	}
	params := slack.PostMessageParameters{AsUser:true}
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:

			logger.Print("msg: ", msg)

			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				containsGo := strings.Contains(strings.ToLower(ev.Msg.Text), "neil is awesome!")
				if containsGo {
					_, _, err := api.PostMessage(ev.Msg.Channel, "Neil is fighting with go!", params)
					if err != nil {
						fmt.Printf("%s\n", err)
						return
					}
				}
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials: %v\n", ev)
				break Loop

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			default:
			}
		}
	}
}
