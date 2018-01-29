package cmd

import (
	"github.com/sbstjn/hanu"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
	"strings"
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

			answer := "No new features to list"

			if len(pr) > 0 {
				answer = "Open Pull Requests:\n"
				for i := 0; i < len(pr); i++ {
					commentBody := getCommentBody(GithubToken, pr[i].Links.Comments.Href)
					answer = answer + "`" + pr[i].Head.Ref + "` " + pr[i].Html_url + " \n" + commentBody
				}
			}
			conv.Reply(answer)
		},
	)
}

type Prlist struct {
	Url		string	`json:"url"`
	Html_url	string	`json:"html_url"`
	State		string	`json:"state"`
	Links		PrLinks	`json:"_links"`
	Head		Prhead	`json:"head"`
}

type Prhead struct {
	Label		string	`json:"label"`
	Ref		string	`json:"ref"`
	Sha		string	`json:"sha"`
}

type ConfigStruct struct {
	GithubToken	string	`json:"GITHUB_ACCESS_TOKEN"`
}

type PrLinks struct {
	Comments	HrefStr	`json:"comments"`
}

type HrefStr struct {
	Href		string	`json:"href"`
}

type Comment struct {
	Body		string	`json:"body"`
}

type ParsedComment struct {
	DeploymentUrl	string
	DeploymentDate	string
	TestCoverage	string
	testPassed	string
	TestTime	string
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

func getCommentBody(GithubToken string, commentUrl string) string {
	url :=  commentUrl +  "?access_token=" + GithubToken

	resp, err :=  http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var comment []Comment
	json.Unmarshal(body, &comment)

	commentBody := ""

	for i := 0; i < len(comment); i++ {
		if strings.Contains(comment[i].Body, "deployment:url") {
			commentBody = comment[i].Body
		}
	}
	stringCommentParser(commentBody)
	return commentBody

}

func stringCommentParser(commentBody string) {
	commentArray := strings.Split(commentBody, "\n")
	fmt.Println(commentArray[0])
}
