package main

import (
	"os"

	"github.com/cibeiwanjia/cobra-viper/cmd"
)

func main() {
	// 初始化并执行 Cobra 根命令
	rootCmd := cmd.InitRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
