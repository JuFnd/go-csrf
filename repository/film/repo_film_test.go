package film

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetFilmsByGenre(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Title", "Poster"})

	expect := []FilmItem{
		{Id: 1, Title: "t1", Poster: "url1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Title, item.Poster)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT film.id, film.title, poster FROM film  JOIN films_genre ON film.id = films_genre.id_film WHERE id_genre = $1 ORDER BY release_date DESC OFFSET $2 LIMIT $3")).
		WithArgs(1, 1, 2).
		WillReturnRows(rows)

	repo := &RepoPostgre{
		db: db,
	}

	films, err := repo.GetFilmsByGenre(1, 1, 2)
	if err != nil {
		t.Errorf("GetFilmsByGenre error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(films, expect) {
		t.Errorf("results not match, want %v, have %v", expect, films)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT film.id, film.title, poster FROM film  JOIN films_genre ON film.id = films_genre.id_film WHERE id_genre = $1 ORDER BY release_date DESC OFFSET $2 LIMIT $3")).
		WithArgs(1, 1, 2).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.GetFilmsByGenre(1, 1, 2)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestGetFilms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Title", "Poster"})

	expect := []FilmItem{
		{Id: 1, Title: "t1", Poster: "url1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Title, item.Poster)
	}

	mock.ExpectQuery("SELECT film.id, film.title, poster FROM film").WithArgs(1, 2).WillReturnRows(rows)

	repo := &RepoPostgre{
		db: db,
	}

	films, err := repo.GetFilms(1, 2)
	if err != nil {
		t.Errorf("GetFilms error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(films, expect) {
		t.Errorf("results not match, want %v, have %v", expect, films)
		return
	}

	mock.
		ExpectQuery("SELECT film.id, film.title, poster FROM film").
		WithArgs(1, 2).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.GetFilms(1, 2)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestGetFilm(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Title", "Info", "Poster", "ReleaseDate", "Country", "Mpaa"})

	expect := []FilmItem{
		{Id: 1, Title: "t1", Info: "i1", Poster: "url1", ReleaseDate: "date1", Country: "c1", Mpaa: "12"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Title, item.Info, item.Poster, item.ReleaseDate, item.Country, item.Mpaa)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM film WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &RepoPostgre{
		db: db,
	}

	films, err := repo.GetFilm(1)
	if err != nil {
		t.Errorf("GetFilm error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(films, &expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect, films)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM film WHERE id = $1")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.GetFilm(1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
