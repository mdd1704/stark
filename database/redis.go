package database

import (
	"os"
	"stark/utils/activity"
	"stark/utils/log"
	"time"

	"github.com/go-redis/redis"
	"github.com/palantir/stacktrace"
)

type Redis struct {
	db *redis.Client
}

func NewRedis() (*Redis, error) {
	ctx := activity.NewContext("init_redis")
	ctx = activity.WithClientID(ctx, "stark_system")

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		// Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't ping redis db"))
		return nil, stacktrace.Propagate(err, "can't ping redis db")
	}

	rdb := &Redis{
		db: client,
	}

	return rdb, nil
}

func (r *Redis) Set(key string, value string, sub time.Duration) error {
	err := r.db.Set(key, value, sub).Err()
	if err != nil {
		return stacktrace.Propagate(err, "can't set to redis db")
	}

	return nil
}

func (r *Redis) Get(key string) (string, error) {
	key, err := r.db.Get(key).Result()
	if err != nil {
		return "", stacktrace.Propagate(err, "can't get from redis DB")
	} else {
		return key, nil
	}
}

func (r *Redis) Delete(key string) (int64, error) {
	deleted, err := r.db.Del(key).Result()
	if err != nil {
		return 0, stacktrace.Propagate(err, "can't delete key redis DB")
	}

	return deleted, nil
}
