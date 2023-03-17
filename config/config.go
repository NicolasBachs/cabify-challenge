package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	App struct {
		Name        string
		Port        string
		Environment string
	}
	Database struct {
		Host             string
		Port             string
		Username         string
		Password         string
		Name             string
		ConnectionString string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	Kafka struct {
		Host   string
		Port   string
		Topics struct {
			NewCarAvailable   string
			NewPendingJourney string
		}
	}
	Etcd struct {
		Host     string
		Port     string
		Username string
		Password string
	}
}

func LoadConfig() error {
	log.Println("Loading config...")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./app/config")

	var config Config

	err := viper.ReadInConfig()
	if err != nil {
		panic("Error reading configuration. Error: " + err.Error())
	}

	log.Println("Config read successfully...")
	log.Println("Getting environment variables to replace values...")

	envPattern := regexp.MustCompile(`\${(\w+)}`)

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		replacedValue := envPattern.ReplaceAllStringFunc(value, func(match string) string {
			envKey := strings.Trim(match, "${}")
			return getEnvOrPanic(envKey)
		})
		viper.Set(k, replacedValue)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("Error decoding configuration: %v", err)
	}

	AppConfig = &config

	log.Println("Configs loaded successfully...")

	return nil
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}

func IsProductionEnv() bool {
	env := AppConfig.App.Environment
	return env == "prod" || env == "production" || env == "prd" || env == "master"
}
