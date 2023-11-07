package csrf

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

func (redisRepo *CsrfRepo) CheckRedisCsrfConnection() {
	ctx := context.Background()
	for {
		_, err := redisRepo.csrfRedisClient.Ping(ctx).Result()
		redisRepo.Connection = err == nil

		time.Sleep(15 * time.Second)
	}
}

func GetCsrfRepo(lg *slog.Logger) (*CsrfRepo, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	ctx := context.Background()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	csrfRepo := CsrfRepo{
		csrfRedisClient: redisClient,
		Connection:      true,
	}

	go csrfRepo.CheckRedisCsrfConnection()

	return &csrfRepo, nil
}

func (redisRepo *CsrfRepo) AddCsrf(active Csrf, lg *slog.Logger) (bool, error) {
	if !redisRepo.Connection {
		lg.Error("Redis csrf connection lost")
		return false, nil
	}

	ctx := context.Background()
	redisRepo.csrfRedisClient.Set(ctx, active.SID, active.SID, 3*time.Hour)

	csrfAdded, err_check := redisRepo.CheckActiveCsrf(active.SID, lg)

	if err_check != nil {
		lg.Error("Error, cannot create csrf token " + err_check.Error())
		return false, err_check
	}

	return csrfAdded, nil
}

func (redisRepo *CsrfRepo) CheckActiveCsrf(sid string, lg *slog.Logger) (bool, error) {
	if !redisRepo.Connection {
		lg.Error("Redis csrf connection lost")
		return false, nil
	}

	ctx := context.Background()

	_, err := redisRepo.csrfRedisClient.Get(ctx, sid).Result()
	if err == redis.Nil {
		lg.Error("Key " + sid + " not found")
		return false, nil
	}

	if err != nil {
		lg.Error("Get request could not be completed ", err)
		return false, err
	}

	return true, nil
}

func (redisRepo *CsrfRepo) DeleteSession(sid string, lg *slog.Logger) (bool, error) {
	ctx := context.Background()

	_, err := redisRepo.csrfRedisClient.Del(ctx, sid).Result()
	if err != nil {
		lg.Error("Delete request could not be completed:", err)
		return false, err
	}

	return true, nil
}
