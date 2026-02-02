package db

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var pgDBClient DBClient

func NewDBClientPostgres(config *DBConfig) (DBClient, error) {
	createDB := func(dsn string) (bun.IDB, error) {
		sqlDB, err := sql.Open("pg", dsn)
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

		bunDB := bun.NewDB(sqlDB, pgdialect.New())

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
	if pgDBClient == nil {
		createDBInstance = true
	}

	if pgDBClient != nil && pgDBClient.ConnReader().NewAddColumn().DB().Ping() != nil {
		createDBInstance = true
	}

	if createDBInstance {
		dsnReader := "postgres://" + config.Read.Username + ":" + config.Read.Password + "@" + config.Read.Host + ":" +
			strconv.Itoa(config.Read.Port) + "/" + config.Read.DBName + "?sslmode=disable"
		dsnWriter := "postgres://" + config.Write.Username + ":" + config.Write.Password + "@" + config.Write.Host + ":" +
			strconv.Itoa(config.Write.Port) + "/" + config.Write.DBName + "?sslmode=disable"

		dbReader, err := createDB(dsnReader)
		if err != nil {
			return nil, err
		}

		dbWriter, err := createDB(dsnWriter)
		if err != nil {
			return nil, err
		}

		pgDBClient = &PostgresDBClient{
			dbReader: dbReader,
			dbWriter: dbWriter,
		}
	}

	return pgDBClient, nil
}

type PostgresDBClient struct {
	dbReader bun.IDB
	dbWriter bun.IDB
	tx       *bun.Tx
}

func (slf *PostgresDBClient) Ping(dbType string) (err error) {
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

func (slf *PostgresDBClient) ConnReader() bun.IDB {
	if slf.tx != nil {
		_ = slf.Ping(DBWrite)
		return slf.tx
	}

	_ = slf.Ping(DBRead)
	return slf.dbReader
}

func (slf *PostgresDBClient) ConnWriter() bun.IDB {
	_ = slf.Ping(DBWrite)

	if slf.tx != nil {
		return slf.tx
	}

	return slf.dbWriter
}

func (slf *PostgresDBClient) BeginTx(ctx context.Context, opts *sql.TxOptions) error {
	tx, err := slf.dbWriter.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	slf.tx = &tx

	return nil
}

func (slf *PostgresDBClient) Rollback() error {
	if slf.tx == nil {
		return nil
	}

	err := slf.tx.Rollback()

	slf.tx = nil

	return err
}

func (slf *PostgresDBClient) Commit() error {
	err := slf.tx.Commit()

	slf.tx = nil

	return err
}

func (slf *PostgresDBClient) Close() error {
	// slf.dbReader.Close()
	// slf.dbWriter.Close()
	return nil
}
