package cmd

import (
	"errors"


	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tx991020/go-boot/generator"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "生成一套完整的gin restful 框架",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		httpPort, _ := cmd.Flags().GetString("http_port")
		typeStr, _ := cmd.Flags().GetString("type")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")
		sslmode, _ := cmd.Flags().GetString("sslmode")
		table, _ := cmd.Flags().GetStringSlice("table")
		isOpenMemoryCache, _ := cmd.Flags().GetBool("memory_cache")
		appPath, _ := cmd.Flags().GetString("appPath")
		tplPath, _ := cmd.Flags().GetString("tplPath")
		if name == "" {
			generator.PrintErr(errors.New("请输入项目名称 -n"))
			return
		}
		if password == "" {
			generator.PrintErr(errors.New("请输入数据库密码 -p"))
			return
		}
		if database == "" {
			generator.PrintErr(errors.New("请输入数据库库名 -d"))
			return
		}
		generator.PrintInfo("------------------开始------------------")
		generator.PrintInfo("cli version: ", viper.GetString("version"))
		generator.PrintInfo("输入的生成信息为: ")
		generator.PrintInfoMap(map[string]interface{}{

			"name":         name,
			"appPath":      appPath,
			"tplPath":      tplPath,
			"type":         typeStr,
			"host":         host,
			"port":         port,
			"user":         user,
			"password":     password,
			"database":     database,
			"sslmode":      sslmode,
			"http_port":    httpPort,
			"memory_cache": isOpenMemoryCache,
		})

		generator.Run(&generator.InputParams{

			Name:              name,
			AppPath:           appPath,
			TemplatePath:      tplPath,
			DBType:            typeStr,
			Host:              host,
			Port:              port,
			User:              user,
			Password:          password,
			Database:          database,
			Sslmode:           sslmode,
			HTTPPort:          httpPort,
			Table:             table,
			IsOpenMemoryCache: isOpenMemoryCache,
		})
		generator.PrintInfo("------------------结束------------------")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.MarkFlagRequired("name")
	newCmd.MarkFlagRequired("appPath")
	newCmd.MarkFlagRequired("tplPath")
	newCmd.MarkFlagRequired("http_port")
	newCmd.MarkFlagRequired("type")
	newCmd.MarkFlagRequired("host")
	newCmd.MarkFlagRequired("driver")
	newCmd.MarkFlagRequired("port")
	newCmd.MarkFlagRequired("user")
	newCmd.MarkFlagRequired("password")
	newCmd.MarkFlagRequired("database")
	newCmd.MarkFlagRequired("sslmode")
	newCmd.MarkFlagRequired("table")
	newCmd.MarkFlagRequired("memory_cache")
}
