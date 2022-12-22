package services

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sync"
)

type (
	FileStore interface {
		Save(fileType string, fileData bytes.Buffer) (string, error)
	}

	DiskFileStore struct {
		mutex      sync.RWMutex
		fileFolder string
	}

	FileInfo struct {
		Type string
		Name string
	}
)

func NewFileStore(fileFolder string) *DiskFileStore {
	return &DiskFileStore{
		fileFolder: fileFolder,
	}
}

func (store *DiskFileStore) Save(fileType string, fileData bytes.Buffer) (string, error) {
	fileId := uuid.New().String()

	fileName := fmt.Sprintf("%s%s", fileId, fileType)

	file, err := os.Create(fmt.Sprintf("%s/%s", store.fileFolder, fileName))
	if err != nil {
		return "", fmt.Errorf("cannot store file: %w", err)
	}
	defer file.Close()

	if _, err = fileData.WriteTo(file); err != nil {
		return "", fmt.Errorf("cannot write to file: %w", err)
	}
	return fileId, nil
}
