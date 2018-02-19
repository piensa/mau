package cmd

import "github.com/sbstjn/hanu"

func init() {
	Register(
		"accept <branch:string>",
		"Merges open Pull Request for specified branch onto master and redeploys master on the staging server.",
		func(conv hanu.ConversationInterface) {
			branch, _ := conv.Match(0)
			conv.Reply("Branch `" + branch + "` has been merged intop master and the staging server is now: " + "piknedixce.api.geosure.tech")
			test_path = "/home/prometheus/test"
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
    cmd := exec.Command("hub","merge",branch)
    cmd.Dir = test_path
    _, err := cmd.Output()
    if err != nil {
         return err
        }
    return nil 
}

func AcceptPull(branch string, test_path string) string {
	msg:= ""
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
    err_hub := HubCheckout("master", test_path)
    if err_hub != nil {
        msg = "Error GitCheckout: "+ fmt.Sprint(err_hub)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg 
    }
    err_merge := MergeBranch(branch, test_path)
    if err_merge != nil {
        msg = "Error GitCheckout: "+ fmt.Sprint(err_merge)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg 
    }
    err_push := HubPush()
    if err_merge != nil {
        msg = "Error HubPush: "+ fmt.Sprint(err_push)+ ":no_entry_sign:"
        fmt.Println(msg)
        return msg 
    }
    return msg
}
