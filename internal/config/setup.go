package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
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
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	v.SetEnvPrefix("app")
	v.AutomaticEnv()

	bindStructKeys(v, reflect.TypeOf(AppConfig{}), "")

	v.AddConfigPath(path)
	if path == "." {
		v.AddConfigPath("..")
		v.AddConfigPath("../..")
	}
	v.SetConfigName(name)

	err = v.ReadInConfig()
	if err != nil {
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundErr) && !os.IsNotExist(err) {
			return config, err
		}
		err = nil
	}

	for _, k := range v.AllKeys() {
		if v.IsSet(k) {
			v.Set(k, v.Get(k))
		}
	}

	err = v.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()
	return config, nil
}

func bindStructKeys(v *viper.Viper, t reflect.Type, prefix string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" || tag == "-" {
			continue
		}

		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Struct {
			bindStructKeys(v, fieldType, key)
		} else {
			_ = v.BindEnv(key)
		}
	}
}

