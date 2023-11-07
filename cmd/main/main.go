package main

import (
	"log/slog"
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

	api.ListenAndServe()
}
