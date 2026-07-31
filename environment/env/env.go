package env

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type EnvConfig = *viper.Viper

var Config = viper.New()

func init() {
	// 1. Загрузка основного файла .env
	Config.SetConfigFile(".env")
	Config.SetConfigType("env")

	if err := Config.MergeInConfig(); err != nil {
		log.Fatal("Не удалось прочитать .env")
	}

	// 2. Наложение .env.local, если он существует
	if _, err := os.Stat(".env.local"); err == nil {
		Config.SetConfigFile(".env.local")
		Config.SetConfigType("env")

		if err := Config.MergeInConfig(); err != nil {
			log.Printf("Предупреждение: .env.local найден, но не прочитан: %v", err)
		}
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
