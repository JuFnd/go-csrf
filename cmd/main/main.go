package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/configs"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/delivery"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/usecase"
)

func main() {
	logFile, _ := os.Create("log.log")
	lg := slog.New(slog.NewJSONHandler(logFile, nil))

	config, err := configs.ReadConfig()
	if err != nil {
		lg.Error("read config error", "err", err.Error())
		return
	}
	core := delivery.GetCore(*config, lg)
	api := usecase.GetApi(core, lg)

	mx := http.NewServeMux()
	mx.HandleFunc("/signup", api.Signup)
	mx.HandleFunc("/signin", api.Signin)
	mx.HandleFunc("/logout", api.LogoutSession)
	mx.HandleFunc("/authcheck", api.AuthAccept)
	mx.HandleFunc("/api/v1/films", api.Films)
	mx.HandleFunc("/api/v1/film", api.Film)
	mx.HandleFunc("/api/v1/actor", api.Actor)
	mx.HandleFunc("/api/v1/comment", api.Comment)
	mx.HandleFunc("/api/v1/comment/add", api.AddComment)
	mx.HandleFunc("/api/v1/csrf", api.GetCsrfToken)
	err = http.ListenAndServe(":8080", mx)
	if err != nil {
		lg.Error("ListenAndServe error", "err", err.Error())
	}
}
