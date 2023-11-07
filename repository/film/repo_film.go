package film

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/configs"

	_ "github.com/jackc/pgx/stdlib"
)

type IFilmsRepo interface {
	GetFilmsByGenre(genre uint64, start uint64, end uint64) ([]FilmItem, error)
	GetFilms(start uint64, end uint64) ([]FilmItem, error)
	GetFilm(filmId uint64) (*FilmItem, error)
}

type RepoPostgre struct {
	db *sql.DB
}

func GetFilmRepo(config configs.DbDsnCfg, lg *slog.Logger) *RepoPostgre {
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
		lg.Error("Repo Film db ping error", "err", err.Error())
	}

	time.Sleep(time.Duration(timer) * time.Second)
}

func (repo *RepoPostgre) GetFilmsByGenre(genre uint64, start uint64, end uint64) ([]FilmItem, error) {
	films := make([]FilmItem, 0, end-start)

	rows, err := repo.db.Query(
		"SELECT film.id, film.title, poster FROM film "+
			"JOIN films_genre ON film.id = films_genre.id_film "+
			"WHERE id_genre = $1 "+
			"ORDER BY release_date DESC "+
			"OFFSET $2 LIMIT $3",
		genre, start, end)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("GetFilmsByGenre err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := FilmItem{}
		err := rows.Scan(&post.Id, &post.Title, &post.Poster)
		if err != nil {
			return nil, fmt.Errorf("GetFilmsByGenre scan err: %w", err)
		}
		films = append(films, post)
	}

	return films, nil
}

func (repo *RepoPostgre) GetFilms(start uint64, end uint64) ([]FilmItem, error) {
	films := make([]FilmItem, 0, end-start)

	rows, err := repo.db.Query(
		"SELECT film.id, film.title, poster FROM film "+
			"ORDER BY release_date DESC "+
			"OFFSET $1 LIMIT $2",
		start, end)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("GetFilms err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := FilmItem{}
		err := rows.Scan(&post.Id, &post.Title, &post.Poster)
		if err != nil {
			return nil, fmt.Errorf("GetFilms scan err: %w", err)
		}
		films = append(films, post)
	}

	return films, nil
}

func (repo *RepoPostgre) GetFilm(filmId uint64) (*FilmItem, error) {
	film := &FilmItem{}
	err := repo.db.QueryRow(
		"SELECT * FROM film "+
			"WHERE id = $1", filmId).
		Scan(&film.Id, &film.Title, &film.Info, &film.Poster, &film.ReleaseDate, &film.Country, &film.Mpaa)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return film, nil
		}

		return nil, fmt.Errorf("GetFilm err: %w", err)
	}

	return film, nil
}
