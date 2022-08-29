package cmd

import (
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
)

var (
	webPort     int
	webRootPath string
)

// web represents the serve command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "web服务",
	Long:  `web服务`, Example: "使用例子： A web",
	Run: func(cmd *cobra.Command, args []string) {
		s := g.Server()
		s.SetIndexFolder(true)
		s.SetServerRoot(webRootPath)
		s.SetPort(webPort)
		s.Run()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().IntVarP(&webPort, "port", "p", 80, "web端口")
	webCmd.Flags().StringVarP(&webRootPath, "root", "r", ".", "设置静态文件服务的目录路径（默认当前路径）")
}
