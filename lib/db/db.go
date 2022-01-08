package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Db interface
type Db interface {
	MakeConnection(dbURL string) (*gorm.DB, error)
}

// New imple
func New(dbURL string) (*gorm.DB, error) {
	return gorm.Open("postgres", dbURL)
	/*
		return &connectionAdapter{
			client: dbURL,
		}
	*/
}

type connectionAdapter struct {
	client string
}

func (m *connectionAdapter) MakeConnection(dbURL string) (*gorm.DB, error) {
	return gorm.Open("postgres", dbURL)
}
