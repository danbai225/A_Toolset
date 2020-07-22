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
	"encoding/base64"
	"github.com/spf13/cobra"
)

var dBase64 bool

// base64 represents the serve command
var baseCmd = &cobra.Command{
	Use:   "base",
	Short: "base64 加解密",
	Long:  `base64 是一种常见的编码 base64用于快速加解密`, Example: "使用例子： A base hi,base64",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			s := args[0]
			if dBase64 {
				sDec, _ := base64.StdEncoding.DecodeString(s)
				println(string(sDec))
			} else {
				b := []byte(s)
				sEnc := base64.StdEncoding.EncodeToString(b)
				println(sEnc)
			}
		} else {
			println(cmd.UsageString())
		}

	},
}

func init() {
	rootCmd.AddCommand(baseCmd)
	baseCmd.Flags().BoolVarP(&dBase64, "decode", "d", false, "是否解码(默认是编码)")
}
