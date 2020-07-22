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
	"github.com/gogf/gf/net/gtcp"
	"github.com/spf13/cobra"
	"os"
	"p00q.cn/A_Toolset/utils"
	"strings"
)

var tcpServer bool
var tcpPort string
var tcpAddress string
var tcpFormat string

// tcpCmd represents the serve command
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "tcp连接测试",
	Long:  `tcp连接测试`, Example: "使用例子： A tcp 127.0.0.1:225",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			s := args[0]
			split := strings.Split(s, ":")
			if len(split) > 1 {
				tcpAddress = args[0]
			} else {
				tcpAddress = s
				tcpAddress += ":" + tcpPort
			}
			if tcpServer {
				gtcp.NewServer(tcpAddress, func(conn *gtcp.Conn) {
					defer conn.Close()
					fmt.Printf("有新连接加入:%s!\n", conn.RemoteAddr().String())
					for {
						bytes := make([]byte, 1024)
						len, err2 := conn.Read(bytes)
						if err2 == nil {
							addr := conn.RemoteAddr
							fmt.Printf("form:%s data:%"+tcpFormat+"\n", addr().String(), bytes[:len])
						}
					}
				}).Run()
			} else {
				conn, err := gtcp.NewConn(tcpAddress)
				if err != nil {
					println("未能连接上服务端")
					os.Exit(0)
				} else {
					println("连接上服务端")
					defer conn.Close()
				}
				go SendTcpMsg(conn)
				for {
					bytes := make([]byte, 1024)
					len, err2 := conn.Read(bytes)
					if err2 == nil {
						addr := conn.RemoteAddr
						fmt.Printf("form:%s data:%"+tcpFormat, addr().String(), bytes[:len])
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
func SendTcpMsg(conn *gtcp.Conn) {
	for {
		var in string
		fmt.Scanf("%"+tcpFormat, &in)
		err := conn.Send([]byte(in))
		utils.Check(err)
	}
}
func init() {
	rootCmd.AddCommand(tcpCmd)
	tcpCmd.Flags().StringVarP(&tcpPort, "port", "p", "2251", "tcp端口")
	tcpCmd.Flags().BoolVarP(&tcpServer, "server", "s", false, "是否为服务端模式")
	tcpCmd.Flags().StringVarP(&tcpFormat, "format", "f", "s", "格式化数据(string) s 字符串 b二进制 x十六进制")
}
