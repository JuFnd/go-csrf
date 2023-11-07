package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logFile, _ := os.Create("log.log")
	lg := slog.New(slog.NewJSONHandler(logFile, nil))

	core := Core{
		lg:         lg.With("module", "core"),
		sessions:   *GetSessionRepo(lg),
		csrfTokens: *GetCsrfRepo(lg),
		users:      make(map[string]User),
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
	}

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
