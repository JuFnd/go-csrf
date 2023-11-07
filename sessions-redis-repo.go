package main

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

type Session struct {
	Login     string
	SID       string
	ExpiresAt time.Time
}

func (redisRepo *SessionRepo) CheckRedisSessionConnection() {
	ctx := context.Background()
	for {

		_, err := redisRepo.sessionRedisClient.Ping(ctx).Result()
		redisRepo.Connection = err == nil

		time.Sleep(15 * time.Second)
	}
}

func GetSessionRepo(lg *slog.Logger) *SessionRepo {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	_, err := redisClient.Ping(ctx).Result()

	sessionRepo := SessionRepo{
		sessionRedisClient: redisClient,
		Connection:         err == nil,
	}

	go sessionRepo.CheckRedisSessionConnection()

	return &sessionRepo

}

func (redisRepo *SessionRepo) AddSession(active Session, lg *slog.Logger) bool {
	if !redisRepo.Connection {
		lg.Error("Redis session connection lost")
		return false
	}

	ctx := context.Background()
	err := redisRepo.sessionRedisClient.Set(ctx, active.SID, active.Login, 24*time.Hour)

	if err != nil {
		lg.Error("Error, cannot create session " + active.SID)
		return false
	}

	return redisRepo.CheckActiveSession(active.SID, lg)
}

func (redisRepo *SessionRepo) CheckActiveSession(sid string, lg *slog.Logger) bool {
	if !redisRepo.Connection {
		lg.Error("Redis session connection lost")
		return false
	}

	ctx := context.Background()

	_, err := redisRepo.sessionRedisClient.Get(ctx, sid).Result()
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

func (redisRepo *SessionRepo) DeleteSession(sid string, lg *slog.Logger) bool {
	ctx := context.Background()

	_, err := redisRepo.sessionRedisClient.Del(ctx, sid).Result()

	if err != nil {
		lg.Error("Delete request could not be completed:", err)
	}

	return err != nil
}
