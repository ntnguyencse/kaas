// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package emcoctl

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/buildkite/interpolate"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	FlagDelete string = "delete"
	FlagApply  string = "apply"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "emco",
	Short: "Emcoctl - CLI for EMCO",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { fmt.Println("emcoctl <command> -f file") },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func ExecWithError() error {
	return rootCmd.Execute()
}
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.emco.yaml)")

}

// Set Config File Path
func SetConfigFilePath(path string) {
	cfgFile = path
}
func SaveValueFile(fileName string, folder string, valueString *string) (string, error) {
	// If folder is empty, save file to current folder
	if folder == "" {
		folder, _ = os.Getwd()
	}
	// Remove file before
	fmt.Println("Remove file: ", folder+fileName)
	e := os.Remove(folder + fileName)
	if e != nil {
		fmt.Println(e, "Error when remove folder")
	}
	if _, err := os.Stat(folder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(folder, os.ModePerm)
		if err != nil {
			fmt.Println(err, "Error when create folder")
		}
	}
	filePath := folder + fileName
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		fmt.Println(err, "Error when open file")
		return "error Open", err
	}
	_, err = file.WriteString(*valueString)
	if err != nil {
		fmt.Println(err, "Error when write file")
		return "error Write", err
	}
	return filePath, nil
}
func CleanUp(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println(err, "Error clean up file: "+filePath)
		return err
	}
	fmt.Println("Cleanup file: " + filePath)
	return nil
}
func InterpolateValueFile(templateString *string, values map[string]string) (string, error) {

	var resultString string
	mapEnv := interpolate.NewMapEnv(values)
	resultString, err := interpolate.Interpolate(mapEnv, *templateString) //"Buildkite... ${HELLO_WORLD}}")
	if err != nil {
		fmt.Println("Error interpolate", err)
		return "", err
	}
	fmt.Println(resultString)
	return resultString, nil
}

// SetArg for root command
func SetArgs(args []string) {
	rootCmd.SetArgs(args)
}
func SetDebugFlags() {
	rootCmd.DebugFlags()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".emco" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".emco")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		err := viper.Unmarshal(&Configurations)
		if err != nil {
			fmt.Printf("Unable to decode into struct, %v", err)
		}
	} else {
		fmt.Println("Warning: No Configuration File found. Using defaults")
		SetDefaultConfiguration()
	}
}
