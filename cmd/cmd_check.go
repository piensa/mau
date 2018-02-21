package cmd

import(
    "github.com/sbstjn/hanu"
    "os/exec"
    "strings"
    "fmt"
    "net/http"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "log"
    "errors"
) 

func init() {
    Register(
        "check <PR:string>",
        "Do the unitests and coverage from a given Pull Request",
        func(conv hanu.ConversationInterface) {
            ConfigFile, err := ioutil.ReadFile("./config.json")
            if err != nil {
                 log.Fatal(err)
            }
            test_path = "/home/prometheus/geosure"
            var config ConfigStruct
            var url string
            var err_sub error 
            json.Unmarshal(ConfigFile, &config)
            GitToken = config.GithubToken
            pull_request, _ := conv.Match(0)
            coverage,time,pass, msg_err := CheckPR(pull_request,test_path)
            if msg_err != "" {
                conv.Reply(msg_err)
                err_co:= GitComment(pull_request,msg_err)
                if err_co != nil {
                    msg_co := "Error GitComment: "+ fmt.Sprint(err_co)
                    conv.Reply(msg_co)
                    fmt.Println(msg_co)
                }

            } else {
                url = "none"
                if pass == true {
                    conv.Reply("Your request is progress... :stopwatch:")
                    err_deploy:= NowDeploy(test_path)
                    if err_deploy != nil {
                        msg_deploy := "Error Now Deploy: "+ fmt.Sprint(err_deploy)+":no_entry_sign:"
                        fmt.Println(msg_deploy)
                        conv.Reply(msg_deploy) 
                    }
                    url,err_sub = GetSubdomain()
                    if err_sub !=nil {
                        msg_sub := "Error Getting Subdomain: "+ fmt.Sprint(err_sub)+":no_entry_sign:"
                        fmt.Println(msg_sub)
                        conv.Reply(msg_sub) 
                    }
                } 
                template_msg :="deployment:url=%s\\ntest:coverage=%s\\ntest:passed=%t\\ntest:time=%s seconds\\n"
                msg:= fmt.Sprintf(template_msg,url,coverage,pass,time)
                conv.Reply(strings.Replace(msg,"\\n","\n",-1))
                git_error:= GitComment(pull_request,msg)
                if git_error != nil {
                    string_error := "Error GitComment: "+ fmt.Sprint(git_error)
                    conv.Reply(string_error)
                    fmt.Println(string_error)
                }
            }
        },
    )
}

func GitClone(test_path string) error {
    _, err := exec.Command("git","clone","git@github.com:geosure/geosure.git",test_path).Output()
    if err != nil {
        return err
    }
        return nil 
    }

func GitCheckout(pull string, test_path string ) error {
    cmd := exec.Command("hub","pr","checkout",pull)
    cmd.Dir = test_path
    _, err := cmd.Output()
    if err != nil {
         return err
        }
    return nil 
}

func MakeTest (test_path string) (string, string,bool,error) {
    coverage := ""
    time:= ""
    split_folder:= "_"+test_path
    cmd_test := exec.Command("make","test")
    cmd_test.Dir = test_path
    test, _ := cmd_test.Output()
    pass := true 
    string_test := string(test)
    if strings.Index(string_test,"coverage:") < 0 && strings.Index(string_test,"failed") > 0 {
        pass = false
        error_build:= errors.New("Golang Error: Build failed. Check your code!")
        return coverage,time,pass,error_build
    }
    if strings.Index(string_test,"FAIL") > 0 {
        pass = false
    }
    split_coverage := strings.Split(string_test, "coverage: ")
    split_percent := strings.Split(split_coverage[1],"%")
    split_dir := strings.Split (split_percent[1],split_folder)
    split_s:= strings.Split(split_dir[1],"s")
    coverage = split_percent[0]
    time = strings.TrimSpace(split_s[0])
    return coverage, time, pass,nil 
}

func DeleteFolder (test_path string) error {
    cmd_rm := exec.Command("rm","-rf",test_path)
    _, err := cmd_rm.Output()
    if err != nil {
        return err
    }
    return nil 
}

func GitComment (PR string, msg string) error {
    raw := fmt.Sprintf(`{"body": "%s"}`,msg)
    url := "https://api.github.com/repos/geosure/geosure/issues/%s/comments?access_token=%s"
    githuburl:=fmt.Sprintf(url,PR,GitToken)
    _ , err := http.Post(githuburl, "application/json", bytes.NewBuffer([]byte(raw)))
    if err != nil {return err}
    return nil 
}
func CheckPR (PR string, test_path string) (string,string,bool,string) {
    err_rm := DeleteFolder(test_path)
    if err_rm != nil {
        msg := "Error DeleteFolder: "+ fmt.Sprint(err_rm)
        fmt.Println(msg)
        return "","",false,msg 
    }
    err_clone := GitClone(test_path)
    if err_clone != nil {
        msg := "Error GitClone: "+ fmt.Sprint(err_clone)+":no_entry_sign:"
        fmt.Println(msg)
        return "","",false,msg 
    }
    err_checkout := GitCheckout(PR, test_path)
    if err_checkout != nil {
        msg := "Invalid Pull Request Number :no_entry_sign:"
        fmt.Println(msg)
        return "","",false,msg 
    }
    coverage, time , pass, err_make := MakeTest(test_path)
    if err_make != nil {
        msg := "Unitest Error: "+ fmt.Sprint(err_make)+":no_entry_sign:"
        fmt.Println(msg)
        return "","",false,msg 
    }
    return coverage,time,pass, ""
}

func NowDeploy(test_path string) error {
    cmd_deploy := exec.Command("make","deploy")
    cmd_deploy.Dir = test_path
    _, err := cmd_deploy.Output()
    if err != nil {
        return err
    }
    return nil 
}

func GetSubdomain() (string,error){
    results, err:= exec.Command("now","ls","geosure").Output()
    if err != nil {
                return "",err
        }
    string_instances := string(results)
    instances:= strings.Split(string_instances,"geosure-")
    first_row := instances[1]
    parsing_row := strings.Split(first_row,".now.sh")
    hash := parsing_row[0]
    url:= fmt.Sprintf("https://geosure-%s.now.sh",hash)
    return url, nil 
}
