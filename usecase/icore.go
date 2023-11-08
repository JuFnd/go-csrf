package usecase

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/comment"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/crew"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/film"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/genre"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/profession"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/profile"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/session"
)

type ICore interface {
	CreateSession(login string) (string, session.Session, error)
	KillSession(sid string) error
	FindActiveSession(sid string) (bool, error)
	CreateUserAccount(login string, password string, name string, birthDate string, email string) error
	FindUserAccount(login string, password string) (*profile.UserItem, bool, error)
	FindUserByLogin(login string) (bool, error)
	GetFilmsByGenre(genre uint64, start uint64, end uint64) ([]film.FilmItem, error)
	GetFilms(start uint64, end uint64) ([]film.FilmItem, error)
	GetFilm(filmId uint64) (*film.FilmItem, error)
	GetFilmGenres(filmId uint64) ([]genre.GenreItem, error)
	GetFilmRating(filmId uint64) (float64, uint64, error)
	GetFilmDirectors(filmId uint64) ([]crew.CrewItem, error)
	GetFilmScenarists(filmId uint64) ([]crew.CrewItem, error)
	GetFilmCharacters(filmId uint64) ([]crew.Character, error)
	GetFilmComments(filmId uint64, first uint64, limit uint64) ([]comment.CommentItem, error)
	GetActor(actorId uint64) (*crew.CrewItem, error)
	GetActorsCareer(actorId uint64) ([]profession.ProfessionItem, error)
	AddComment(filmId uint64, userLogin string, rating uint16, text string) error
	GetUserName(sid string) (string, error)
	GetUserProfile(login string) (*profile.UserItem, error)
	GetGenre(genreId uint64) (string, error)
	CheckCsrfToken(token string, r *http.Request) (bool, error)
	CreateCsrfToken(r *http.Request) (string, error)
}
