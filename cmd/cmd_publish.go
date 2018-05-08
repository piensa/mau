package cmd

import (

    "github.com/sbstjn/hanu"
    "errors"
    "strings"
    "fmt"
    "io/ioutil"
    "log"
    "github.com/ghodss/yaml"
    "encoding/json"
)

func init() {
    Register(
        "publish <hash:string>",
        "Promotes the latest staging server to production. Requires the instance hash as confirmation.",
        func(conv hanu.ConversationInterface) {
            test_path = "/home/prometheus/geosure"
            hash, _ := conv.Match(0)
            validate,err,stage_url := ValidateHash(hash)
            if validate == false {
                if err != nil {conv.Reply(err.Error())}
                conv.Reply("Mismatching the given hash with the staging hash. :no_entry_sign:")     
            } else {
                conv.Reply("Your request is in progress... Waiting... :stopwatch:")
                msg:= PublishUrl(hash, test_path,stage_url,fmt.Sprintf("%s.api.geosure.tech",hash))
                if msg != ""{
                    fmt.Println(msg)
                    conv.Reply(msg)
                }
                conv.Reply(fmt.Sprintf("Production environment is now running https://%s.api.geosure.tech :smile:",hash))

            }
           
        },
    )
}

func ValidateHash (hash string) (bool, error,string){
    url := "https://api.geosure.tech/openapi.yaml"
    openapi,err := GetOpenApi(url)
    if err !=nil {
        return false,errors.New(fmt.Sprintf("Error Getting OpenApi: %v",err)),""
    }
    servers,err_server:=GetServers(openapi)
    if err_server !=nil {
        return false,errors.New(fmt.Sprintf("Error: %v",err_server)),""
    }
    staging_url:= servers["staging"]
    remove_geosure:= strings.Replace(staging_url, "https://geosure-", "", -1)
    hash_stage:=strings.Replace(remove_geosure, ".now.sh", "", -1)
    if hash_stage != hash {
        return false, nil, staging_url
    }
    return true,nil,staging_url
}
func ProductionOpenApi(test_path string,production_url string) map[string]interface{} {
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
    production := servers[2]
    production_map := production.(map[string]interface{})
    production_map["url"] = production_url
    servers[2]=production_map
    obj_map["servers"]= servers
    return obj_map
}

func PublishUrl(hash string, test_path string,staging_url string,production_url string) string {
    msg:= ""
    err_alias := NowAlias(test_path+"/api",staging_url,production_url)
    if err_alias !=nil {
        msg = "Error Now Alias: "+ fmt.Sprint(err_alias)+":no_entry_sign:"
        fmt.Println(msg)
        return msg
    }
    err_rm := DeleteFolder(test_path)
    if err_rm != nil {
        msg = "Error DeleteFolder: "+ fmt.Sprint(err_rm)
        fmt.Println(msg)
        return msg 
    }
    err_clone := GitClone(test_path)
    if err_clone != nil {
        msg = "Error GitClone: "+ fmt.Sprint(err_clone)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg
    }
    obj_map:=ProductionOpenApi(test_path,"https://"+production_url)
    err_yml:= WriteYAML (test_path+"/api/openapi.yaml", obj_map) 
    if err_yml !=nil {
        msg = "Error Saving openapi.yaml: "+ fmt.Sprint(err_yml)+":no_entry_sign:"
        fmt.Println(msg)
        return msg
    }
    err_api:= NowApi(test_path+"/api")
    if err_api != nil {
        msg = "Error OpenAPi Now Deploy: "+ fmt.Sprint(err_api)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg
    }
    api_url,err_url := GetApiUrl(test_path)
    if err_url !=nil {
        msg = "Error Getting Api Url: "+ fmt.Sprint(err_url)+":no_entry_sign:"
        fmt.Println(msg)
        return msg
    }
    err_alias = NowAlias(test_path+"/api",api_url,"api.geosure.tech")
    if err_alias !=nil {
        msg = "Error Now Alias: "+ fmt.Sprint(err_alias)+":no_entry_sign:"
        fmt.Println(msg)
        return msg
    }
    err_commit:= GitCommit(test_path,"api/openapi.yaml", "New endpoint in production version")
    if err_commit !=nil {
        msg = "Error GitCommit: "+ fmt.Sprint(err_commit)+":no_entry_sign:"
        fmt.Println(msg)
        return msg
    }
    err_push:= HubPush(test_path)
    if err_push != nil {
        msg = "Error HubPush: "+ fmt.Sprint(err_push)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg
    }
    return ""
}
