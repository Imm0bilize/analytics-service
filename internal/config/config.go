package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type httpCfg struct {
	Port         string        `yaml:"port"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
}

type appCfg struct {
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type loggerCgf struct {
	Level    string `yaml:"level"`
	TsFormat string `yaml:"ts_format"`
}

type authCfg struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type grpcCfg struct {
	Network string `yaml:"network"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

type kafkaCfg struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Topic    string `yaml:"topic"`
	ClientID string `yaml:"client_id"`
}

type dbCfg struct {
	IsNeedMigration    bool   `yaml:"is_need_migration"`
	NAttemptsToConnect int    `yaml:"n_attempts_to_connect"`
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	User               string `yaml:"user" env:"DB_POSTGRES_USER" `
	Password           string `yaml:"password" env:"DB_POSTGRES_PASSWORD" `
	DbName             string `yaml:"db_name" env:"DB_POSTGRES_DBNAME"`
}

type sentryCfg struct {
	Debug        bool          `yaml:"debug"`
	Dsn          string        `yaml:"dsn"`
	FlushTimeout time.Duration `yaml:"flush_timeout"`
}

type prometheusCfg struct {
	Port string `yaml:"port"`
}

type Config struct {
	Http       httpCfg       `yaml:"http"`
	Logger     loggerCgf     `yaml:"logger"`
	Auth       authCfg       `yaml:"auth_service"`
	Db         dbCfg         `yaml:"database"`
	Grpc       grpcCfg       `yaml:"grpc_server"`
	Kafka      kafkaCfg      `yaml:"kafka"`
	Sentry     sentryCfg     `yaml:"sentry"`
	Prometheus prometheusCfg `yaml:"prometheus"`
	AppParams  appCfg        `yaml:"application_params"`
}

func New(path string) (*Config, error) {
	var config Config

	if err := readParamsFromConfigFile(path, &config); err != nil {
		return nil, err
	}
	if err := readParamsFromEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func GetPathToCfgFromEnv() (string, error) {
	err := godotenv.Load()
	if err == nil {
		log.Println(".env file found")
	}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return "", ErrPathNotSpecified
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", ErrGetFullPath
	}
	return absPath, nil
}

func readParamsFromConfigFile(pathToFile string, config *Config) error {
	file, err := os.Open(pathToFile)
	if err != nil {
		return fmt.Errorf("error: %s\nInfo: %w", ErrOpenCfgFile.Error(), err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(ErrCloseCfgFile)
		}
	}(file)

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error: %s\t%w", ErrReadCfgFile.Error(), err)
	}

	if err = yaml.Unmarshal(fileContent, config); err != nil {
		return err
	}
	return nil
}

func readParamsFromEnv(config *Config) error {
	if err := envconfig.Process("", config); err != nil {
		return ErrParseFile
	}
	return nil
}
