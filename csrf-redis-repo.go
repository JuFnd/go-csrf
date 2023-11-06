package main

import (
	"context"
	"fmt"
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

func (redisRepo *CsrfRepo) AddCsrf(active Csrf) bool {
	if !redisRepo.Connection {
		fmt.Println("Redis connection lost")
		return false
	}

	ctx := context.Background()
	redisRepo.csrfRedisClient.Set(ctx, active.SID, active.SID, 3*time.Hour)
	if !redisRepo.CheckActiveCsrf(active.SID) {
		fmt.Println("Error create session", active.SID)
		return false
	}
	return true
}

func (redisRepo *CsrfRepo) CheckActiveCsrf(sid string) bool {
	if !redisRepo.Connection {
		fmt.Println("Redis connection lost")
		return false
	}

	ctx := context.Background()

	session, err := redisRepo.csrfRedisClient.Get(ctx, sid).Result()
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

func (redisRepo *CsrfRepo) DeleteSession(sid string) bool {
	// Создание контекста
	ctx := context.Background()

	// Выполнение запроса на удаление
	result, err := redisRepo.csrfRedisClient.Del(ctx, sid).Result()
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса на удаление:", err)
	} else {
		fmt.Println("Количество удаленных ключей:", result)
	}
	return err != nil
}
