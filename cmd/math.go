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
	"github.com/dengsgo/math-engine/engine"
	"github.com/spf13/cobra"
)

// mathCmd represents the serve command
var mathCmd = &cobra.Command{
	Use:   "=",
	Short: "计算数学公式",
	Long: `该命令用于简单数学公式计算
该命令能够处理的表达式样例：
	1+127-21+(3-4)*6/2.5
	(88+(1+8)*6)/2+99
	123_345_456 * 1.5 - 2 ^ 4
	-4 * 6 + 2e2 - 1.6e-3
	sin(pi/2)+cos(45-45*1)+tan(pi/4)
	99+abs(-1)-ceil(88.8)+floor(88.8)
	max(min(2^3, 3^2), 10*1.5-7)
	double(6) + 3
详情：https://github.com/dengsgo/math-engine
`, Example: "使用例子： A = 1+1/2 \n输出：1+1/2 = 1.5",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			s := args[0]
			r, err := engine.ParseAndExec(s)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s = %v\r\n", s, r)
		} else {
			fmt.Println(cmd.Long)
		}

	},
}

func init() {
	rootCmd.AddCommand(mathCmd)
}
