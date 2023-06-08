package user

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
	selectCountUserQuery = "SELECT COUNT(*) FROM users"
	insertUserQuery      = `
		INSERT INTO users (id, name, email, username, contact, password, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ? ,?)
	`
	updateUserQuery = `
		UPDATE users SET
			name = ?,
			email = ?,
			username = ?,
			contact = ?,
			password = ?,
			updated_at = ?
		WHERE id = ?
	`
	updateProfileQuery = `
		UPDATE users SET
			name = ?,
			username = ?,
			contact = ?,
			updated_at = ?
		WHERE id = ?
	`
)

func NewSQLRepository(mysqlDB *database.MySQL) Repository {
	return &sqlRepository{mysqlDB}
}

func (repo *sqlRepository) Store(data *User) error {
	exist, err := repo.existByID(data.ID)
	if err != nil {
		return err
	}

	if exist {
		return repo.update(data)
	}

	return repo.insert(data)
}

func (repo *sqlRepository) StoreProfile(data *User) error {
	exist, err := repo.existByID(data.ID)
	if err != nil {
		return err
	}

	if exist {
		return repo.updateProfile(data)
	} else {
		return errors.New("user ID not exists")
	}
}

func (repo *sqlRepository) FindByID(id uuid.UUID) (result *User, err error) {
	var data User
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("users")
	dataset = dataset.Where(goqu.Ex{
		"id": id.String(),
	})

	sql, _, err := dataset.ToSQL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "sql error")
	}

	err = repo.mysqlDB.Get(&data, sql)
	if err != nil {
		return nil, stacktrace.Propagate(err, "can't read user by id")
	}

	return &data, nil
}

func (repo *sqlRepository) FindByFilter(filter Filter) (result []*User, err error) {
	if filter.IsEmpty() {
		return
	}

	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("users")
	if len(filter.Emails) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"email": filter.Emails,
		})
	}

	if len(filter.Usernames) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"username": filter.Usernames,
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

func (repo *sqlRepository) FindPage(offset int, limit int) (result []*User, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("users")
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
	dataset := dialect.From("users")
	dataset = dataset.Select(goqu.COUNT("*"))
	if len(filter.Emails) != 0 {
		dataset = dataset.Where(goqu.ExOr{
			"email": filter.Emails,
		})
	}

	if len(filter.Usernames) != 0 {
		dataset = dataset.Where(goqu.ExOr{
			"username": filter.Usernames,
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
	err := repo.mysqlDB.Get(&total, selectCountUserQuery+" WHERE id = ?", id)
	if err != nil {
		return false, stacktrace.Propagate(err, "select count fails")
	}

	return total > 0, nil
}

func (repo *sqlRepository) insert(data *User) error {
	_, err := repo.mysqlDB.WithTransaction(func(tx *sqlx.Tx) (interface{}, error) {
		res, err := tx.Exec(insertUserQuery,
			data.ID,
			data.Name,
			data.Email,
			data.Username,
			data.Contact,
			data.Password,
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
			return nil, errors.New("insert user fails")
		}

		return nil, nil
	})

	return err
}

func (repo *sqlRepository) update(data *User) error {
	_, err := repo.mysqlDB.WithTransaction(func(tx *sqlx.Tx) (interface{}, error) {
		res, err := tx.Exec(updateUserQuery,
			data.Name,
			data.Email,
			data.Username,
			data.Contact,
			data.Password,
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
			return nil, errors.New("update user fails")
		}

		return nil, nil
	})

	return err
}

func (repo *sqlRepository) updateProfile(data *User) error {
	_, err := repo.mysqlDB.WithTransaction(func(tx *sqlx.Tx) (interface{}, error) {
		res, err := tx.Exec(updateProfileQuery,
			data.Name,
			data.Username,
			data.Contact,
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
			return nil, errors.New("update profile fails")
		}

		return nil, nil
	})

	return err
}
