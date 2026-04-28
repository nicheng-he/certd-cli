package cmd

import (
	"certd-cli/constant"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	Version   = "1.0.0"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  fmt.Sprintf("显示 %s 的详细版本信息", constant.Name),
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan, color.Bold)
		green := color.New(color.FgGreen)

		cyan.Printf("%s 命令行工具\n", constant.Name)
		green.Printf("版本: %s\n", Version)
		green.Printf("Git提交: %s\n", GitCommit)
		green.Printf("构建时间: %s\n", BuildTime)
	},
}
