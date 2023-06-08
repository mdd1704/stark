package client

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"

	"stark/database"
)

type sqlRepository struct {
	mysqlDB *database.MySQL
}

const (
	selectCountQuery = "SELECT COUNT(*) FROM clients"
	insertQuery      = `
		INSERT INTO clients (id, name, bearer_key, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)
	`
	updateQuery = `
		UPDATE clients SET
			name = ?,
			bearer_key = ?,
			updated_at = ?
		WHERE id = ?
	`
)

func NewSQLRepository(mysqlDB *database.MySQL) Repository {
	return &sqlRepository{mysqlDB}
}

func (repo *sqlRepository) Store(data *Client) error {
	exist, err := repo.existByID(data.ID)
	if err != nil {
		return err
	}

	if exist {
		return repo.update(data)
	}

	return repo.insert(data)
}

func (repo *sqlRepository) FindByID(id uuid.UUID) (result *Client, err error) {
	var data Client
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("clients")
	dataset = dataset.Where(goqu.Ex{
		"id": id.String(),
	})

	sql, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "sql error")
	}

	err = repo.mysqlDB.Get(&data, sql)
	if err != nil {
		return nil, stacktrace.Propagate(err, "can't read client by id")
	}

	return &data, nil
}

func (repo *sqlRepository) FindByFilter(filter Filter) (result []*Client, err error) {
	if filter.IsEmpty() {
		return
	}

	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("clients")
	if len(filter.Names) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"name": filter.Names,
		})
	}

	if len(filter.BearerKeys) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"bearer_key": filter.BearerKeys,
		})
	}

	sql, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "sql error")
	}

	err = repo.mysqlDB.Select(&result, sql)
	if err != nil {
		return nil, stacktrace.Propagate(err, "select rows fails")
	}

	return result, nil
}

func (repo *sqlRepository) FindPage(offset int, limit int) (result []*Client, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("clients")
	sql, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "sql error")
	}

	err = repo.mysqlDB.Select(&result, sql+" LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, stacktrace.Propagate(err, "select rows fails")
	}

	return
}

func (repo *sqlRepository) FindTotalByFilter(filter Filter) (total int, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("clients")
	dataset = dataset.Select(goqu.COUNT("*"))
	if len(filter.Names) != 0 {
		dataset = dataset.Where(goqu.ExOr{
			"name": filter.Names,
		})
	}

	if len(filter.BearerKeys) != 0 {
		dataset = dataset.Where(goqu.ExOr{
			"bearer_key": filter.BearerKeys,
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

func (repo *sqlRepository) existByID(id uuid.UUID) (bool, error) {
	var total int
	err := repo.mysqlDB.Get(&total, selectCountQuery+" WHERE id = ?", id)
	if err != nil {
		return false, stacktrace.Propagate(err, "select count fails")
	}

	return total > 0, nil
}

func (repo *sqlRepository) insert(data *Client) error {
	_, err := repo.mysqlDB.WithTransaction(func(tx *sqlx.Tx) (interface{}, error) {
		res, err := tx.Exec(insertQuery,
			data.ID,
			data.Name,
			data.BearerKey,
			data.CreatedAt,
			data.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		if rowsAffected <= 0 {
			return nil, errors.New("insert client fails")
		}

		return nil, nil
	})

	return err
}

func (repo *sqlRepository) update(data *Client) error {
	_, err := repo.mysqlDB.WithTransaction(func(tx *sqlx.Tx) (interface{}, error) {
		res, err := tx.Exec(updateQuery,
			data.Name,
			data.BearerKey,
			data.UpdatedAt,
			data.ID,
		)

		if err != nil {
			return nil, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		if rowsAffected <= 0 {
			return nil, errors.New("update client fails")
		}

		return nil, nil
	})

	return err
}
