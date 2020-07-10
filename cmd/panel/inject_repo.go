package main

import (
	"github.com/go-redis/redis"
	"github.com/h3isenbug/url-panel/config"
	"github.com/h3isenbug/url-panel/repositories/url"
	"github.com/h3isenbug/url-panel/repositories/user"
	"github.com/jmoiron/sqlx"
)

func provideUserRepository() (user.UserRepository, func(), error) {
	var con, err = sqlx.Open("postgres", config.Config.UserDatabaseDSN)
	if err != nil {
		return nil, nil, err
	}
	repo, err := user.NewPostgresUserRepository(con)
	return repo, func() { con.Close() }, err
}

func provideURLRepository(redisClient *redis.Client) (url.URLRepository, func(), error) {
	var con, err = sqlx.Open("postgres", config.Config.URLDatabaseDSN)
	if err != nil {
		return nil, nil, err
	}
	pgRepo, err := url.NewPostgresURLRepository(con)

	redisRepo := url.NewRedisCacheWriteRepositoryV1(redisClient, pgRepo)
	return redisRepo, func() { con.Close() }, err
}

func provideRedisClient() (*redis.Client, func()) {
	var client = redis.NewClient(&redis.Options{
		Addr:     config.Config.URLRedisServer,
		Password: config.Config.URLRedisPassword,
		DB:       config.Config.URLRedisDB,
	})
	return client, func() { client.Close() }
}
