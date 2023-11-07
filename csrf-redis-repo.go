package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
)

type CsrfRepo struct {
	csrfRedisClient *redis.Client
	Connection      bool
}

type Csrf struct {
	SID       string
	ExpiresAt time.Time
}

func (redisRepo *CsrfRepo) CheckRedisCsrfConnection() {
	ctx := context.Background()
	for {

		_, err := redisRepo.csrfRedisClient.Ping(ctx).Result()
		redisRepo.Connection = err == nil

		time.Sleep(15 * time.Second)
	}
}

func GetCsrfRepo(lg *slog.Logger) *CsrfRepo {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	ctx := context.Background()

	_, err := redisClient.Ping(ctx).Result()

	csrfRepo := CsrfRepo{
		csrfRedisClient: redisClient,
		Connection:      err == nil,
	}

	go csrfRepo.CheckRedisCsrfConnection()

	return &csrfRepo

}

func (redisRepo *CsrfRepo) AddCsrf(active Csrf, lg *slog.Logger) bool {
	if !redisRepo.Connection {
		lg.Error("Redis csrf connection lost")
		return false
	}

	ctx := context.Background()
	err := redisRepo.csrfRedisClient.Set(ctx, active.SID, active.SID, 3*time.Hour)

	if err != nil {
		lg.Error("Error, cannot create csrf token ", err.Err())
		return false
	}

	return redisRepo.CheckActiveCsrf(active.SID, lg)
}

func (redisRepo *CsrfRepo) CheckActiveCsrf(sid string, lg *slog.Logger) bool {
	if !redisRepo.Connection {
		lg.Error("Redis csrf connection lost")
		return false
	}

	ctx := context.Background()

	_, err := redisRepo.csrfRedisClient.Get(ctx, sid).Result()
	if err == redis.Nil {
		lg.Error("Key" + sid + "not found")
		return false
	}

	if err != nil {
		lg.Error("Get request could not be completed ", err)
		return false
	}

	return true
}

func (redisRepo *CsrfRepo) DeleteSession(sid string, lg *slog.Logger) bool {
	ctx := context.Background()

	_, err := redisRepo.csrfRedisClient.Del(ctx, sid).Result()

	if err != nil {
		lg.Error("Delete request could not be completed:", err)
	}

	return err != nil
}
