package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 定义配置结构体（与 yaml 配置字段对应）
type AppConfig struct {
	HTTP struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"http"`
	Mysql struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"mysql"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
}

var (
	cfg     AppConfig // 全局配置实例
	envFlag string    // 接收环境参数的变量（dev/beta/pre）
	rootCmd = &cobra.Command{
		Use:   "app",
		Short: "一个基于 Cobra + Viper 的简易配置演示程序",
		Long:  "通过环境变量指定运行环境，自动读取对应 yaml 配置文件",
		Run:   runRootCmd, // 命令执行入口
	}
)

// runRootCmd：命令执行逻辑（读取配置 + 打印配置）
func runRootCmd(cmd *cobra.Command, args []string) {
	// 1. 配置 Viper 读取对应环境的配置文件
	viper.SetConfigName(envFlag)     // 配置文件名（不带后缀，对应 dev/beta/pre）
	viper.SetConfigType("yaml")      // 配置文件类型
	viper.AddConfigPath("./configs") // 配置文件存放目录（项目根目录下的 configs）

	// 2. 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("读取配置文件失败：%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("成功读取配置文件：%s\n", viper.ConfigFileUsed())

	// 3. 将配置反序列化到 AppConfig 结构体
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("解析配置文件失败：%v\n", err)
		os.Exit(1)
	}

	// 4. 打印配置信息（验证读取结果）
	fmt.Println("------------------- 配置信息 -------------------")
	fmt.Printf("HTTP 地址：%s\n", cfg.HTTP.Addr)
	fmt.Printf("MySQL DSN：%s\n", cfg.Mysql.DSN)
	fmt.Printf("日志级别：%s\n", cfg.Log.Level)
	fmt.Println("-----------------------------------------------")
}

// InitRootCmd：初始化根命令（对外暴露，供 main.go 调用）
func InitRootCmd() *cobra.Command {
	// 定义命令行参数（--env / -e），接收环境值，默认 dev
	rootCmd.Flags().StringVarP(&envFlag, "env", "e", "dev", "指定运行环境（dev/beta/pre）")

	// 支持通过环境变量 APP_ENV 指定（优先级：命令行参数 > 环境变量 > 默认值）
	viper.BindPFlag("env", rootCmd.Flags().Lookup("env"))
	viper.SetEnvPrefix("APP") // 环境变量前缀：APP_
	viper.BindEnv("env")      // 绑定环境变量 APP_ENV
	envFlag = viper.GetString("env")

	return rootCmd
}
