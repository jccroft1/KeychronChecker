package main

import (
	"flag"
	"fmt"

	"github.com/jccroft1/KeychronChecker/keychron"
	"github.com/jccroft1/KeychronChecker/telegram"
)

var (
	token   = flag.String("token", "", "the bot token")
	channel = flag.String("channel", "", "the target channel for the stock alert")
)

func main() {
	flag.Parse()

	telegram.Token = *token
	telegram.Channel = *channel

	err := telegram.GetMe()
	if err != nil {
		fmt.Println(err)
		return
	}

	keychron.Alert = telegram.SendMessage
	keychron.Start()
}
