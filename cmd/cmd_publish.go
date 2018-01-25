package cmd

import "github.com/sbstjn/hanu"

func init() {
	Register(
		"publish <hash:string>",
		"Promotes the latest staging server to production. Requires the instance hash as confirmation.",
		func(conv hanu.ConversationInterface) {
			hash, _ := conv.Match(0)
			conv.Reply("Production environment is now running " + hash)
		},
	)
}
