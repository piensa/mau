package main

import (
	"fmt"
	"github.com/piensa/mau/cmd"
	"github.com/sbstjn/hanu"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/github"
	"log"
	"strconv"
	"time"
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

const (
	path = "/webhooks"
	port = 3016
)

// HandlePullRequest handles GitHub pull_request events
func HandlePullRequest(payload interface{}, header webhooks.Header) {

	fmt.Println("Handling Pull Request")

	pl := payload.(github.PullRequestPayload)

	// Do whatever you want from here...
	fmt.Printf("%+v", pl)
}

func github_main() {

	hook := github.New(&github.Config{Secret: "hola-ariel"})
	hook.RegisterEvents(HandlePullRequest, github.PullRequestEvent)

	err := webhooks.Run(hook, ":"+strconv.Itoa(port), path)
	if err != nil {
		fmt.Println(err)
	}
}

func bot_main() {
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

func main() {
	go github_main()
	bot_main()
}
