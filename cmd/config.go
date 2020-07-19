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
	"p00q.cn/A_Toolset/itself"
	"p00q.cn/A_Toolset/utils"
)

var ini bool
var iniName string
// confCmd represents the serve command
var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "Cil配置管理",
	Long: `Cil配置管理`,Example: "使用例子： A conf key val",
	Run: func(cmd *cobra.Command, args []string) {
		if ini{
			utils.AddInitData(iniName)
			return
		}
		if len(args)>0{
			k := args[0]
			if len(args)>1{
				v := args[1]
				itself.Put(k,v)
				fmt.Printf("%s = %s\r\n", k, v)
			}else {
				fmt.Printf("%s = %s\r\n", k, itself.Get(k))
			}
		}else{
			fmt.Println(cmd.Long)
		}
	},
}

func init() {
	rootCmd.AddCommand(confCmd)
	confCmd.Flags().BoolVarP(&ini,"ini","i",false,"初始化0")
	confCmd.Flags().StringVarP(&iniName,"iniName","n","./A","初始化文件名")
}
