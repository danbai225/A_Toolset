package main

import (
	"log"
	"p00q.cn/A_Toolset/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	cmd.Execute()
}
