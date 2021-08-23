/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// buildCmd represents the serve command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "go编译",
	Long:  `go编译`, Example: "使用例子： A build",
	Run: func(cmd *cobra.Command, args []string) {
		file := ""
		outFile := ""
		if len(args) > 0 {
			file = args[0]
		}
		if len(args) > 1 {
			outFile = args[1]
		}
		Build(file, outFile)
	},
}
var oss string
var arch string
var release bool

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&oss, "os", "o", "", "目标系统")
	buildCmd.Flags().StringVarP(&arch, "arch", "a", "", "目标架构")
	buildCmd.Flags().BoolVarP(&release, "release", "r", false, "发行")
}
func Build(file, outFile string) {
	switch oss {
	case "mac":
		oss = "darwin"
	case "win":
		oss = "windows"
	}
	if oss != "" {
		os.Setenv("CGO_ENABLED", "0")
		os.Setenv("GOOS", oss)
	}
	if arch != "" {
		os.Setenv("GOARCH", arch)
	}
	var command *exec.Cmd
	if release {
		command = exec.Command("go", "build", `-ldflags="-s -w"`, "-o", outFile, file)
	} else {
		command = exec.Command("go", "build", "-o", outFile, file)
	}
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Run()
}
