package user

import (
	"certd-cli/constant"
	settingsContext "certd-cli/context"
	"certd-cli/database"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var passwordCmd = &cobra.Command{
	Use:   "password <username> [new-password]",
	Short: "重置用户密码",
	Long:  "重置指定用户的密码，如果不提供密码则自动生成随机密码",
	Example: fmt.Sprintf(`  # 为用户生成随机密码
  %s user password admin

  # 为用户设置指定密码
  %s user password admin 123456

  # 详细模式
  %s user password admin -v`, constant.Name, constant.Name, constant.Name),
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")

		settings, _ := settingsContext.GetSettingsFromContext(cmd)

		passwordVersion := 2
		username := args[0]
		var newPassword string

		if len(args) == 2 {
			newPassword = args[1]
		} else {
			newPassword = generateRandomPassword(8)
		}

		container := database.GetContainer()

		if !container.IsInitialized() {
			return fmt.Errorf("数据库未初始化，请检查配置文件或环境变量")
		}

		db, err := container.GetDatabase()
		if err != nil {
			return fmt.Errorf("获取数据库连接失败: %w", err)
		}

		if verbose {
			fmt.Printf("正在为用户 '%s' 重置密码...\n", username)
		}

		_password, err := genPassword(newPassword, passwordVersion, settings.SiteId)

		result, err := db.Exec("UPDATE sys_user SET password_version = ?, password = ? WHERE username = ?", passwordVersion, _password, username)
		if err != nil {
			return fmt.Errorf("更新密码失败: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("获取影响行数失败: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("用户 '%s' 不存在", username)
		}

		fmt.Printf("✓ 用户 '%s' 的密码已重置\n", username)
		fmt.Printf("新密码: %s\n", newPassword)

		if verbose {
			fmt.Printf("影响行数: %d\n", rowsAffected)
		}

		return nil
	},
}

func init() {
	UserCmd.AddCommand(passwordCmd)
}

func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	password := make([]byte, length)
	for i := range password {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[num.Int64()]
	}

	return string(password)
}

func genPassword(rawPassword string, passwordVersion int, siteId string) (string, error) {
	if passwordVersion <= 1 {
		hash := md5.Sum([]byte(rawPassword))
		return hex.EncodeToString(hash[:]), nil
	}

	plainPassword, err := buildPlainPassword(rawPassword, siteId)
	if err != nil {
		return "", err
	}

	salt, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	if err != nil {
		return "", fmt.Errorf("生成盐值失败: %w", err)
	}

	return string(salt), nil
}

func buildPlainPassword(rawPassword string, siteId string) (string, error) {
	if siteId == "" || len(siteId) < 5 {
		return "", fmt.Errorf("站点ID还未初始化或长度不足")
	}

	prefixSiteId := siteId[1:5]
	return rawPassword + prefixSiteId, nil
}
