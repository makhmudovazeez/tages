package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type (
	FileModel interface {
		Save(ctx context.Context, id, name string) error
		GetAll(ctx context.Context) ([]*File, error)
	}

	fileModel struct {
		db *gorm.DB
	}

	File struct {
		Id        string
		Name      string
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}
)

func NewFileModel(db *gorm.DB) FileModel {
	file := &fileModel{
		db: db,
	}
	fileM := &File{}

	if !file.db.Migrator().HasTable(fileM) {
		if err := file.db.AutoMigrate(fileM); err != nil {
			panic(err)
		}
	}

	return file
}

func (m *fileModel) Save(ctx context.Context, id, name string) error {
	file := &File{
		Id:   id,
		Name: fmt.Sprintf("%s%s", id, name),
	}
	m.db.Create(file)
	return nil
}

func (m *fileModel) GetAll(ctx context.Context) ([]*File, error) {
	files := []*File{}
	err := m.db.Find(&files).Error
	return files, err
}
