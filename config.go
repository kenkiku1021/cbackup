package main

import (
	"io/ioutil"
	"path/filepath"
	"os"
	"log"
	"os/user"
	"gopkg.in/yaml.v2"
)

const appDir = ".cbackup"
const defaultConfigFile = "cbackup_config.yaml"

type Config struct {
	AppDir string `yaml:"app_dir"`
	CfgDir string
	LogDir string
	BackupDir []string `yaml:"backup_dir"`
	GcpProjectID string `yaml: "gcpprojectid"`
	GcpCredentials string `yaml:"gcpcredentials"`
	GcpKind string `yaml: "gcpkind"`
	BucketName string `yaml: "bucketname"`
}

func GetDefaultConfigFilePath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	path := filepath.Join(user.HomeDir, appDir, defaultConfigFile)
	return path, nil
}

func LoadConfig(configFile string) (Config, error) {
	cfgdir := filepath.Dir(configFile)
	cfg := Config{}
	cfg.CfgDir = cfgdir
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("Cannot load config file (%v)", configFile)
		return cfg, err
	}
	yaml.Unmarshal(data, &cfg)

	if cfg.GcpCredentials != "" {
		gcpCredentialsFile := filepath.Join(cfgdir, cfg.GcpCredentials)
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", gcpCredentialsFile)
	}

	cfg.LogDir = filepath.Join(cfg.CfgDir, "log")
	err = os.MkdirAll(cfg.LogDir, 0755)
	if err != nil {
		log.Fatal("Cannot create log dir (%v): %v", cfg.LogDir, err)
	}

	return cfg, nil
}

