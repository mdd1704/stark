package database

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/palantir/stacktrace"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"stark/utils/activity"
	"stark/utils/log"
)

func NewMongo() (*mongo.Client, error) {
	ctx := activity.NewContext("init_mongo")
	ctx = activity.WithClientID(ctx, "stark_system")

	uri := "mongodb://" + os.Getenv("MONGO_USERNAME") + ":" + os.Getenv("MONGO_PASSWORD") + "@" + os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't connect mongo db"))
		return nil, stacktrace.Propagate(err, "can't connect mongo db")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't ping mongo db"))
		return nil, stacktrace.Propagate(err, "can't ping mongo db")
	}

	driver, err := mongodb.WithInstance(client, &mongodb.Config{
		DatabaseName: os.Getenv("MONGO_DATABASE"),
	})

	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "error create instance mongo"))
		return nil, stacktrace.Propagate(err, "error create instance mongo")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migrations/mongodb/",
		os.Getenv("MONGO_DATABASE"),
		driver,
	)

	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "error init migration mongo"))
		return nil, stacktrace.Propagate(err, "error init migration mongo")
	}

	err = m.Steps(1)
	if err != nil {
		if err.Error() != "file does not exist" {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "error step migration mongo"))
			return nil, stacktrace.Propagate(err, "error step migration mongo")
		}
	}

	return client, nil
}
