package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

func main() {
	logFile, _ := os.Create("log.log")
	lg := slog.New(slog.NewJSONHandler(logFile, nil))

	core := Core{
		sessions: SessionRepo{
			sessionRedisClient: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379", // адрес и порт Redis сервера
				Password: "",               // пароль, если требуется
				DB:       0,                // номер базы данных
			}),
			Connection: true,
		},
		csrfTokens: CsrfRepo{
			csrfRedisClient: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379", // адрес и порт Redis сервера
				Password: "",               // пароль, если требуется
				DB:       1,                // номер базы данных
			}),
			Connection: true,
		},
		users: make(map[string]User),
		collections: map[string]string{
			"new":       "Новинки",
			"action":    "Боевик",
			"comedy":    "Комедия",
			"ru":        "Российский",
			"eu":        "Зарубежный",
			"war":       "Военный",
			"kids":      "Детский",
			"detective": "Детектив",
			"drama":     "Драма",
			"crime":     "Криминал",
			"melodrama": "Мелодрама",
			"horror":    "Ужас",
		},
		lg: lg.With("module", "core"),
	}
	go core.CheckRedisSessionsConnection()
	go core.CheckRedisCsrfConnection()
	api := API{core: &core, lg: lg.With("module", "api")}

	mx := http.NewServeMux()
	mx.HandleFunc("/signup", api.Signup)
	mx.HandleFunc("/signin", api.Signin)
	mx.HandleFunc("/logout", api.LogoutSession)
	mx.HandleFunc("/authcheck", api.AuthAccept)
	mx.HandleFunc("/api/v1/films", api.Films)
	mx.HandleFunc("/api/v1/csrf", api.GetCsrfToken)
	err := http.ListenAndServe(":8080", mx)
	if err != nil {
		api.lg.Error("ListenAndServe error", "err", err.Error())
	}
	select {}
}
