package cmd

import (
	"github.com/sbstjn/hanu"
	"github.com/ghodss/yaml"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
	"strings"
	"regexp"
	"os"
	"io"
	"os/exec"
)


func init() {
	config_Prometheus := "/home/prometheus/data/"
	Register(
		"prometheus <hash:string>",
		"Set up config file in Prometheus to get the metrics",
		func(conv hanu.ConversationInterface) {
			hash, _ := conv.Match(0)
			msg:=SetupPrometheus(hash, config_Prometheus)
			if msg != "" {
			   conv.Reply("%s",msg)
			}
		},
	)
}

func ReloadPrometheus() {
    pid, err := exec.Command("pgrep","prometheus").Output()
    if err != nil {
        log.Fatal(err)
    }
    arg := strings.Replace(string(pid), "\n", "", -1)
    _ ,err = exec.Command("kill", "-HUP",arg).Output()

    if err != nil {
        log.Fatal(err)
    }

    _, err = exec.Command("curl","-X","POST","http://nunez.co:9090/-/reload").Output()
    if err != nil {
        log.Fatal(err)
    }
}

func SetupPrometheus (hash string, config_path string) string {

	var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
	if IsLetter(hash) == false {
		return "Hash must only contain letters. :worried:"
	}
	if len(hash) != 10 {
		return "Invalid hash lenght. It must be lenght 10. :worried:"
	}
	slice := []string{hash,".api.geosure.tech"}
	endpoint :=strings.Join(slice, "")
	s := []string{config_path,"config.yml"};
	yamlFile, err := ioutil.ReadFile(strings.Join(s, ""))
    if err != nil {
        log.Fatal("yamlFile.Get err   #%v ", err)
    }
    string_yaml := string(yamlFile)
    j2, err := yaml.YAMLToJSON([]byte(string_yaml))
    if err != nil {
    	log.Fatal("Yaml parsing error...")
    }
    var objmap map[string]interface{} 
	err = json.Unmarshal(j2, &objmap)
	if err != nil {
		log.Fatal("Json parsing error...")

	}
	scrape_configs := objmap["scrape_configs"].([]interface{})
	scrape := scrape_configs[0].(map[string]interface{})
	static_conf := scrape["static_configs"].([]interface{})
	static := static_conf[0].(map[string]interface{})
	targets := static["targets"].([]interface{})
	is_value := isValueInList(endpoint,targets)
	if is_value == true {
		return fmt.Sprintf("The endpoint `%s` has already been configured in Prometheus. Not need to update it. :v:", endpoint)
	}

	// adding new endpoint
	targets = append(targets,endpoint)
	static["targets"] = targets
	var static_to_inter []interface {}
	static_to_inter = append(static_to_inter,static)
	scrape["static_configs"] = static_to_inter
	var scrape_to_inter []interface {}
	scrape_to_inter = append(scrape_to_inter,scrape)
	objmap["scrape_configs"] = scrape_to_inter 
	save:= saveYML(config_path,objmap)
	if save == true {
		ReloadPrometheus()
		return fmt.Sprintf("The endpoint `%s` has updated. :smile:", endpoint)
	}
	fmt.Println("save: ",save)
	return ""

}

func saveYML (folder_path string, objmap map[string]interface{}) bool {
	json, err := json.Marshal(objmap)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
		}
	yml, err := yaml.JSONToYAML(json)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
		}
	string_yml := strings.Replace(string(yml), "null", "", -1)
	filepath := []string{folder_path,"config.yml"}
	fo, err := os.Create(strings.Join(filepath, ""))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
		}
	_, err = io.Copy(fo, strings.NewReader(string_yml))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
		}
	return true


}

func isValueInList(value string, list []interface{}) bool {
    for _, v := range list {
        if v == value {
            return true
        }
    }
    return false
}
