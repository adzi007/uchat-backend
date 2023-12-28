package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	MONGODB_IRI string
	PORT_AP     string
}

var (
	ENV        Config
	_, b, _, _ = runtime.Caller(0)

	ProjectRootPath = filepath.Join(filepath.Dir(b), "../")
)

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		log.Fatal(err)
	}

	log.Println("Load server successfully")
}
