package database

import (
	"database/sql"
	"io/ioutil"
	"os"
	"runtime"
	"stark/database/migrations"
	"stark/utils/activity"
	"stark/utils/log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/palantir/stacktrace"
)

const INIT_STEP = 7
const APP_SCHEMA_VERSION = 7

var seeds = []string{
	"user",
}

func initMigration(connectionString string) *sql.DB {
	ctx := activity.NewContext("init_mysql_migration")
	ctx = activity.WithClientID(ctx, "stark_system")

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "mysql connection error"))
		return nil
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "init instance db migration error"))
		return nil
	}

	sqlStatement := "SELECT version, dirty FROM schema_migrations"
	stmt, err := db.Prepare(sqlStatement)

	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "select schema migration error"))
		err = m.Steps(INIT_STEP)
		return db
	} else {
		row := stmt.QueryRow()

		var schema migrations.Schema_migration
		err = row.Scan(&schema.Version, &schema.Dirty)
		deltaSchema := APP_SCHEMA_VERSION - schema.Version

		switch err {
		case sql.ErrNoRows:
			err = m.Steps(INIT_STEP)
			if err != nil {
				log.WithContext(ctx).Error(stacktrace.Propagate(err, "init step migration error"))
				return db
			}

			dir, err := os.Getwd()
			if err != nil {
				log.WithContext(ctx).Error(stacktrace.Propagate(err, "get dir migration db error"))
				return db
			}

			if runtime.GOOS == "windows" {
				for _, seed := range seeds {
					file := dir + "\\migrations\\seeds\\" + seed + ".sql"
					sqls, _ := ioutil.ReadFile(file)
					sqlArr := strings.Split(string(sqls), ";")
					for i := 0; i < len(sqlArr)-1; i++ {
						sqlQuery := sqlArr[i]
						_, err := db.Exec(sqlQuery)
						if err != nil {
							log.WithContext(ctx).Error(stacktrace.Propagate(err, "seeding database error"))
							return db
						}
					}
				}
			} else {
				for _, seed := range seeds {
					file := dir + "/migrations/seeds/" + seed + ".sql"
					sqls, _ := ioutil.ReadFile(file)
					sqlArr := strings.Split(string(sqls), ";")
					for i := 0; i < len(sqlArr)-1; i++ {
						sqlQuery := sqlArr[i]
						_, err := db.Exec(sqlQuery)
						if err != nil {
							log.WithContext(ctx).Error(stacktrace.Propagate(err, "seeding database error"))
							return db
						}
					}
				}
			}

			return db
		case nil:
			err = m.Steps(deltaSchema)
			if err != nil {
				if err.Error() != "no change" {
					log.WithContext(ctx).Error(stacktrace.Propagate(err, "init step migration error"))
					return db
				}
			}
			return db
		default:
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "migration mysql error"))
			return db
		}
	}
}
