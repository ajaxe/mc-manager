package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func LoadAppConfig() (config AppConfig) {
	config, err := loadAppConfigInternal(".", "config")
	if err != nil {
		log.Fatalf("failed to get app config: %v", err)
	}
	return
}

func loadAppConfigInternal(path, name string) (config AppConfig, err error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.SetEnvPrefix("app")

	viper.AutomaticEnv()

	viper.AddConfigPath(path)
	viper.SetConfigName(name)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
	return
}
