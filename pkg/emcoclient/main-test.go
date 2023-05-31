package emcoctl

import (
	"fmt"
)

const (
	config = "/home/ubuntu/emco-client/.emco.yaml"
)

func maintest() {
	fmt.Println("Beginning of Cmd")
	fmt.Println("Execute the cmd")
	// cfgFile = config
	SetConfigFilePath(config)
	args := []string{"apply", "-f", "/home/ubuntu/emco-client/test/register-cluster.yml"}
	SetArgs(args)
	// rootCmd.SetArgs(args)
	// rootCmd.DebugFlags()
	SetDebugFlags()
	Execute()
	resqq := GetResponseOutput()
	fmt.Println(resqq.StatusCode())
	fmt.Println(string(resqq.Body()))
}
