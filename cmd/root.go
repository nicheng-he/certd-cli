package cmd

import (
	"certd-cli/cmd/user"
	"certd-cli/constant"
	settingsContext "certd-cli/context"
	"certd-cli/database"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool

	rootCmd = &cobra.Command{
		Use:   constant.Name,
		Short: "一个强大的命令行工具",
		Long:  fmt.Sprintf("%s 是certd命令行工具", constant.Name),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// version 命令不需要初始化配置和数据库
			if cmd.Name() == "version" || cmd.Use == "version" {
				return nil
			}
			// 只有当不是根命令时才初始化配置
			if cmd.HasParent() {
				initConfig()

				container := database.GetContainer()

				if !container.IsInitialized() {
					return fmt.Errorf("数据库未初始化，请检查配置文件或环境变量")
				}

				db, err := container.GetDatabase()
				if err != nil {
					return fmt.Errorf("获取数据库连接失败: %w", err)
				}

				settings, err := loadSystemSettings(db)
				if err != nil {
					verbose, _ := cmd.Flags().GetBool("verbose")
					if verbose {
						fmt.Printf("警告: 加载系统设置失败: %v\n", err)
					}
					settings = &settingsContext.SettingsContext{
						RawSettings: make(map[string]json.RawMessage),
					}
				}

				ctx := settingsContext.SetSettingsFromContext(cmd, settings)
				cmd.SetContext(ctx)
				fmt.Printf("站点ID: %s\n", settings.SiteId)
				fmt.Println()
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				fmt.Println("以详细模式启动")
			}
			fmt.Printf("欢迎使用 %s！使用 --help 查看帮助\n", constant.Name)
		},
	}
)

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 全局标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")

	// 添加子命令
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(user.UserCmd)
}

// initConfig 初始化配置
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		if verbose {
			fmt.Println("未指定配置文件，使用环境变量配置")
		}
	}

	viper.AutomaticEnv()
	bindEnvVariables()

	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Println("使用配置文件:", viper.ConfigFileUsed())
		}
	} else {
		if verbose && cfgFile == "" {
			fmt.Println("使用环境变量: certd_typeorm_dataSource_default_*")
		}
	}

	// 初始化数据库容器
	if err := initializeDatabaseContainer(); err != nil {
		if verbose {
			fmt.Printf("警告: 数据库初始化失败: %v\n", err)
		}
	}
}

// bindEnvVariables 绑定环境变量到配置项
func bindEnvVariables() {
	// 数据库配置
	_ = viper.BindEnv("database.type", "certd_typeorm_dataSource_default_type")
	_ = viper.BindEnv("database.host", "certd_typeorm_dataSource_default_host")
	_ = viper.BindEnv("database.port", "certd_typeorm_dataSource_default_port")
	_ = viper.BindEnv("database.username", "certd_typeorm_dataSource_default_username")
	_ = viper.BindEnv("database.password", "certd_typeorm_dataSource_default_password")
	_ = viper.BindEnv("database.database", "certd_typeorm_dataSource_default_database")
}

func initializeDatabaseContainer() error {
	var config database.AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	if verbose {
		fmt.Println("\n=== 数据库配置 ===")
		fmt.Printf("类型: %s\n", config.Database.Type)
		if config.Database.Type == "sqlite" {
			fmt.Printf("路径: %s\n", config.Database.Path)
			if config.Database.Path == "" && config.Database.Database != "" {
				fmt.Printf("数据库字段(兼容): %s\n", config.Database.Database)
			}
		} else {
			fmt.Printf("主机: %s\n", config.Database.Host)
			fmt.Printf("端口: %d\n", config.Database.Port)
			fmt.Printf("用户名: %s\n", config.Database.Username)
			fmt.Printf("密码: %s\n", maskPassword(config.Database.Password))
			fmt.Printf("数据库名: %s\n", config.Database.Database)
		}
		fmt.Println("==================\n")
	}

	container := database.GetContainer()

	if container.IsInitialized() {
		if verbose {
			fmt.Println("数据库容器已初始化，跳过")
		}
		return nil
	}

	if verbose {
		fmt.Println("正在初始化数据库容器...")
	}

	if err := container.Initialize(&config); err != nil {
		return fmt.Errorf("初始化数据库容器失败: %w", err)
	}

	if verbose {
		fmt.Println("✓ 数据库容器初始化成功")
	}

	return nil
}

func maskPassword(password string) string {
	if password == "" {
		return "(空)"
	}
	if len(password) <= 2 {
		return "**"
	}
	return password[:1] + strings.Repeat("*", len(password)-2) + password[len(password)-1:]
}

func loadSystemSettings(db database.Database) (*settingsContext.SettingsContext, error) {
	query := "SELECT `key`, setting FROM sys_settings"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询系统设置失败: %w", err)
	}
	defer rows.Close()

	settings := &settingsContext.SettingsContext{
		RawSettings: make(map[string]json.RawMessage),
	}

	for rows.Next() {
		var key string
		var settingBytes []byte
		if err := rows.Scan(&key, &settingBytes); err != nil {
			return nil, fmt.Errorf("扫描设置数据失败: %w", err)
		}
		setting := json.RawMessage(settingBytes)
		settings.RawSettings[key] = setting

		var value map[string]interface{}

		switch key {
		case "sys.secret":
			if err := json.Unmarshal(setting, &value); err == nil {
				if siteId, ok := value["siteId"].(string); ok {
					settings.SiteId = siteId
				}
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return settings, nil
}
