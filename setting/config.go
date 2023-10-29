package setting

import (
	"fmt"
	"go-micro-toDoList/global"

	"log"
	"os"

	"github.com/spf13/viper"
)

func InitConfig() {
	if global.Config != nil {
		return
	}
	workDir, _ := os.Getwd()
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&global.Config)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Configuration file loads successfully.")
}
