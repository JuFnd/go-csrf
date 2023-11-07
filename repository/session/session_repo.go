package session

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
)

type SessionRepo struct {
	sessionRedisClient *redis.Client
	Connection         bool
}

func (redisRepo *SessionRepo) CheckRedisSessionConnection() {
	ctx := context.Background()
	for {
		_, err := redisRepo.sessionRedisClient.Ping(ctx).Result()
		redisRepo.Connection = err == nil

		time.Sleep(15 * time.Second)
	}
}

func GetSessionRepo(lg *slog.Logger) (*SessionRepo, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	sessionRepo := SessionRepo{
		sessionRedisClient: redisClient,
		Connection:         true,
	}

	go sessionRepo.CheckRedisSessionConnection()

	return &sessionRepo, nil
}

func (redisRepo *SessionRepo) AddSession(active Session, lg *slog.Logger) (bool, error) {
	if !redisRepo.Connection {
		lg.Error("Redis session connection lost")
		return false, nil
	}

	ctx := context.Background()
	err := redisRepo.sessionRedisClient.Set(ctx, active.SID, active.Login, 24*time.Hour)
	if err != nil {
		lg.Error("Error, cannot create session " + active.SID)
		return false, err.Err()
	}

	sessionAdded, err_check := redisRepo.CheckActiveSession(active.SID, lg)

	if err_check != nil {
		return false, err_check
	}

	return sessionAdded, nil
}

func (redisRepo *SessionRepo) CheckActiveSession(sid string, lg *slog.Logger) (bool, error) {
	if !redisRepo.Connection {
		lg.Error("Redis session connection lost")
		return false, nil
	}

	ctx := context.Background()

	_, err := redisRepo.sessionRedisClient.Get(ctx, sid).Result()
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

func (redisRepo *SessionRepo) DeleteSession(sid string, lg *slog.Logger) (bool, error) {
	ctx := context.Background()

	_, err := redisRepo.sessionRedisClient.Del(ctx, sid).Result()
	if err != nil {
		lg.Error("Delete request could not be completed:", err)
		return false, err
	}

	return true, nil
}
