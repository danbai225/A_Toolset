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
