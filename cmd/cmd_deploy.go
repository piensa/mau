package cmd

import "github.com/sbstjn/hanu"

func init() {
	Register(
		"deploy <branch:string>",
		"Deploy a git branch or commit.",
		func(conv hanu.ConversationInterface) {
			branch, _ := conv.Match(0)
                        url := "https://blablah.geosure.tech"
			conv.Reply("Deployed " + branch + " at the following url " + url)
		},
	)
}
