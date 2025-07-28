package environment

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type EnvConfig = *viper.Viper

var Config = viper.New()

func InitYamlConfig(file string) EnvConfig {
	fileBase := path.Base(file)
	fileExt := path.Ext(fileBase)
	path, _ := strings.CutSuffix(file, fileBase)
	name, _ := strings.CutSuffix(fileBase, fileExt)

	fmt.Println("path: ", path, "name: ", name, "ext: ", fileExt)

	Config.SetConfigName(name)
	Config.SetConfigType(fileExt)
	Config.AddConfigPath(path)

	if err := Config.ReadInConfig(); err != nil {
		log.Fatal("Не удалось выполнить чтение config", file)
	}

	InitEnvConfig()

	return Config
}

func InitEnvConfig() EnvConfig {
	Config.SetConfigFile(".env")
	Config.SetConfigType("env")

	if err := Config.MergeInConfig(); err != nil {
		log.Fatal("Не удалось прочитать .env")
	}

	replaceDots()

	return Config
}

func replaceDots() {
	for _, key := range Config.AllKeys() {
		dotKey := strings.ReplaceAll(key, "_", ".")
		if dotKey != key {
			Config.RegisterAlias(dotKey, key)
		}
	}
}
