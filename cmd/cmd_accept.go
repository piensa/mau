package cmd

import(
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
    "io"
    "github.com/sbstjn/hanu"
    "strings"
    "github.com/ghodss/yaml"
    "fmt"
    "os/exec"
)

func init() {
	Register(
		"accept <branch:string>",
		"Merges open Pull Request for specified branch onto master and redeploys master on the staging server.",
		func(conv hanu.ConversationInterface) {
            test_path = "/home/prometheus/test"
			branch, _ := conv.Match(0)
            msg,stage_url := AcceptPull(branch,test_path)
            conv.Reply("Your request is in progress... Waiting... :stopwatch:")
            if msg != ""{
                conv.Reply(msg)
            } else {
                conv.Reply("Branch `" + branch + "` has been merged into master and the staging server is now: " + stage_url)
            }	
			
		},
	)
}

func HubCheckout (branch string, test_path string) error {
    cmd := exec.Command("hub","checkout", branch)
    cmd.Dir = test_path
    _, err := cmd.Output()
    if err != nil {
         return err
        }
    return nil    
}
func GitCommit (test_path string, file string, msg string) error {
    cmd := exec.Command("git","add", file)
    cmd.Dir = test_path
    _, err := cmd.Output()
    if err != nil {
         return err
    }
    cmd_commit := exec.Command("git","commit","-m",msg)
    cmd_commit.Dir = test_path
    _, err_commit := cmd_commit.Output()
    if err_commit != nil {
         return err_commit
    }
    return nil

}
func HubPush (test_path string) error {
    cmd:= exec.Command("hub","push")
    cmd.Dir = test_path
     _, err := cmd.Output()
    if err != nil {
         return err
        }
    return nil 
}
func MergeBranch (branch string, test_path string) error {
    branch_ori:=fmt.Sprintf("origin/%s",branch)
    cmd := exec.Command("hub","merge",branch_ori)
    cmd.Dir = test_path
    _, err := cmd.Output()
    fmt.Println(err)
    if err != nil {
         return err
        }
    return nil 
}

func EditOpenApi(test_path string,stage_url string) map[string]interface{} {
    s := []string{test_path,"/api/openapi.yaml"};
    yamlFile, err := ioutil.ReadFile(strings.Join(s, ""))
    if err != nil { log.Fatal("yamlFile.Get err   #%v ", err)}
    string_yaml := string(yamlFile)
    j2, err := yaml.YAMLToJSON([]byte(string_yaml))
    if err != nil { log.Fatal("Yaml parsing error...")}
    var obj_map map[string]interface{} 
    err = json.Unmarshal(j2, &obj_map)
    if err != nil { log.Fatal("Json parsing error...") }
    servers := obj_map["servers"].([]interface{})
    stage := servers[1]
    stage_map := stage.(map[string]interface{})
    stage_map["url"] = stage_url
    servers[1]=stage_map
    obj_map["servers"]= servers
    return obj_map
}

func WriteYAML (filepath string, objmap map[string]interface{}) error {
    json, err := json.Marshal(objmap)
    if err != nil {
        fmt.Printf("err: %v\n", err)
        return err
        }
    yml, err := yaml.JSONToYAML(json)
    if err != nil {
        fmt.Printf("err: %v\n", err)
        return err
        }
    string_yml := strings.Replace(string(yml), "null", "", -1)
    //filepath := []string{folder_path,"config.yml"}
    fo, err := os.Create(filepath)
    if err != nil {
        fmt.Printf("err: %v\n", err)
        return err
        }
    _, err = io.Copy(fo, strings.NewReader(string_yml))
    if err != nil {
        fmt.Printf("err: %v\n", err)
        return err
        }
    return nil
}
func NowApi(folder_path string) error {
    cmd:= exec.Command("now")
    cmd.Dir = folder_path
     _, err := cmd.Output()
    if err != nil {
         return err
        }
    return nil 
}

func GetApiUrl(folder_path string) (string,error) {
    cmd:= exec.Command("now","ls","api")
    cmd.Dir = folder_path
    now_ls, err := cmd.Output()
    if err != nil {
         return "",err
        }
   string_now := string(now_ls)
   instances:= strings.Split(string_now,"api-")
   first_row := instances[1]
   parsing_row := strings.Split(first_row,".now.sh")
   hash := parsing_row[0]
   url:= fmt.Sprintf("https://api-%s.now.sh",hash)
   return url,nil
}

func NowAlias (folder_path string, source_url string, destination_url string) error {
    cmd:= exec.Command("now","alias",source_url,destination_url)
    cmd.Dir = folder_path
     _, err := cmd.Output()
    if err != nil {
         return err
        }
    return nil
}

func AcceptPull(branch string, test_path string) (string,string) {
	msg:= ""
	err_rm := DeleteFolder(test_path)
    if err_rm != nil {
        msg = "Error DeleteFolder: "+ fmt.Sprint(err_rm)
        fmt.Println(msg)
        return msg,"" 
    }
    err_clone := GitClone(test_path)
    if err_clone != nil {
        msg = "Error GitClone: "+ fmt.Sprint(err_clone)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg,"" 
    }
    err_hub := HubCheckout("master", test_path)
    if err_hub != nil {
        msg = "Error GitCheckout: "+ fmt.Sprint(err_hub)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg,"" 
    }
    err_merge := MergeBranch(branch, test_path)
    if err_merge != nil {
        msg = "Error MergeBranch: "+ fmt.Sprint(err_merge)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg,""
    }
    err_push := HubPush(test_path)
    if err_merge != nil {
        msg = "Error HubPush: "+ fmt.Sprint(err_push)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg,""
    }
    err_now:= NowDeploy(test_path)
    if err_now != nil {
        msg = "Error Now Deploy: "+ fmt.Sprint(err_push)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg,"" 
    }
    url,err_sub := GetSubdomain()
    if err_sub !=nil {
        msg = "Error Getting Subdomain: "+ fmt.Sprint(err_sub)+":no_entry_sign:"
        fmt.Println(msg)
        return msg,""
    }
    obj_map:=EditOpenApi(test_path,url)
    file_slice := []string{test_path,"/api/openapi.yaml"}
    filepath:= strings.Join(file_slice, "")
    err_yml:= WriteYAML (filepath, obj_map) 
    if err_yml !=nil {
        msg = "Error Saving openapi.yaml: "+ fmt.Sprint(err_yml)+":no_entry_sign:"
        fmt.Println(msg)
        return msg,""
    }
    err_commit:= GitCommit(test_path,"api/openapi.yaml", "New endpoint in stage version")
    if err_commit !=nil {
        msg = "Error GitCommit: "+ fmt.Sprint(err_commit)+":no_entry_sign:"
        fmt.Println(msg)
        return msg,""
    }
    err_push = HubPush(test_path)
    if err_merge != nil {
        msg = "Error HubPush: "+ fmt.Sprint(err_push)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg,""
    }
    path_slice := []string{test_path,"/api"}
    path:= strings.Join(path_slice, "")
    err_api:= NowApi(path)
    if err_api != nil {
        msg = "Error OpenAPi Now Deploy: "+ fmt.Sprint(err_api)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg,""
    }
    api_url,err_url := GetApiUrl(test_path)
    if err_url !=nil {
        msg = "Error Getting Api Url: "+ fmt.Sprint(err_url)+":no_entry_sign:"
        fmt.Println(msg)
        return msg,""
    }
    destination_url:= "https://api.geosure.tech/"
    err_alias := NowAlias(path,api_url,destination_url)
    if err_url !=nil {
        msg = "Error Now Alias: "+ fmt.Sprint(err_alias)+":no_entry_sign:"
        fmt.Println(msg)
        return msg,""
    }

    return msg,url
}
