package generator

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func ConfigInit() {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("Get os.Getwd error : %s", err))
	}

	viper.SetConfigName("bootConfig")
	viper.AddConfigPath(dir)
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}
