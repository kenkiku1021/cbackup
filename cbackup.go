package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type BackupAdapter interface {
	BackupFile(filename string) error
}

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal("Cannot get current user", err)
	}
	configFile, err := GetDefaultConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := LoadConfig(configFile)
	if err != nil {
		log.Fatal("Cannot load config file", err)
	}

	initiateBackup(cfg)
	fmt.Printf("user dir: %v\n", user.HomeDir)
}

func initiateBackup(cfg Config) {
	adapter := NewGcpBackup(cfg)
	for i := 0; i < len(cfg.BackupDir); i++ {
		backupDir(cfg.BackupDir[i], adapter)
	}
}

func backupDir(path string, adapter BackupAdapter) {
	dir, err := os.Open(path)
	if err != nil {
		log.Fatal("Cannot open path", err)
	}
	bufSize := 10

	files, err := dir.Readdir(bufSize)
	for len(files) > 0 {
		if err != nil {
			log.Fatal(err)
		}
		count := len(files)
		for i := 0; i < count; i++ {
			filename := filepath.Join(path, files[i].Name())
			if(files[i].IsDir()) {
				backupDir(filename, adapter)
			} else {
				adapter.BackupFile(filename)
			}

		}
		files, err = dir.Readdir(bufSize)
	}
}
