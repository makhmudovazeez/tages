package svc

import (
	"github.com/makhmudovazeez/tages/internal/config"
	"github.com/makhmudovazeez/tages/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	FileModel models.FileModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(sqlite.Open(c.SqlLite.File), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:    c,
		FileModel: models.NewFileModel(db),
	}
}
