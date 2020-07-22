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
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/crypto/gsha1"
	"github.com/spf13/cobra"
	"io/ioutil"
	"p00q.cn/A_Toolset/utils"
)

var hashFile string
var sha string

// hashCmd represents the serve command
var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "hash 计算",
	Long:  `hash 计算`, Example: "使用例子： A hash hi,hash",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 || hashFile != "" {
			var err error
			var str string
			switch sha {
			case "md5":
				if hashFile != "" {
					str, err = gmd5.EncryptFile(hashFile)
				} else {
					str, err = gmd5.EncryptString(args[0])
				}
				break
			case "sha1":
				if hashFile != "" {
					str, err = gsha1.EncryptFile(hashFile)
				} else {
					str = gsha1.Encrypt(args[0])
				}
				break
			case "sha256":
				hash := sha256.New()
				var bs []byte
				if hashFile != "" {
					bs, _ = ioutil.ReadFile(hashFile)
				} else {
					bs = []byte(args[0])
				}
				hash.Write(bs)
				fmt.Printf("%x\n", hash.Sum(nil))
				return
			case "sha512":
				hash := sha512.New()
				var bs []byte
				if hashFile != "" {
					bs, _ = ioutil.ReadFile(hashFile)
				} else {
					bs = []byte(args[0])
				}
				hash.Write(bs)
				fmt.Printf("%x\n", hash.Sum(nil))
				return
			}
			utils.Check(err)
			println(str)
		} else {
			println(cmd.UsageString())
		}
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.Flags().StringVarP(&hashFile, "file", "f", "", "文件路径")
	hashCmd.Flags().StringVarP(&sha, "sha", "s", "md5", "算法(md5,sha1,sha256,sha512)")
}
