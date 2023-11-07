package comment

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/configs"

	_ "github.com/jackc/pgx/stdlib"
)

type ICommentRepo interface {
	GetFilmRating(filmId uint64) (float64, uint64, error)
	GetFilmComments(filmId uint64, first uint64, limit uint64) ([]CommentItem, error)
	AddComment(filmId uint64, userId uint64, rating uint16, text string) error
}

type RepoPostgre struct {
	DB *sql.DB
}

func GetCommentRepo(config configs.DbDsnCfg, lg *slog.Logger) *RepoPostgre {
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

	postgreDb := RepoPostgre{DB: db}

	go postgreDb.pingDb(config.Timer, lg)
	return &postgreDb
}

func (repo *RepoPostgre) pingDb(timer uint32, lg *slog.Logger) {
	err := repo.DB.Ping()
	if err != nil {
		lg.Error("Repo Comment db ping error", "err", err.Error())
	}

	time.Sleep(time.Duration(timer) * time.Second)
}

func (repo *RepoPostgre) GetFilmRating(filmId uint64) (float64, uint64, error) {
	var rating float64
	var number uint64
	err := repo.DB.QueryRow(
		"SELECT AVG(rating), COUNT(rating) FROM users_comment "+
			"WHERE id_film = $1", filmId).Scan(&rating, &number)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}

	return rating, number, nil
}

func (repo *RepoPostgre) GetFilmComments(filmId uint64, first uint64, limit uint64) ([]CommentItem, error) {
	var comments []CommentItem

	rows, err := repo.DB.Query(
		"SELECT profile.login, rating, comment FROM users_comment "+
			"JOIN profile ON users_comment.id_user = profile.id "+
			"WHERE id_film = $1 "+
			"OFFSET $2 LIMIT $3", filmId, first, limit)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		post := CommentItem{}
		err := rows.Scan(&post.Username, &post.Rating, &post.Comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, post)
	}

	return comments, nil
}

func (repo *RepoPostgre) AddComment(filmId uint64, userId uint64, rating uint16, text string) error {
	_, err := repo.DB.Exec(
		"INSERT INTO users_comment(id_film, id_user, rating, comment) "+
			"VALUES($1, $2, $3, $4) ", filmId, userId, rating, text)
	if err != nil {
		return err
	}

	return nil
}
