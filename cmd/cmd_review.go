package cmd

import (
	"github.com/sbstjn/hanu"
	"net/url"
	"net/http"
	"fmt"
	"log"
	"encoding/json"
//	"io/ioutil"
)

func init() {
	Register(
		"review",
		"Alegría al mejor precio!",
		func(conv hanu.ConversationInterface) {

			githubUrl := "https://api.github.com/repos/geosure/geosure/issues/29/comments?access_token=9e19598fa3a3433a1742b9fb24b21ade8e8536a4"
			body := `deployment:url=https://geosure-boohaaa.now.sh
				deployment:date=20180130204532
				test:coverage=94.5
				test:passed=true
				test:time=65`

			formData := url.Values{
				"body": {"body"},
			}

			resp, err := http.PostForm(githubUrl, formData)
			if err != nil {
				log.Fatalln(err)
			}

			var result map[string]interface{}

			json.NewDecoder(resp.Body).Decode(&result)

			fmt.Println(result["form"])

			conv.Reply(body)
		},
	)
}
