package user_detail

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
	selectCountUserQuery = "SELECT COUNT(*) FROM user_details"
	insertUserQuery      = `
		INSERT INTO user_details (id, device_token, device_os, avatar_url, avatar_path, source, oauth_id, id_card_url, id_card_path, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	updateUserQuery = `
		UPDATE user_details SET
			device_token = ?,
			device_os = ?,
			avatar_url = ?,
			avatar_path = ?,
			source = ?,
			oauth_id = ?,
			id_card_url = ?,
			id_card_path = ?,
			updated_at = ?
		WHERE id = ?
	`
)

func NewSQLRepository(mysqlDB *database.MySQL) Repository {
	return &sqlRepository{mysqlDB}
}

func (repo *sqlRepository) Store(data *UserDetail) error {
	exist, err := repo.existByID(data.ID)
	if err != nil {
		return err
	}

	if exist {
		return repo.update(data)
	}

	data = New(
		data.ID,
		data.DeviceToken,
		data.DeviceOS,
		data.AvatarUrl,
		data.AvatarPath,
		data.Source,
		data.OAuthId,
		data.IDCardUrl,
		data.IDCardPath,
	)

	return repo.insert(data)
}

func (repo *sqlRepository) FindByID(id uuid.UUID) (result *UserDetail, err error) {
	var data UserDetail
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("user_details")
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

func (repo *sqlRepository) FindByFilter(filter Filter) (result []*UserDetail, err error) {
	if filter.IsEmpty() {
		return
	}

	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("user_details")
	if len(filter.IDs) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"id": filter.IDs,
		})
	}

	if len(filter.OAuthIDs) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"oauth_id": filter.OAuthIDs,
		})
	}

	if len(filter.DeviceTokens) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"device_token": filter.DeviceTokens,
		})
	}

	if len(filter.Sources) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"source": filter.Sources,
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

func (repo *sqlRepository) FindPage(offset int, limit int) (result []*UserDetail, err error) {
	dialect := goqu.Dialect("mysql")
	dataset := dialect.From("user_details")
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
	dataset := dialect.From("user_details")
	dataset = dataset.Select(goqu.COUNT("*"))
	if len(filter.IDs) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"id": filter.IDs,
		})
	}

	if len(filter.OAuthIDs) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"oauth_id": filter.OAuthIDs,
		})
	}

	if len(filter.DeviceTokens) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"device_token": filter.DeviceTokens,
		})
	}

	if len(filter.Sources) != 0 {
		dataset = dataset.Where(goqu.Ex{
			"source": filter.Sources,
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

func (repo *sqlRepository) insert(data *UserDetail) error {
	_, err := repo.mysqlDB.WithTransaction(func(tx *sqlx.Tx) (interface{}, error) {
		res, err := tx.Exec(insertUserQuery,
			data.ID,
			data.DeviceToken,
			data.DeviceOS,
			data.AvatarUrl,
			data.AvatarPath,
			data.Source,
			data.OAuthId,
			data.IDCardUrl,
			data.IDCardPath,
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
			return nil, errors.New("insert user detail fails")
		}

		return nil, nil
	})

	return err
}

func (repo *sqlRepository) update(data *UserDetail) error {
	_, err := repo.mysqlDB.WithTransaction(func(tx *sqlx.Tx) (interface{}, error) {
		res, err := tx.Exec(updateUserQuery,
			data.DeviceToken,
			data.DeviceOS,
			data.AvatarUrl,
			data.AvatarPath,
			data.Source,
			data.OAuthId,
			data.IDCardUrl,
			data.IDCardPath,
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
			return nil, errors.New("update user detail fails")
		}

		return nil, nil
	})

	return err
}
