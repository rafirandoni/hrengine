package db

import (
	"context"
	"database/sql"

	"github.com/rafirandoni/hrengine/constant"

	"github.com/uptrace/bun"
)

type DBClient interface {
	ConnReader() bun.IDB
	ConnWriter() bun.IDB
	BeginTx(ctx context.Context, opts *sql.TxOptions) error
	Rollback() error
	Commit() error
	Close() error
}

type DBConfig struct {
	DBConnection string
	Debug        bool
	Read         *DBConfigCredential
	Write        *DBConfigCredential
}

type DBConfigCredential struct {
	Username string
	Password string
	Host     string
	Port     int
	DBName   string
	Path     string // when using sqlite, file path need to be specified here
	Schema   string
}

const (
	DBRead  = "read"
	DBWrite = "write"
)

func NewDBClient(config *DBConfig) (DBClient, error) {
	switch config.DBConnection {
	case constant.DBConnectionPostgres:
		return NewDBClientPostgres(config)
	// case constant.DBConnectionMysql:
	// return NewDBClientMysql(config)

	case constant.DBConnectionSqlite:
		return NewDBClientSqlite(config)
	}

	panic("db driver not found")
}
