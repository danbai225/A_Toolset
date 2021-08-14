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
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
)

//https://restapi.amap.com/v5/ip?parameters?&key=9fe97f8e6302c01f087b12f0a5cdfb1c&type=4&ip=171.217.139.247
const gdKey = "9fe97f8e6302c01f087b12f0a5cdfb1c"

// tqCmd represents the serve command
var tqCmd = &cobra.Command{
	Use:   "tq",
	Short: "天气",
	Long:  `获取当前天气`,
	Run: func(cmd *cobra.Command, args []string) {
		info := GetIpInfo("")
		get, err := http.Get(fmt.Sprintf("https://restapi.amap.com/v5/ip?parameters?&key=9fe97f8e6302c01f087b12f0a5cdfb1c&type=4&ip=%s", info.IP))
		if err == nil {
			all, err := ioutil.ReadAll(get.Body)
			if err != nil {
				log.Println(err.Error())
				return
			}
			location := gjson.GetBytes(all, "location").String()
			url := fmt.Sprintf("https://api.caiyunapp.com/v2.5/ujp0HddE4bY2SwRc/%s/realtime.json", location)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err.Error())
				return
			}
			readAll, _ := ioutil.ReadAll(resp.Body)
			log.Println("温度", gjson.GetBytes(readAll, "result.realtime.temperature").String())
			log.Println("体感温度", gjson.GetBytes(readAll, "result.realtime.apparent_temperature").String())
		} else {
			log.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(tqCmd)
}
