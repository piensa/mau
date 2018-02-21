package cmd

import (
    "github.com/sbstjn/hanu"
    "fmt"
    "net/http"
    "io/ioutil"
    "github.com/ghodss/yaml"
    "encoding/json"
    "errors"
)

func init() {
    Register(
        "servers",
        "List development, staging and production servers.",
        func(conv hanu.ConversationInterface) {
            url := "https://api.geosure.tech/openapi.yaml"
            msg, string_err := HandleServers(url)
            if string_err != "" {
                fmt.Println(string_err)
                conv.Reply(string_err)
            }
            if msg != "" {
                conv.Reply(msg)
            }             
        },
    )
}

func HandleServers (url string) (string,string){
    openapi,err := GetOpenApi(url)
    if err !=nil {
        
    }
    servers,err_server:=GetServers(openapi)
    if err_server !=nil {
        return "",fmt.Sprintf("Error: %v",err_server)
    }
    return FormatMessage(servers),""

}
func GetOpenApi(url string) (string,error) {
    resp, err := http.Get(url)
    if err != nil {
        return "",err
    }
    defer resp.Body.Close()
    body, err_get := ioutil.ReadAll(resp.Body)
    if err_get != nil {
        return "", err_get
    }
    string_body := string(body)
    return string_body,nil
}
func GetServers (body string) (map[string]string,error) {
    output := make(map[string]string)
    j2, err := yaml.YAMLToJSON([]byte(body))
    var msg string 
    if err != nil {
        msg = fmt.Sprintf("Yaml parsing error in OpenApi.yaml: %v",err)
        return output,errors.New(msg)
    }
    var obj_map map[string]interface{} 
    err = json.Unmarshal(j2, &obj_map)
    if err != nil { 
        msg = fmt.Sprintf("Json parsing error in OpenApi: %v",err)
        return output,errors.New(msg)
    }
    servers := obj_map["servers"].([]interface{})
    for _, item := range servers {
        temp:= item.(map[string]interface{})
        server_type:=temp["description"].(string)
        server_url := temp["url"].(string)
        output[server_type] = server_url
    }
    return output,nil
}

func FormatMessage(servers map[string]string) string {
    var message string 
    message = "The API servers are currently:\n"
    for key := range servers {
        message = message+fmt.Sprintf("`%s`: %s\n",key,servers[key])
    }
    return message
}
