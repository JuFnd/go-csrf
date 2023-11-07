package delivery

import (
	"fmt"
	"log/slog"
	"math/rand"
	"regexp"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/configs"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/errors"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/comment"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/crew"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/film"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/csrf"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/genre"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/profession"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/profile"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/session"
)

type Core struct {
	sessions   session.SessionRepo
	csrfTokens csrf.CsrfRepo
	mutex      sync.RWMutex
	lg         *slog.Logger
	Films      film.IFilmsRepo
	Users      profile.IUserRepo
	Genres     genre.IGenreRepo
	Comments   comment.ICommentRepo
	Crew       crew.ICrewRepo
	Profession profession.IProfessionRepo
}

type Session struct {
	Login     string
	ExpiresAt time.Time
}

func GetCore(cfg configs.DbDsnCfg, lg *slog.Logger) *Core {
	core := Core{
		sessions:   make(map[string]Session),
		lg:         lg.With("module", "core"),
		films:      film.GetFilmRepo(cfg, lg),
		users:      profile.GetUserRepo(cfg, lg),
		genres:     genre.GetGenreRepo(cfg, lg),
		comments:   comment.GetCommentRepo(cfg, lg),
		crew:       crew.GetCrewRepo(cfg, lg),
		profession: profession.GetProfessionRepo(cfg, lg),
	}
	return &core
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (core *Core) CreateSession(login string) (string, Session, error) {
	SID := RandStringRunes(32)

	session := Session{
		Login:     login,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	core.mutex.Lock()
	core.sessions[SID] = session
	core.mutex.Unlock()

	return SID, session, nil
}

func (core *Core) KillSession(sid string) error {
	core.mutex.Lock()
	delete(core.sessions, sid)
	core.mutex.Unlock()
	return nil
}

func (core *Core) FindActiveSession(sid string) (bool, error) {
	core.mutex.RLock()
	_, found := core.sessions[sid]
	core.mutex.RUnlock()
	return found, nil
}

func (core *Core) CreateUserAccount(login string, password string, name string, birthDate string, email string) error {
	if matched, _ := regexp.MatchString(`^\w@\w$`, email); !matched {
		return errors.InvalideEmail
	}
	err := core.users.CreateUser(login, password, name, birthDate, email)
	if err != nil {
		core.lg.Error("create user error", "err", err.Error())
		return fmt.Errorf("CreateUserAccount err: %w", err)
	}

	return nil
}

func (core *Core) FindUserAccount(login string, password string) (*profile.UserItem, bool, error) {
	user, found, err := core.users.GetUser(login, password)
	if err != nil {
		core.lg.Error("find user error", "err", err.Error())
		return nil, false, fmt.Errorf("FindUserAccount err: %w", err)
	}
	return user, found, nil
}

func (core *Core) FindUserByLogin(login string) (bool, error) {
	found, err := core.users.FindUser(login)
	if err != nil {
		core.lg.Error("find user error", "err", err.Error())
		return false, fmt.Errorf("FindUserByLogin err: %w", err)
	}

	return found, nil
}

func RandStringRunes(seed int) string {
	symbols := make([]rune, seed)
	for i := range symbols {
		symbols[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(symbols)
}

func (core *Core) GetFilmsByGenre(genre uint64, start uint64, end uint64) ([]film.FilmItem, error) {
	films, err := core.films.GetFilmsByGenre(genre, start, end)
	if err != nil {
		core.lg.Error("failed to get films from db", "err", err.Error())
		return nil, fmt.Errorf("GetFilmsByGenre err: %w", err)
	}

	return films, nil
}

func (core *Core) GetFilms(start uint64, end uint64) ([]film.FilmItem, error) {
	films, err := core.films.GetFilms(start, end)
	if err != nil {
		core.lg.Error("failed to get films from db", "err", err.Error())
		return nil, fmt.Errorf("GetFilms err: %w", err)
	}

	return films, nil
}

func (core *Core) GetFilm(filmId uint64) (*film.FilmItem, error) {
	film, err := core.films.GetFilm(filmId)
	if err != nil {
		core.lg.Error("Get Film error", "err", err.Error())
		return nil, fmt.Errorf("GetFilm err: %w", err)
	}

	return film, nil
}

func (core *Core) GetFilmGenres(filmId uint64) ([]genre.GenreItem, error) {
	genres, err := core.genres.GetFilmGenres(filmId)
	if err != nil {
		core.lg.Error("Get Film Genres error", "err", err.Error())
		return nil, fmt.Errorf("GetFilmGenres err: %w", err)
	}

	return genres, nil
}

func (core *Core) GetFilmRating(filmId uint64) (float64, uint64, error) {
	rating, number, err := core.comments.GetFilmRating(filmId)
	if err != nil {
		core.lg.Error("Get Film Rating error", "err", err.Error())
		return 0, 0, fmt.Errorf("GetFilmRating err: %w", err)
	}

	return rating, number, nil
}

func (core *Core) GetFilmDirectors(filmId uint64) ([]crew.CrewItem, error) {
	directors, err := core.crew.GetFilmDirectors(filmId)
	if err != nil {
		core.lg.Error("Get Film Directors error", "err", err.Error())
		return nil, fmt.Errorf("GetFilmDirectors err: %w", err)
	}

	return directors, nil
}

func (core *Core) GetFilmScenarists(filmId uint64) ([]crew.CrewItem, error) {
	scenarists, err := core.crew.GetFilmScenarists(filmId)
	if err != nil {
		core.lg.Error("Get Film Scenarists error", "err", err.Error())
		return nil, fmt.Errorf("GetFilmScenarists err: %w", err)
	}

	return scenarists, nil
}

func (core *Core) GetFilmCharacters(filmId uint64) ([]crew.Character, error) {
	characters, err := core.crew.GetFilmCharacters(filmId)
	if err != nil {
		core.lg.Error("Get Film Characters error", "err", err.Error())
		return nil, fmt.Errorf("GetFilmCharacters err: %w", err)
	}

	return characters, nil
}

func (core *Core) GetFilmComments(filmId uint64, first uint64, limit uint64) ([]comment.CommentItem, error) {
	comments, err := core.comments.GetFilmComments(filmId, first, limit)
	if err != nil {
		core.lg.Error("Get Film Comments error", "err", err.Error())
		return nil, fmt.Errorf("GetFilmComments err: %w", err)
	}

	return comments, nil
}

func (core *Core) GetActor(actorId uint64) (*crew.CrewItem, error) {
	actor, err := core.crew.GetActor(actorId)
	if err != nil {
		core.lg.Error("Get Actor error", "err", err.Error())
		return nil, fmt.Errorf("GetActor err: %w", err)
	}

	return actor, nil
}

func (core *Core) GetActorsCareer(actorId uint64) ([]profession.ProfessionItem, error) {
	career, err := core.profession.GetActorsProfessions(actorId)
	if err != nil {
		core.lg.Error("Get Actors Career error", "err", err.Error())
		return nil, fmt.Errorf("GetActorsCareer err: %w", err)
	}

	return career, nil
}

func (core *Core) AddComment(filmId uint64, userLogin string, rating uint16, text string) error {
	err := core.comments.AddComment(filmId, userLogin, rating, text)
	if err != nil {
		core.lg.Error("Add Comment error", "err", err.Error())
		return fmt.Errorf("GetActorsCareer err: %w", err)
	}

	return nil
}

func (core *Core) GetUserName(sid string) (string, error) {
	return "", nil
}

func (core *Core) GetUserProfile(login string) (*profile.UserItem, error) {
	profile, err := core.users.GetUserProfile(login)
	if err != nil {
		core.lg.Error("GetUserProfile error", "err", err.Error())
		return nil, fmt.Errorf("GetUserProfile err: %w", err)
	}

	return profile, nil
}

func (core *Core) GetGenre(genreId uint64) (string, error) {
	genre, err := core.genres.GetGenreById(genreId)
	if err != nil {
		core.lg.Error("GetGenre error", "err", err.Error())
		return "", fmt.Errorf("GetGenre err: %w", err)
	}

	return genre, nil
}
