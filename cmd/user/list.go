package user

import (
	"certd-cli/constant"
	"certd-cli/database"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出用户",
	Long:  "从数据库中查询并列出用户列表",
	Example: fmt.Sprintf(`  # 列出所有用户（默认前20条）
  %s user list

  # 列出前10个用户
  %s user list -l 10

  # 分页：第2页，每页10条
  %s user list -l 10 -o 10

  # 详细模式查看执行信息
  %s user list -v`, constant.Name, constant.Name, constant.Name, constant.Name),
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")

		container := database.GetContainer()

		if !container.IsInitialized() {
			return fmt.Errorf("数据库未初始化，请检查配置文件或环境变量")
		}

		db, err := container.GetDatabase()
		if err != nil {
			return fmt.Errorf("获取数据库连接失败: %w", err)
		}

		if verbose {
			fmt.Println("正在查询用户列表...")
		}

		users, err := queryUsers(db, limit, offset)
		if err != nil {
			return err
		}

		if len(users) == 0 {
			fmt.Println("没有找到用户")
			return nil
		}

		outputTable(users)

		if verbose {
			fmt.Printf("\n共返回 %d 条记录 (偏移量: %d)\n", len(users), offset)
		}

		return nil
	},
}

func init() {
	// 添加子命令
	UserCmd.AddCommand(listCmd)

	listCmd.Flags().IntVarP(&limit, "limit", "l", 20, "每页显示数量")
	listCmd.Flags().IntVarP(&offset, "offset", "o", 0, "偏移量")
}

func queryUsers(db database.Database, limit, offset int) ([]User, error) {
	query := "SELECT id, username, nick_name, email, create_time FROM sys_user ORDER BY id ASC LIMIT ? OFFSET ?"

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.UserName, &user.NickName, &user.Email, &user.CreatedTime); err != nil {
			return nil, fmt.Errorf("扫描数据失败: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func outputTable(users []User) {
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\t账号\t昵称\t邮箱\t创建时间")
	fmt.Fprintln(w, "---\t----\t----\t----\t--------")

	for _, user := range users {
		createdTime := user.CreatedTime.Format("2006-01-02 15:04:05")
		email := "---"
		if user.Email.Valid {
			email = user.Email.String
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", user.ID, user.UserName, user.NickName, email, createdTime)
	}
	w.Flush()
}
