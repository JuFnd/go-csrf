package profile

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/configs"

	_ "github.com/jackc/pgx/stdlib"
)

type IUserRepo interface {
	GetUser(login string, password string) (*UserItem, bool, error)
	FindUser(login string) (bool, error)
	CreateUser(login string, password string, name string, birthDate string, email string) error
	GetUserProfile(login string) (*UserItem, error)
}

type RepoPostgre struct {
	db *sql.DB
}

func GetUserRepo(config configs.DbDsnCfg, lg *slog.Logger) *RepoPostgre {
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
		lg.Error("Repo Profile db ping error", "err", err.Error())
	}

	time.Sleep(time.Duration(timer) * time.Second)
}

func (repo *RepoPostgre) GetUser(login string, password string) (*UserItem, bool, error) {
	post := &UserItem{}

	err := repo.db.QueryRow(
		"SELECT login, photo FROM profile "+
			"WHERE login = $1 AND password = $2", login, password).Scan(&post.Login, &post.Photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("GetUser err: %w", err)
	}

	return post, true, nil
}

func (repo *RepoPostgre) FindUser(login string) (bool, error) {
	post := &UserItem{}

	err := repo.db.QueryRow(
		"SELECT login FROM profile "+
			"WHERE login = $1", login).Scan(&post.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("FindUser err: %w", err)
	}

	return true, nil
}

func (repo *RepoPostgre) CreateUser(login string, password string, name string, birthDate string, email string) error {
	_, err := repo.db.Exec(
		"INSERT INTO profile(name, birth_date, photo, login, password, email, registration_date) "+
			"VALUES($1, $2, '../../user_avatars/default.jpg', $3, $4, $5, CURRENT_TIMESTAMP)",
		name, birthDate, login, password, email)
	if err != nil {
		return fmt.Errorf("CreateUser err: %w", err)
	}

	return nil
}

func (repo *RepoPostgre) GetUserProfile(login string) (*UserItem, error) {
	post := &UserItem{}

	err := repo.db.QueryRow(
		"SELECT name, birth_date, login, email, photo FROM profile "+
			"WHERE login = $1", login).Scan(&post.Name, &post.Birthdate, &post.Login, &post.Email, &post.Photo)
	if err != nil {
		return nil, fmt.Errorf("GetUserProfile err: %w", err)
	}

	return post, nil
}
