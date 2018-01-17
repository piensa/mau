package cmd

import "github.com/sbstjn/hanu"

func init() {
	Register(
		"alias <hash:string> <environment:string>",
		"Alias a given environment to production, staging or development.",
		func(conv hanu.ConversationInterface) {
			hash, _ := conv.Match(0)
			environment, _ := conv.Match(1)
                        //TODO: check environment is one of development, staging, production.

                        // Update api.geosure.tech and v3.api.geosure.tech with the right value on openapi.yaml
                        // and commit it back to the main git repo.
			conv.Reply(environment + " environment is now running " + hash)
		},
	)
}
