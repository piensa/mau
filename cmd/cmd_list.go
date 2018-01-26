package cmd

import (
	"github.com/sbstjn/hanu"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
)

func init() {
	ConfigFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config ConfigStruct
	json.Unmarshal(ConfigFile, &config)
	GithubToken := config.GithubToken

	Register(
		"list",
		"List open pull requests.",
		func(conv hanu.ConversationInterface) {
			pr := getPrList(GithubToken)
			fmt.Println(pr)
			answer := "No new features to list"

			if len(pr) > 0 {
				answer = "Open Pull Requests:\n"
				for i := 0; i < len(pr); i++ {
					answer = answer + "`" + pr[i].Head.Ref + "` " + pr[i].Html_url + " \n"
				}
			}
			conv.Reply(answer)
		},
	)
}

type Prlist struct {
	Url string `json:"url"`
	Html_url string `json:"html_url"`
	State string `json:"state"`
	Head Prhead `json:"head"`
}

type Prhead struct {
	Label string `json:"label"`
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}

type ConfigStruct struct {
	GithubToken string `json:"GITHUB_ACCESS_TOKEN"`
}

func getPrList(GithubToken string) []Prlist {
	url := "https://api.github.com/repos/geosure/geosure/pulls?access_token=" + GithubToken

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var pr []Prlist
	json.Unmarshal(body, &pr)

	return pr
}
