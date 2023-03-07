package cmd

import (
	"bufio"
	"fmt"
	"github.com/solywsh/chatgpt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

var iai bool

// base64 represents the serve command
var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "ai chat",
	Long:  `ai助理基于chatgpt`, Example: "使用例子： A ai 你好请你帮我计算下PI的第100位",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			question := args[0]
			key := os.Getenv("ChatGPTKey")
			if key == "" {
				println("请设置ChatGPTKey环境变量")
				return
			}
			chat := chatgpt.New(key, "danbai", 10*time.Second)
			defer chat.Close()
			if iai {
				input := bufio.NewScanner(os.Stdin)
				fmt.Println("开始会话输入.exit 退出")
				for {
					input.Scan()
					if strings.Compare(strings.TrimSpace(input.Text()), "") == 0 {
						continue
					}
					if strings.Compare(strings.TrimSpace(input.Text()), ".exit") == 0 {
						os.Exit(0)
					}
					question = input.Text()
					answer, err := chat.Chat(question)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Printf("A: %s\n", answer)
				}
			} else {
				fmt.Printf("Q: %s\n", question)
				answer, err := chat.Chat(question)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("A: %s\n", answer)
			}
		} else {
			println(cmd.UsageString())
		}
	},
}

func init() {
	rootCmd.AddCommand(aiCmd)
	aiCmd.Flags().BoolVarP(&iai, "interaction", "i", false, "交互")
}
