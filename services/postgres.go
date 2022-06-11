package services

import (
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresDB *PostgresClient

type PostgresClient struct {
	Client *gorm.DB
}

func NewPostgresClient(dns string) (*PostgresClient, error) {
	db, err := gorm.Open(postgres.Open(dns))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &PostgresClient{
		Client: db,
	}, nil
}
