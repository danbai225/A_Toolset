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
	"github.com/antchfx/htmlquery"
	"github.com/leekchan/accounting"
	"github.com/spf13/cobra"
	"math"
	"math/big"
	"p00q.cn/A_Toolset/itself"
	"p00q.cn/A_Toolset/utils"
	"strconv"
	"strings"
	"time"
)

// erCmd represents the serve command
var erCmd = &cobra.Command{
	Use:   "er [转换代码] [目标代码] [数量]",
	Short: "汇率",
	Long: `汇率转换 
美元:USD 人民币:USD 欧元:EUR 日元:JPY 新台币:TWD 港币:HKD 英镑:GBP 韩元:KRW
货币代码:http://www.cnhuilv.com/currency/`, Example: "使用例子： A re CNY USD 100",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 2 {
			hl := getHl(args[0], args[1])
			m, err := strconv.ParseFloat(args[2], 64)
			if err == nil {
				ac := accounting.Accounting{Symbol: accounting.LocaleInfo[args[1]].ComSymbol, Precision: 2}
				mul := big.NewFloat(math.MaxFloat64).Mul(big.NewFloat(hl), big.NewFloat(m))
				println(ac.FormatMoneyBigFloat(mul))
			}
		} else {
			fmt.Println(cmd.UsageString())
		}
	},
}

func init() {
	rootCmd.AddCommand(erCmd)
}
func getHl(s1 string, s2 string) float64 {
	hlmap := make(map[string]string)
	hl := itself.Get("re-" + s1 + "-" + s2)
	json.Unmarshal([]byte(hl), &hlmap)
	int64, _ := strconv.ParseInt(hlmap["time"], 10, 64)
	if time.Now().Unix()-int64 < 100 {
		float, err := strconv.ParseFloat(hlmap["hl"], 64)
		utils.Check(err)
		return float
	}
	strhl := re(s1, s2)
	float, err := strconv.ParseFloat(strhl, 64)
	utils.Check(err)
	hlmap["hl"] = strhl
	hlmap["time"] = strconv.FormatInt(time.Now().Unix(), 10)
	data, err := json.Marshal(hlmap)
	itself.Put("re-"+s1+"-"+s2, string(data))
	return float

	return 0
}
func re(scur, tcur string) (r string) {
	if tcur != "" && scur != "" {
		rq, err := htmlquery.LoadURL("http://www.cnhuilv.com/" + strings.ToLower(scur) + "/" + strings.ToLower(tcur))
		if err == nil {
			node := htmlquery.FindOne(rq, "/html/body/div[4]/div/div[1]/div[3]/div[2]/span[1]")
			if node != nil {
				r = htmlquery.InnerText(node)
				return r
			}
		}
	}
	return ""
}
