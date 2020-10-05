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
	"encoding/json"
	"github.com/axgle/mahonia"
	"github.com/gogf/gf/net/ghttp"
	"github.com/spf13/cobra"
)

type ipInfo struct {
	IP          string `json:"ip"`
	Pro         string `json:"pro"`
	ProCode     string `json:"proCode"`
	City        string `json:"city"`
	CityCode    string `json:"cityCode"`
	Region      string `json:"region"`
	RegionCode  string `json:"regionCode"`
	Addr        string `json:"addr"`
	RegionNames string `json:"regionNames"`
	Err         string `json:"err"`
}

// hashCmd represents the serve command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "ip查询",
	Long:  `ip查询`, Example: "使用例子： A ip",
	Run: func(cmd *cobra.Command, args []string) {
		ip := ""
		if len(args) > 0 {
			ip = args[0]
		}
		info := GetIpInfo(ip)
		println("IP:", info.IP)
		println("地址:", info.Addr)
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}
func GetIpInfo(ip string) ipInfo {
	info := ipInfo{}
	get, err := ghttp.Get("http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip)
	if err != nil {
		println(err)
	} else {
		json.Unmarshal([]byte(mahonia.NewDecoder("gbk").ConvertString(get.ReadAllString())), &info)
	}
	return info
}
