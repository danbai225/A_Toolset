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
	"github.com/gogf/gf/net/gudp"
	"github.com/spf13/cobra"
	"os"
	"p00q.cn/A_Toolset/utils"
	"strings"
)

var udpServer bool
var udpPort string
var udpAddress string
var udpFormat string

// udpCmd represents the serve command
var udpCmd = &cobra.Command{
	Use:   "udp",
	Short: "udp连接测试",
	Long:  `udp连接测试`, Example: "使用例子： A udp 127.0.0.1:225",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			s := args[0]
			split := strings.Split(s, ":")
			if len(split) > 1 {
				udpAddress = args[0]
			} else {
				udpAddress = s
				udpAddress += ":" + udpPort
			}
			if udpServer {
				gudp.NewServer(udpAddress, func(conn *gudp.Conn) {
					defer conn.Close()
					for {
						bytes := make([]byte, 1024)
						len, add, err2 := conn.ReadFrom(bytes)
						if err2 == nil {
							fmt.Printf("form:%s data:%"+udpFormat+"\n", add.String(), bytes[:len])
						}
					}
				}).Run()
			} else {
				conn, err := gudp.NewConn(udpAddress)

				if err != nil {
					println("未能连接服务端")
					os.Exit(0)
				} else {
					println("连接上服务端")
					defer conn.Close()
				}
				go SendUdpMsg(conn)
				for {
					bytes := make([]byte, 1024)
					len, add, err2 := conn.ReadFrom(bytes)
					if err2 == nil {
						fmt.Printf("form:%s data:%"+udpFormat, add.String(), bytes[:len])
					} else {
						println("连接断开")
						os.Exit(0)
					}
				}
			}
		} else {
			println(cmd.UsageString())
		}
	},
}

// 向服务器端发消息
func SendUdpMsg(conn *gudp.Conn) {
	for {
		var in string
		fmt.Scanf("%"+udpFormat, &in)
		err := conn.Send([]byte(in))
		utils.Check(err)
	}
}
func init() {
	rootCmd.AddCommand(udpCmd)
	udpCmd.Flags().StringVarP(&udpPort, "port", "p", "2251", "udp端口")
	udpCmd.Flags().BoolVarP(&udpServer, "server", "s", false, "是否为服务端模式")
	udpCmd.Flags().StringVarP(&udpFormat, "format", "f", "s", "格式化数据(string) s 字符串 b二进制 x十六进制")
}
