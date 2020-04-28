package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "api-generator",
	Short: "gin restful 脚手架工具",
	Long: `golang restful 脚手架生成工具
示例: WebGenerator new -H 172.25.128.114 -P 3306 -u root -p Qwert123! -n gm-go -d gm1 --http_port 12000 --table challenge_chance_record,challenge_problem_record`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringP("name", "n", "", "应用名称")
	rootCmd.PersistentFlags().StringP("appPath", "a", "/src/walter/", "应用path")
	rootCmd.PersistentFlags().StringP("tplPath", "b", "./template", "模板path")
	rootCmd.PersistentFlags().StringP("http_port", "", "18000", "应用监听端口")
	rootCmd.PersistentFlags().StringP("type", "t", "postgres", "数据库类型")
	rootCmd.PersistentFlags().StringP("host", "H", "127.0.0.1", "数据库连接地址")
	rootCmd.PersistentFlags().StringP("port", "P", "5432", "数据库端口")
	rootCmd.PersistentFlags().StringP("user", "u", "postgres", "数据库用户名")
	rootCmd.PersistentFlags().StringP("password", "p", "", "数据库密码")
	rootCmd.PersistentFlags().StringP("database", "d", "", "数据库库名")
	rootCmd.PersistentFlags().StringP("sslmode", "s", "disable", "数据库连接方式")
	rootCmd.PersistentFlags().BoolP("memory_cache", "m", false, "是否启用内存级缓存")
	rootCmd.PersistentFlags().StringSliceP("table", "", []string{}, "生成几张表,没有默认全部表格;多个表使用,分割")
	rootCmd.MarkFlagRequired("name")
	rootCmd.MarkFlagRequired("password")
	rootCmd.MarkFlagRequired("database")

}

