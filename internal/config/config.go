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

type Config struct {
	Http      httpCfg   `yaml:"http"`
	Logger    loggerCgf `yaml:"logger"`
	Auth      authCfg   `yaml:"auth_service"`
	AppParams appCfg    `yaml:"application_params"`
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
