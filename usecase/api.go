package usecase

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/delivery"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/film"
)

type API struct {
	core delivery.ICore
	lg   *slog.Logger
}

func GetApi(c *delivery.Core, l *slog.Logger) *API {
	return &API{core: c, lg: l.With("module", "api")}
}

func (a *API) GetCsrfToken(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Body: nil}

	csrfToken := r.Header.Get("x-csrf-token")

	found, err := a.core.CheckCsrfToken(csrfToken)
	if csrfToken != "" && found {
		w.Header().Set("X-CSRF-Token", csrfToken)
		a.SendResponse(w, response)
		return
	}

	token, err := a.core.CreateCsrfToken()

	if err != nil {
		w.Header().Set("X-CSRF-Token", "null")
		response.Status = 502
		a.SendResponse(w, response)
	}

	w.Header().Set("X-CSRF-Token", token)
	a.SendResponse(w, response)
	return

}

func (a *API) SendResponse(w http.ResponseWriter, response Response) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.lg.Error("failed to pack json", "err", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		a.lg.Error("failed to send response", "err", err.Error())
	}
}

func (a *API) Films(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, response)
		return
	}

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.ParseUint(r.URL.Query().Get("page_size"), 10, 64)
	if err != nil {
		pageSize = 8
	}

	var films []film.FilmItem
	collectionId := r.URL.Query().Get("collection_id")
	if collectionId == "" {
		films, err = a.core.GetFilms(uint64((page-1)*pageSize), pageSize)
	} else {
		films, err = a.core.GetFilmsByGenre(collectionId, uint64((page-1)*pageSize), pageSize)
	}
	if err != nil {
		a.lg.Error("Films error", "err", err.Error())
	}
	filmsResponse := FilmsResponse{
		Page:           page,
		PageSize:       pageSize,
		Total:          uint64(len(films)),
		CollectionName: collectionId,
		Films:          films,
	}
	response.Body = filmsResponse

	a.SendResponse(w, response)
}

func (a *API) LogoutSession(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Body: nil}

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		response.Status = http.StatusUnauthorized
		a.SendResponse(w, response)
		return
	}

	found, _ := a.core.FindActiveSession(session.Value)
	if !found {
		response.Status = http.StatusUnauthorized
		a.SendResponse(w, response)
		return
	} else {
		err := a.core.KillSession(session.Value)
		if err != nil {
			a.lg.Error("failed to kill session", "err", err.Error())
		}
		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	}

	a.SendResponse(w, response)
}

func (a *API) AuthAccept(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Body: nil}
	var authorized bool

	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized, _ = a.core.FindActiveSession(session.Value)
	}

	if !authorized {
		response.Status = http.StatusUnauthorized
		a.SendResponse(w, response)
		return
	}

	a.SendResponse(w, response)
}

func (a *API) Signin(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Body: nil}
	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, response)
		return
	}
	csrfToken := r.Header.Get("x-csrf-token")

	found, err := a.core.CheckCsrfToken(csrfToken)
	if !found || err != nil {
		response.Status = 412
		a.SendResponse(w, response)
		return
	}

	var request SigninRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		a.SendResponse(w, response)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		response.Status = http.StatusBadRequest
		a.SendResponse(w, response)
		return
	}

	user, found, err := a.core.FindUserAccount(request.Login, request.Password)
	if err != nil {
		a.lg.Error("Signin error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}
	if !found {
		response.Status = http.StatusUnauthorized
		a.SendResponse(w, response)
		return
	} else {
		sid, session, _ := a.core.CreateSession(user.Login)
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    sid,
			Path:     "/",
			Expires:  session.ExpiresAt,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}

	a.SendResponse(w, response)
}

func (a *API) Signup(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Body: nil}
	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, response)
		return
	}

	csrfToken := r.Header.Get("x-csrf-token")

	found, err_token := a.core.CheckCsrfToken(csrfToken)
	if !found || err_token != nil {
		response.Status = 412
		a.SendResponse(w, response)
		return
	}

	var request SignupRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		a.SendResponse(w, response)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Status = http.StatusBadRequest
		a.SendResponse(w, response)
		return
	}

	found, err_find := a.core.FindUserByLogin(request.Login)
	if err_find != nil {
		a.lg.Error("Signup error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}
	if found {
		response.Status = http.StatusConflict
		a.SendResponse(w, response)
		return
	}
	err = a.core.CreateUserAccount(request.Login, request.Password, request.Name, request.BirthDate, request.Email)
	if err == delivery.InvalideEmail {
		a.lg.Error("create user error", "err", err.Error())
		response.Status = http.StatusBadRequest
	}
	if err != nil {
		a.lg.Error("failed to create user account", "err", err.Error())
		response.Status = http.StatusBadRequest
	}

	a.SendResponse(w, response)
}

func (a *API) Film(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Body: nil}
	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, response)
		return
	}

	filmId, err := strconv.ParseUint(r.URL.Query().Get("film_id"), 10, 64)
	if err != nil {
		response.Status = http.StatusBadRequest
		a.SendResponse(w, response)
		return
	}

	film, err := a.core.GetFilm(filmId)
	if err != nil {
		a.lg.Error("Film error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}
	if film == nil {
		response.Status = http.StatusNotFound
		a.SendResponse(w, response)
		return
	}
	genres, err := a.core.GetFilmGenres(filmId)
	if err != nil {
		a.lg.Error("Film error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}
	rating, number, err := a.core.GetFilmRating(filmId)
	if err != nil {
		a.lg.Error("Film error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}
	directors, err := a.core.GetFilmDirectors(filmId)
	if err != nil {
		a.lg.Error("Film error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}
	scenarists, err := a.core.GetFilmScenarists(filmId)
	if err != nil {
		a.lg.Error("Film error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}
	characters, err := a.core.GetFilmCharacters(filmId)
	if err != nil {
		a.lg.Error("Film error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}

	filmResponse := FilmResponse{
		Film:       *film,
		Genres:     genres,
		Rating:     rating,
		Number:     number,
		Directors:  directors,
		Scenarists: scenarists,
		Characters: characters,
	}
	response.Body = filmResponse

	a.SendResponse(w, response)
}

func (a *API) Actor(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Body: nil}
	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, response)
		return
	}

	actorId, err := strconv.ParseUint(r.URL.Query().Get("actor_id"), 10, 64)
	if err != nil {
		response.Status = http.StatusBadRequest
		a.SendResponse(w, response)
		return
	}

	actor, err := a.core.GetActor(actorId)
	if err != nil {
		a.lg.Error("Actor error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}
	if actor == nil {
		response.Status = http.StatusNotFound
		a.SendResponse(w, response)
		return
	}
	career, err := a.core.GetActorsCareer(actorId)
	if err != nil {
		a.lg.Error("Actor error", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}

	actorResponse := ActorResponse{
		Name:      actor.Name,
		Photo:     actor.Photo,
		BirthDate: actor.Birthdate,
		Country:   actor.Country,
		Info:      actor.Info,
		Career:    career,
	}

	response.Body = actorResponse
	a.SendResponse(w, response)
}

func (a *API) Comment(w http.ResponseWriter, r *http.Request) {
	response := Response{Status: http.StatusOK, Body: nil}
	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, response)
		return
	}

	filmId, err := strconv.ParseUint(r.URL.Query().Get("film_id"), 10, 64)
	if err != nil {
		response.Status = http.StatusBadRequest
		a.SendResponse(w, response)
		return
	}
	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.ParseUint(r.URL.Query().Get("per_page"), 10, 64)
	if err != nil {
		pageSize = 10
	}

	comments, err := a.core.GetFilmComments(filmId, (page-1)*pageSize, pageSize)
	if err != nil {
		a.lg.Error("Comment", "err", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, response)
		return
	}

	commentsResponse := CommentResponse{Comments: comments}

	response.Body = commentsResponse
	a.SendResponse(w, response)
}

func (a *API) AddComment(w http.ResponseWriter, r *http.Request) {

}
