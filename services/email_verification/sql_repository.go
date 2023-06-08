package email_verification

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"

	"stark/database"
)

type sqlRepository struct {
	mysqlDB *database.MySQL
}

const (
	selectCountEmailVerificationQuery = "SELECT COUNT(*) FROM email_verifications"
	insertEmailVerificationQuery      = `
		INSERT INTO email_verifications (email, token, created_at) 
		VALUES (?, ?, ?)
	`
)

func NewSQLRepository(mysqlDB *database.MySQL) Repository {
	return &sqlRepository{mysqlDB}
}

func (repo *sqlRepository) Store(data *EmailVerification) error {
	return repo.insert(data)
}

func (repo *sqlRepository) FindByToken(token string) (result *EmailVerification, err error) {
	var data EmailVerification
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("email_verifications")
	dataset = dataset.Where(goqu.Ex{
		"token": token,
	})

	sql, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "sql error")
	}

	err = repo.mysqlDB.Get(&data, sql)
	if err != nil {
		return nil, stacktrace.Propagate(err, "can't read email verification by token")
	}

	return &data, nil
}

func (repo *sqlRepository) FindTotalByFilter(filter Filter) (total int, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("email_verifications")
	dataset = dataset.Select(goqu.COUNT("*"))
	if len(filter.Tokens) != 0 {
		dataset = dataset.Where(goqu.ExOr{
			"token": filter.Tokens,
		})
	}

	sql, _, err := dataset.ToSQL()
	if err != nil {
		return 0, stacktrace.Propagate(err, "sql error")
	}

	err = repo.mysqlDB.Get(&total, sql)
	if err != nil {
		return 0, stacktrace.Propagate(err, "select row fails")
	}

	return total, nil
}

func (repo *sqlRepository) insert(data *EmailVerification) error {
	_, err := repo.mysqlDB.WithTransaction(func(tx *sqlx.Tx) (interface{}, error) {
		res, err := tx.Exec(insertEmailVerificationQuery,
			data.Email,
			data.Token,
			data.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		if rowsAffected <= 0 {
			return nil, errors.New("insert email verification fails")
		}

		return nil, nil
	})

	return err
}
