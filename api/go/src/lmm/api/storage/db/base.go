package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var (
	ErrConnDone = sql.ErrConnDone
	ErrNoRows   = sql.ErrNoRows
	ErrTxDonw   = sql.ErrTxDone
)

type base struct {
	src *sql.DB
}

func newBase(driver string, config Config) DB {
	if config.User == "" {
		config.User = "root"
	}
	if config.Host == "" {
		config.Host = "127.0.0.1"
	}
	if config.Port == "" {
		config.Port = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	var (
		db    DB
		err   error
		retry = config.Retry
	)

	for {
		db, err = tryToOpenDB(driver, dsn)
		if err == nil {
			break
		}

		if retry == 0 {
			break
		} else {
			zap.L().Warn("retry connecting to db...",
				zap.String("host", config.Host),
				zap.String("port", config.Port),
				zap.String("db", config.Database),
			)
			<-time.After(5 * time.Second)
		}
		if retry > 0 {
			retry--
		}
	}
	if err != nil {
		zap.L().Panic(err.Error())
	}
	return db
}

func tryToOpenDB(driver, dsn string) (DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return &base{src: db}, nil
}

func (db *base) Begin(c context.Context, opts *sql.TxOptions) (Tx, error) {
	txn, err := db.src.BeginTx(c, opts)
	if err != nil {
		return nil, err
	}
	return &tx{Tx: txn}, nil
}

func (db *base) Conn(c context.Context) (*sql.Conn, error) {
	return db.src.Conn(c)
}

func (db *base) Close() error {
	return db.src.Close()
}

func (db *base) Driver() driver.Driver {
	return db.src.Driver()
}

func (db *base) Exec(c context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.src.ExecContext(c, query, args...)
}

func (db *base) Ping(c context.Context) error {
	return db.src.PingContext(c)
}

func (db *base) Prepare(c context.Context, query string) Stmt {
	st, err := db.src.PrepareContext(c, query)
	if err != nil {
		panic(err)
	}
	return &stmt{Stmt: st}
}

func (db *base) Query(c context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.src.QueryContext(c, query, args...)
}

func (db *base) QueryRow(c context.Context, query string, args ...interface{}) *sql.Row {
	return db.src.QueryRowContext(c, query, args...)
}

func (db *base) SetConnMaxLifetime(d time.Duration) {
	db.src.SetConnMaxLifetime(d)
}

func (db *base) SetMaxIdleConns(n int) {
	db.src.SetMaxIdleConns(n)
}

func (db *base) SetMaxOpenConns(n int) {
	db.src.SetMaxIdleConns(n)
}

func (db *base) Stats() sql.DBStats {
	return db.src.Stats()
}