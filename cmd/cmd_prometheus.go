package cmd

import (
	"github.com/sbstjn/hanu"
	"github.com/ghodss/yaml"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
	"strings"
)


func init() {
	config_Prometheus ="/Users/waybarrios/Documents/code/prometheus"
	Register(
		"prometheus <endpoint:string>",
		"Set up config file in Prometheus  to ge the metrics",
		func(conv hanu.ConversationInterface) {
			endpoint, _ := conv.Match(0)
			conv.Reply("Endpoint is `%s`", endpoint)
			conv.Reply("Prometheus config path is `%s`",config_Prometheus)
			parsing_yml(endpoint, config_Prometheus)
		},
	)
}

func parsing_yml (endpoint string, config_path string) {
	s := []string{config_path,"/","config.yml"};
	yamlFile, err := ioutil.ReadFile(strings.Join(s, ""))
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    string_yaml := string(yamlFile)
    j2, err := yaml.YAMLToJSON([]byte(string_yaml))
    if err != nil {
    	log.Printf("Yaml parsing error...")
    }
    var objmap map[string]interface{} 
	err = json.Unmarshal(j2, &objmap)
	if err != nil {
		log.Printf("Json parsing error...")
	}
	
	scrape_configs := objmap["scrape_configs"].([]interface{})
	scrape := scrape_configs[0].(map[string]interface{})
	static_conf := scrape["static_configs"].([]interface{})
	static := static_conf[0].(map[string]interface{})
	targets := static["targets"].([]interface{})

	for _, value := range targets {
		fmt.Println(value)
	}
	
    	
}