package cmd

import(
	"github.com/sbstjn/hanu"
	"os/exec"
    "strings"

) 

func init() {
	Register(
		"check <PR:number>",
		"Do the unitests and coverage from a given Pull Request",
		func(conv hanu.ConversationInterface) {
			test_path = "/home/prometheus/test"
			pull_request, _ := conv.Match(0)
			coverage,time,fail, msg_err = CheckPR(pull,test_path)
			if msg_err != "" {
				conv.Reply(msg_err)
			} else {
				template_msg :=`test:coverage=%s
				test:passed=%t
				test:time=%s seconds`
				msg = fmt.Sprintf(template_msg,coverage,fail,time)

			}
			conv.Reply("Branch `" + branch + "` has been merged intop master and the staging server is now: " + "piknedixce.api.geosure.tech")
		},
	)
}

func GitClone(test_path string){
	_, err := exec.Command("git","clone","https://github.com/geosure/geosure.git",test_path).Output()
   	if err != nil {
    	return err
    }
    return nil 
}

func GitCheckout(pull string, test_path string ) error {
	cmd := exec.Command("hub","pr","checkout",pull)
	cmd.Dir = test_path
	_, err = cmd.Output()
	if err != nil {
         //fmt.Println("error checkout: "+ fmt.Sprint(err))
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
    test, err := cmd_test.Output()
    fail := false 
    if err != nil {
        //fmt.Println("error make test: "+ fmt.Sprint(err))
        return coverage, time,true,err
    }
    string_test := string(test)
    if strings.Index(string_test,"FAIL") > 0 {
    	fail = true
    }
    split_coverage := strings.Split(string_test, "coverage: ")
    split_percent := strings.Split(split_coverage[1],"%")
    split_dir := strings.Split (split_percent[1],split_folder)
    split_s:= strings.Split(split_dir[1],"s")
    coverage = split_percent[0]
    time = strings.TrimSpace(split_s[0])
    return coverage, time, fail,nil 
}

func DeleteFolder (test_path string) error {
	cmd_rm := exec.Command("rm","-rf",test_path)
	_, err := cmd_rm.Output()
	if err != nil {
		return err
	}
	return nil 
}

func CheckPR (PR string, test_path string) (string,string,string) {
	err_clone := GitClone(test_path)
	if err != nil {
		msg := "Error GitClone: "+ fmt.Sprint(err_clone)
		fmt.Println(msg)
		return "","",msg 
	}
	err_checkout := GitCheckout(PR, test_path)
	if err != nil {
		msg := "Error GitCheckout: "+ fmt.Sprint(err_checkout)
		fmt.Println(msg)
		return "","",msg 
	}
	coverage, time , err_make := MakeTest(test_path)
	if err != nil {
		msg := "Error MakeTest: "+ fmt.Sprint(err_make)
		fmt.Println(msg)
		return "","",msg 
	}
	err_rm := DeleteFolder(test_path)
	if err != nil {
		msg := "Error DeleteFolder: "+ fmt.Sprint(err_rm)
		fmt.Println(msg)
		return "","",msg 
	}
	return coverage,time, ""
}

