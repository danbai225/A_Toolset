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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
	"io/ioutil"
	"os"
	"p00q.cn/A_Toolset/utils"
	"strings"
)

var dirPath = "./"
var rNum int
var size bool

// treeCmd represents the serve command
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "目录树生成",
	Long:  `生成指定目录的目录树`, Example: "A tree -l 5 -s",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			dirPath = args[0]
			l := dirPath[len(dirPath)-1:]
			if l != string(os.PathSeparator) {
				dirPath += string(os.PathSeparator)
			}
		}
		if utils.PathExists(dirPath) {
			tree := treeprint.New()
			add(dirPath, tree)
			fmt.Println(tree.String())
		}
	},
}

func add(path string, tree treeprint.Tree) {
	if strings.Count(strings.ReplaceAll(path, dirPath, ""), utils.PathSeparator()) < rNum {
		files, err := ioutil.ReadDir(path)
		utils.Check(err)
		for _, f := range files {
			if f.IsDir() {
				var branch treeprint.Tree
				if size {
					branch = tree.AddMetaBranch(utils.FormatFileSize(utils.DirSizeB(utils.AddPath(path, f.Name()))), f.Name())
				} else {
					branch = tree.AddBranch(f.Name())
				}
				add(utils.AddPath(path, f.Name()), branch)
			} else {
				if size {
					tree.AddMetaNode(utils.FormatFileSize(f.Size()), f.Name())
				} else {
					tree.AddNode(f.Name())
				}
			}
		}
	}
}
func init() {
	rootCmd.AddCommand(treeCmd)
	treeCmd.Flags().IntVarP(&rNum, "level", "l", 4, "递归的目录层级(默认4)")
	treeCmd.Flags().BoolVarP(&size, "size", "s", false, "显示大小")
}
