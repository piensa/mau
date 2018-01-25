package cmd

import "github.com/sbstjn/hanu"

func init() {
	Register(
		"servers",
		"List development, staging and production servers.",
		func(conv hanu.ConversationInterface) {
                        development := "fhgfdsagifu"
                        staging := "tutufddfjs"
                        production := "fvtoofdsafd"
			conv.Reply("The API servers are currently:\n*development*: " + development +"\n*staging*: " + staging + "\n*production*: " +
 production)
		},
	)
}
