package user

import (
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "用户管理",
	Long:  "用户管理子命令，支持添加、列出、删除等操作",
}

func init() {

}
