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

// hashCmd represents the serve command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "go编译",
	Long:  `go编译`, Example: "使用例子： A build",
	Run: func(cmd *cobra.Command, args []string) {
		file := ""
		if len(args) > 0 {
			file = args[0]
		}
		Build(file)
	},
}
var oss string
var arch string

func init() {
	rootCmd.AddCommand(buildCmd)
	hashCmd.Flags().StringVarP(&oss, "os", "o", "linux", "目标系统")
	hashCmd.Flags().StringVarP(&arch, "arch", "a", "amd64", "目标架构")
}
func Build(file string) {
	switch oss {
	case "mac":
		oss = "darwin"
	case "win":
		oss = "windows"
	}
	os.Setenv("CGO_ENABLED", "0")
	os.Setenv("GOOS", oss)
	os.Setenv("GOARCH", arch)
	exec.Command("go", "build", file).Run()
}
