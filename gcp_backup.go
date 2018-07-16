package main

import (
	"io"
	"os"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/datastore"
        "golang.org/x/net/context"
	_ "google.golang.org/api/iterator"
)

type GcpBackup struct {
	projectId string
	ctx context.Context
	dsClient *datastore.Client
	storageClient *storage.Client
	bucketName string
	bucket *storage.BucketHandle
	dsKind string
}

func NewGcpBackup(cfg Config) *GcpBackup {
	backup := GcpBackup{}
	backup.projectId = cfg.GcpProjectID
	backup.ctx = context.Background()
	dsClient, err := datastore.NewClient(backup.ctx, backup.projectId)
	if err != nil {
		logger.Fatalf("Cannot create Google Datastore Client object (project id: %v): %v", backup.projectId, err)
	}
	backup.dsClient = dsClient
	backup.dsKind = cfg.GcpKind
	storageClient, err := storage.NewClient(backup.ctx)
	if err != nil {
		logger.Fatalf("Cannot create Google Cloud Storage Client object: %v", err)
	}
	backup.storageClient = storageClient
	backup.bucketName = cfg.BucketName
	bucket := backup.storageClient.Bucket(cfg.BucketName)
	backup.bucket = bucket

	return &backup
}

func (backup GcpBackup) BackupFile(filename string) error {
	var fdb FileDB
	logger.Printf("[info] Processing file %v", filename)
	hash, err := MakeFileHash(filename)
	if err != nil {
		logger.Fatalf("Cannot open file %v : %v", filename, err)
	}
	stat, err := os.Stat(filename)
	if err != nil {
		logger.Fatalf("Cannot stat file %v : %v", filename, err)
	}
	key := datastore.NameKey(backup.dsKind, hash, nil)
	err = backup.dsClient.Get(backup.ctx, key, &fdb)
	if err != nil {
		// no data in datastore
		err = backup.UploadFile(filename, hash)
		if err != nil {
			logger.Printf("[warning] Upload file error (%v): %v",
				filename, err)
			return nil
		} else {
			fdb.Paths = AppendPath(fdb.Paths, filename)
			fdb.Size = stat.Size()
			fdb.ModTime = stat.ModTime()
		}
	} else {
		fdb.Paths = AppendPath(fdb.Paths, filename)
	}
	key, err = backup.dsClient.Put(backup.ctx, key, &fdb)
	if err != nil {
		logger.Printf("[warning] Cannot put filedb (%v): %v",
			filename, err)
	}

	return nil
}

func (backup GcpBackup) UploadFile(filename string, hash string) error {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	
	wc := backup.bucket.Object(hash).NewWriter(backup.ctx)
	logger.Printf("[info] Uploading file to GCP (%v)", filename)
	if _, err = io.Copy(wc, f); err != nil {
		logger.Printf("[warning] Copy error (%v): %v", filename, err)
		return err
	}
	if err = wc.Close(); err != nil {
		logger.Printf("[warning] Close error (%v): %v", filename, err)
		return err
	}
	
	return nil
}

