package postgres

import (
	"fmt"
)

// DSNFactory dsn factory
type DSNFactory struct{}

// NewDSNFactory ...
func NewDSNFactory() *DSNFactory {
	return &DSNFactory{}
}

// Create creates dsn (Data Source Name) string
func (f DSNFactory) Create(
	host string,
	port int,
	user string,
	password string,
	dbname string,
	connectTimeout int,
) string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable connect_timeout=%d",
		host, port, user, dbname, connectTimeout,
	)
	if password != "" {
		dsn = dsn + " password=" + password
	}

	return dsn
}
