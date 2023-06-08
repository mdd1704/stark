package database

import (
	"fmt"
	"os"
	"stark/utils/activity"
	"stark/utils/log"

	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"
)

type Block func(tx *sqlx.Tx) (result interface{}, err error)

type MySQL struct {
	db *sqlx.DB
}

func NewMySQL() (*MySQL, error) {
	ctx := activity.NewContext("init_mysql")
	ctx = activity.WithClientID(ctx, "stark_system")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)

	initMigration(connectionString)
	db, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't open mysql connection"))
		return nil, stacktrace.Propagate(err, "can't open mysql connection")
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't ping mysql db"))
		return nil, stacktrace.Propagate(err, "can't ping mysql db")
	}

	return &MySQL{db}, nil
}

func (m *MySQL) WithTransaction(block Block) (interface{}, error) {
	tx, err := m.db.Beginx()
	if err != nil {
		return nil, stacktrace.Propagate(err, "can't start mysql DB transaction")
	}

	result, err := block(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, stacktrace.Propagate(err, "can't rollback mysql DB transaction")
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, stacktrace.Propagate(err, "commit transaction mysql DB fails")
	}

	return result, nil
}

func (m *MySQL) In(query string, params map[string]interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return "", nil, err
	}

	return sqlx.In(query, args...)
}

func (m *MySQL) Get(dest interface{}, query string, args ...interface{}) error {
	return m.db.Get(dest, query, args...)
}

func (m *MySQL) Select(dest interface{}, query string, args ...interface{}) error {
	return m.db.Select(dest, query, args...)
}

func (m *MySQL) Rebind(query string) string {
	return m.db.Rebind(query)
}
