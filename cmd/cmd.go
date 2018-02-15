package cmd

import (
	"time"
	"github.com/sbstjn/hanu"
)

var commandList []hanu.CommandInterface

// Version stores the chatbot version. Will be updated by `main.go`
var Version string

// Start contains the Time when the process started
var Start time.Time

// String contains prometheus config path 
var config_Prometheus string 

// String contains test folder for PR test
var test_path string 

// String access token 
var GitToken string 
GitToken = "998bcfb5d565dd23bd04276093e3f6b71c06bf5d"

// Register adds a new command to commandList
func Register(cmd string, description string, handler hanu.Handler) {
	commandList = append(commandList, hanu.NewCommand(cmd, description, handler))
}

// List returns commandList
func List() []hanu.CommandInterface {
	return commandList
}
