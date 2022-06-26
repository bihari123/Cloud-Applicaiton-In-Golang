// Package main is the entry point to the server. It reads configuration, sets up logging and error handling,
// handles signals from the OS, and starts and stops the server.
package main

import (
	"fmt"

	"github.com/bihari123/cloud-application-in-golang/helper"
	"github.com/bihari123/cloud-application-in-golang/loghelper"
)

func main() {
	myConfig := helper.GetConfig()
	fmt.Printf("\n\n%+v\n\n", myConfig)
	fmt.Println("LogFolder --", myConfig.AppConfig.LogFolder)

	loghelper.ConfigLogurus(myConfig)
	
	fmt.Println("🤓")
}
