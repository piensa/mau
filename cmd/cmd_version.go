package cmd

import "github.com/sbstjn/hanu"

func init() {
	Register(
		"version",
		"Reply with the current chatbot version",
		func(conv hanu.ConversationInterface) {
			conv.Reply("I'm running with `%s`", Version)
		},
	)
}
