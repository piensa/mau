package cmd

import "github.com/sbstjn/hanu"

func init() {
	Register(
		"list",
		"List open pull requests.",
		func(conv hanu.ConversationInterface) {
			conv.Reply("Open Pull Requests:\n`nationality` 74/100 coverage: https://bloblo.api.geosure.tech https://github.com/piensa/geosure/pulls/54 \n")
		},
	)
}
