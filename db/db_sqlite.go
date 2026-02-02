package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "github.com/uptrace/bun/dialect/sqlitedialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var sqliteDBClient DBClient

func NewDBClientSqlite(config *DBConfig) (DBClient, error) {
	createDB := func(dsn string) (bun.IDB, error) {
		sqlDB, err := sql.Open(sqliteshim.ShimName, dsn)
		if err != nil {
			panic(err)
		}

		// sqlDB.SetMaxIdleConns(10)
		// sqlDB.SetMaxOpenConns(5)
		// sqlDB.SetConnMaxLifetime(0)

		err = sqlDB.Ping()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		bunDB := bun.NewDB(sqlDB, sqlitedialect.New())

		bunOptions := make([]bundebug.Option, 0)
		if config != nil && config.Debug {
			bunOptions = append(bunOptions, bundebug.WithVerbose(true))
		}

		bunDB.AddQueryHook(
			bundebug.NewQueryHook(
				bunOptions...,
			),
		)

		return bunDB, nil
	}

	createDBInstance := false
	if sqliteDBClient == nil {
		createDBInstance = true
	}

	if sqliteDBClient != nil && sqliteDBClient.ConnReader().NewAddColumn().DB().Ping() != nil {
		createDBInstance = true
	}

	if createDBInstance {
		dsnReader := "file:" + config.Read.Path + "?cache=shared"
		dsnWriter := "file:" + config.Write.Path + "?cache=shared"

		dbReader, err := createDB(dsnReader)
		if err != nil {
			return nil, err
		}

		dbWriter, err := createDB(dsnWriter)
		if err != nil {
			return nil, err
		}

		sqliteDBClient = &SqliteDBClient{
			dbReader: dbReader,
			dbWriter: dbWriter,
		}
	}

	return sqliteDBClient, nil
}

type SqliteDBClient struct {
	dbReader bun.IDB
	dbWriter bun.IDB
	tx       *bun.Tx
}

func (slf *SqliteDBClient) Ping(dbType string) (err error) {
	switch dbType {
	case DBRead:
		if db, ok := slf.dbReader.(*bun.DB); ok {
			err = db.Ping()
		}

	case DBWrite:
		if db, ok := slf.dbWriter.(*bun.DB); ok {
			err = db.Ping()
		}
	}

	return
}

func (slf *SqliteDBClient) ConnReader() bun.IDB {
	if slf.tx != nil {
		_ = slf.Ping(DBWrite)
		return slf.tx
	}

	_ = slf.Ping(DBRead)
	return slf.dbReader
}

func (slf *SqliteDBClient) ConnWriter() bun.IDB {
	_ = slf.Ping(DBWrite)

	if slf.tx != nil {
		return slf.tx
	}

	return slf.dbWriter
}

func (slf *SqliteDBClient) BeginTx(ctx context.Context, opts *sql.TxOptions) error {
	tx, err := slf.dbWriter.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	slf.tx = &tx

	return nil
}

func (slf *SqliteDBClient) Rollback() error {
	if slf.tx == nil {
		return nil
	}

	err := slf.tx.Rollback()

	slf.tx = nil

	return err
}

func (slf *SqliteDBClient) Commit() error {
	err := slf.tx.Commit()

	slf.tx = nil

	return err
}

func (slf *SqliteDBClient) Close() error {
	// slf.dbReader.Close()
	// slf.dbWriter.Close()
	return nil
}
