package main

import (
	"log"
	"time"
	"github.com/sbstjn/hanu"
	"github.com/piensa/mau/cmd"
	"github.com/spf13/viper"
)

var (
	// Version is the bot version
	Version = "0.0.2"
	// SlackToken will be set by ENV or config file in init()
	SlackToken = ""
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/go/src/github.com/piensa/mau")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

         if err != nil {
              log.Fatal("Could not load config file")
          }

	SlackToken = viper.GetString("MAU_SLACK_TOKEN")
}

func main() {
	bot, err := hanu.New(SlackToken)

	if err != nil {
		log.Fatal(err)
	}

	cmd.Version = Version
	cmd.Start = time.Now()
	cmdList := cmd.List()
	for i := 0; i < len(cmdList); i++ {
		bot.Register(cmdList[i])
	}

	bot.Listen()
}
