package main

import (
	"log"
	"os"
	"path/filepath"
)

type BackupAdapter interface {
	BackupFile(filename string) error
}

func main() {
	configFile, err := GetDefaultConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := LoadConfig(configFile)
	if err != nil {
		log.Fatal("Cannot load config file", err)
	}
	InitLogger(cfg)

	initiateBackup(cfg)
	logger.Print("Backup completed")
	log.Print("Backup completed")
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
		logger.Fatal("Cannot open path", err)
	}
	bufSize := 10

	files, err := dir.Readdir(bufSize)
	for len(files) > 0 {
		if err != nil {
			logger.Fatal(err)
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
