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
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/liushuochen/gotable"
	"github.com/spf13/cobra"
	"time"
)

var tqHours int64

type tq struct {
	Results []struct {
		Location struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			Country        string `json:"country"`
			Path           string `json:"path"`
			Timezone       string `json:"timezone"`
			TimezoneOffset string `json:"timezone_offset"`
		} `json:"location"`
		Hourly []struct {
			Time          time.Time `json:"time"`
			Text          string    `json:"text"`
			Code          string    `json:"code"`
			Temperature   string    `json:"temperature"`
			Humidity      string    `json:"humidity"`
			WindDirection string    `json:"wind_direction"`
			WindSpeed     string    `json:"wind_speed"`
		} `json:"hourly"`
	} `json:"results"`
}

// hashCmd represents the serve command
var tqCmd = &cobra.Command{
	Use:   "tq",
	Short: "天气查询",
	Long:  `天气查询`, Example: "使用例子： A tq ",
	Run: func(cmd *cobra.Command, args []string) {
		info := tq{}
		city := ""
		if len(args) > 0 {
			city = args[0]
		} else {
			city = GetIpInfo("").City
		}
		get, err := ghttp.Get("https://api.seniverse.com/v3/weather/hourly.json?key=SaKrHjSjVx09-euRZ&location=" + city + "&language=zh-Hans&unit=c&start=0&hours=" + gconv.String(tqHours))
		if err != nil {
			println(err)
		} else {
			json.Unmarshal(get.ReadAll(), &info)
			if len(info.Results) < 1 {
				println("获取失败")
				return
			}
			println("位置:" + info.Results[0].Location.Name + "未来12小时天气")
			formHead := []string{"时间", "天气", "温度", "湿度", "风向", "风速"}
			demo, err := gotable.Create(formHead...)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			values := []map[string]string{}
			for _, h := range info.Results[0].Hourly {
				values = append(values, map[string]string{
					"时间": gconv.String(h.Time.Day()) + "日" + gconv.String(h.Time.Hour()) + "时",
					"天气": h.Text,
					"温度": h.Temperature + "℃",
					"湿度": h.Humidity + "%",
					"风向": h.WindDirection,
					"风速": h.WindSpeed + "km/h",
				})
			}
			for _, value := range values {
				err := demo.AddRow(value)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
			demo.PrintTable()
		}
	},
}

func init() {
	rootCmd.AddCommand(tqCmd)
	tqCmd.Flags().Int64VarP(&tqHours, "hours", "s", 12, "指定查看未来几小时")
}
