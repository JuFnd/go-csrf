package csrf

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/configs"
	"github.com/go-redis/redis/v8"
)

type CsrfRepo struct {
	csrfRedisClient *redis.Client
	mutex           sync.RWMutex
	Connection      bool
}

func (redisRepo *CsrfRepo) CheckRedisCsrfConnection(csrfCfg configs.DbRedisCfg) {
	ctx := context.Background()
	for {
		_, err := redisRepo.csrfRedisClient.Ping(ctx).Result()
		redisRepo.mutex.Lock()
		redisRepo.Connection = err == nil
		redisRepo.mutex.Unlock()

		time.Sleep(time.Duration(csrfCfg.Timer) * time.Second)
	}
}

func GetCsrfRepo(csrfConfigs configs.DbRedisCfg, lg *slog.Logger) (*CsrfRepo, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     csrfConfigs.Host,
		Password: csrfConfigs.Password,
		DB:       csrfConfigs.DbNumber,
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

	go csrfRepo.CheckRedisCsrfConnection(csrfConfigs)

	return &csrfRepo, nil
}

func (redisRepo *CsrfRepo) AddCsrf(active Csrf, lg *slog.Logger, r *http.Request) (bool, error) {
	if !redisRepo.Connection {
		lg.Error("Redis csrf connection lost")
		return false, nil
	}

	redisRepo.csrfRedisClient.Set(r.Context(), active.SID, active.SID, 3*time.Hour)

	csrfAdded, err_check := redisRepo.CheckActiveCsrf(active.SID, lg, r)

	if err_check != nil {
		lg.Error("Error, cannot create csrf token " + err_check.Error())
		return false, err_check
	}

	return csrfAdded, nil
}

func (redisRepo *CsrfRepo) CheckActiveCsrf(sid string, lg *slog.Logger, r *http.Request) (bool, error) {
	if !redisRepo.Connection {
		lg.Error("Redis csrf connection lost")
		return false, nil
	}

	_, err := redisRepo.csrfRedisClient.Get(r.Context(), sid).Result()
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

func (redisRepo *CsrfRepo) DeleteSession(sid string, lg *slog.Logger, r *http.Request) (bool, error) {
	ctx := context.Background()

	_, err := redisRepo.csrfRedisClient.Del(ctx, sid).Result()
	if err != nil {
		lg.Error("Delete request could not be completed:", err)
		return false, err
	}

	return true, nil
}
