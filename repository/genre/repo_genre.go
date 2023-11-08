package genre

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/configs"

	_ "github.com/jackc/pgx/stdlib"
)

type IGenreRepo interface {
	GetFilmGenres(filmId uint64) ([]GenreItem, error)
	GetGenreById(genreId uint64) (string, error)
}

type RepoPostgre struct {
	db *sql.DB
}

func GetGenreRepo(config configs.DbDsnCfg, lg *slog.Logger) *RepoPostgre {
	dsn := fmt.Sprintf("user=%s dbname=%s password= %s host=%s port=%d sslmode=%s",
		config.User, config.DbName, config.Password, config.Host, config.Port, config.Sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		lg.Error("sql open error", "err", err.Error())
		return nil
	}
	err = db.Ping()
	if err != nil {
		lg.Error("sql ping error", "err", err.Error())
		return nil
	}
	db.SetMaxOpenConns(config.MaxOpenConns)

	postgreDb := RepoPostgre{db: db}

	go postgreDb.pingDb(config.Timer, lg)
	return &postgreDb
}

func (repo *RepoPostgre) pingDb(timer uint32, lg *slog.Logger) {
	err := repo.db.Ping()
	if err != nil {
		lg.Error("Repo Genre db ping error", "err", err.Error())
	}

	time.Sleep(time.Duration(timer) * time.Second)
}

func (repo *RepoPostgre) GetFilmGenres(filmId uint64) ([]GenreItem, error) {
	genres := []GenreItem{}

	rows, err := repo.db.Query(
		"SELECT genre.id, genre.title FROM genre "+
			"JOIN films_genre ON genre.id = films_genre.id_genre "+
			"WHERE films_genre.id_film = $1", filmId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("GetFilmGenres err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := GenreItem{}
		err := rows.Scan(&post.Id, &post.Title)
		if err != nil {
			return nil, fmt.Errorf("GetFilmGenres scan err: %w", err)
		}
		genres = append(genres, post)
	}

	return genres, nil
}

func (repo *RepoPostgre) GetGenreById(genreId uint64) (string, error) {
	var genre string

	err := repo.db.QueryRow(
		"SELECT title FROM genre "+
			"WHERE id = $1", genreId).Scan(&genre)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return genre, nil
}
