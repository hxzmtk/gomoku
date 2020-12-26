package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Model struct {
	ID        int `gorm:"primarykey"`
	CreatedAt time.Time
}

func (m *Model) BeforeCreate(db *gorm.DB) error {
	m.CreatedAt = time.Now()
	return nil
}

const DB_NAME = "localhost.db"

var DB *gorm.DB

func Start() {
	var err error
	DB, err = gorm.Open(sqlite.Open(DB_NAME), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
