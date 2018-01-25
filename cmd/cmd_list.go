package cmd

import "github.com/sbstjn/hanu"

func init() {
	Register(
		"list",
		"List open pull requests.",
		func(conv hanu.ConversationInterface) {
			conv.Reply("Todo es una farsa: \n Open Pull Requests:\n`mdf20180110` 95/100 coverage: https://tbypidhfto.api.geosure.tech/ https://github.com/piensa/geosure/pulls/11 \n")
		},
	)
}
