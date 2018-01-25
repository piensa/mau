package cmd

import "github.com/sbstjn/hanu"

func init() {
	Register(
		"accept <branch:string>",
		"Merges open Pull Request for specified branch onto master and redeploys master on the staging server.",
		func(conv hanu.ConversationInterface) {
			branch, _ := conv.Match(0)
			conv.Reply("Branch `" + branch + "` has been merged intop master and the staging server is now: " + "piknedixce.api.geosure.tech")
		},
	)
}
