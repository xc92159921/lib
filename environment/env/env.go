package env

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type EnvConfig = *viper.Viper

var Config = viper.New()

func init() {
	Config.SetConfigFile(".env")
	Config.SetConfigType("env")

	if err := Config.MergeInConfig(); err != nil {
		log.Fatal("Не удалось прочитать .env")
	}

	replaceDots()
}

func replaceDots() {
	for _, key := range Config.AllKeys() {
		dotKey := strings.ReplaceAll(key, "_", ".")
		if dotKey != key {
			Config.RegisterAlias(dotKey, key)
		}
	}
}
