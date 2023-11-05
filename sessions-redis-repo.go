package main

import (
	"context"
	"fmt"
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

func (redisRepo *SessionRepo) AddSession(active Session) bool {
	if !redisRepo.Connection {
		//fmt.Println("Redis connection lost")
		//return false
	}

	ctx := context.Background()
	redisRepo.sessionRedisClient.Set(ctx, active.SID, active.Login, 24*time.Hour)
	if !redisRepo.CheckActiveSession(active.SID) {
		fmt.Println("Error create session", active.SID)
		return false
	}
	return true
}

func (redisRepo *SessionRepo) CheckActiveSession(sid string) bool {
	if !redisRepo.Connection {
		//fmt.Println("Redis connection lost")
		//return false
	}

	ctx := context.Background()

	session, err := redisRepo.sessionRedisClient.Get(ctx, sid).Result()
	if err == redis.Nil {
		fmt.Println("Ключ не найден")
		return false
	} else if err != nil {
		fmt.Println("Ошибка при выполнении GET-запроса:", err)
		return false
	} else {
		fmt.Println("Значение:", session)
		return true
	}
}

func (redisRepo *SessionRepo) DeleteSession(sid string) bool {
	// Создание контекста
	ctx := context.Background()

	// Выполнение запроса на удаление
	result, err := redisRepo.sessionRedisClient.Del(ctx, sid).Result()
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса на удаление:", err)
	} else {
		fmt.Println("Количество удаленных ключей:", result)
	}
	return err != nil
}
