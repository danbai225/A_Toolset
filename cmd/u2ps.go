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
	"flag"
	"github.com/spf13/cobra"
	"log"
	"p00q.cn/A_Toolset/u2ps/client"
	"p00q.cn/A_Toolset/u2ps/common"
	"p00q.cn/A_Toolset/u2ps/server"
)

// tcpCmd represents the serve command
var u2psCmd = &cobra.Command{
	Use:   "u2ps",
	Short: "内网穿透客户端",
	Long:  `内网穿透客户端 官网(使用前注册获取key):u2ps.com`, Example: "使用例子： A u2ps 123456789",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			s := args[0]

			log.Println("版本:", common.Versions)
			flag.Parse()
			if common.NodeMode {
				log.Println("节点模式运行...")
				common.Token = s
				server.Conn()
			} else {
				log.Println("客户端模式运行...")

				common.Key = s
				println(common.Key)
				client.Conn()
			}
		} else {
			println(cmd.UsageString())
		}
	},
}

func init() {
	rootCmd.AddCommand(u2psCmd)
	u2psCmd.Flags().BoolVarP(&common.NodeMode, "Node", "n", false, "是否为node模式")
	u2psCmd.Flags().IntVar(&common.MaxRi, "r", 10, "最多重连次数")
	u2psCmd.Flags().StringVar(&common.HostInfo, "h", "server.u2ps.com:2251", "连接地址")
}
