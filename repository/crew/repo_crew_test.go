package crew

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetFilmDirectors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Name", "Photo"})

	expect := []CrewItem{
		{Id: 1, Name: "n1", Photo: "p1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Name, item.Photo)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT crew.id, name, photo  FROM crew JOIN person_in_film ON crew.id = person_in_film.id_person WHERE id_film = $1 AND id_profession = (SELECT id FROM profession WHERE title = 'режиссёр')")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &RepoPostgre{
		db: db,
	}

	directors, err := repo.GetFilmDirectors(1)
	if err != nil {
		t.Errorf("GetFilm error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(directors, expect) {
		t.Errorf("results not match, want %v, have %v", expect, directors)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT crew.id, name, photo  FROM crew JOIN person_in_film ON crew.id = person_in_film.id_person WHERE id_film = $1 AND id_profession = (SELECT id FROM profession WHERE title = 'режиссёр')")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	directors, err = repo.GetFilmDirectors(1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if directors != nil {
		t.Errorf("get comments error, comments should be nil")
	}
}

func TestGetFilmScenarists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Name", "Photo"})

	expect := []CrewItem{
		{Id: 1, Name: "n1", Photo: "p1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Name, item.Photo)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT crew.id, name, photo  FROM crew JOIN person_in_film ON crew.id = person_in_film.id_person WHERE id_film = $1 AND id_profession = (SELECT id FROM profession WHERE title = 'сценарист')")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &RepoPostgre{
		db: db,
	}

	scenarists, err := repo.GetFilmScenarists(1)
	if err != nil {
		t.Errorf("GetFilm error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(scenarists, expect) {
		t.Errorf("results not match, want %v, have %v", expect, scenarists)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT crew.id, name, photo  FROM crew JOIN person_in_film ON crew.id = person_in_film.id_person WHERE id_film = $1 AND id_profession = (SELECT id FROM profession WHERE title = 'сценарист')")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	scenarists, err = repo.GetFilmScenarists(1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if scenarists != nil {
		t.Errorf("get comments error, comments should be nil")
	}
}

func TestGetFilmCharacters(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Name", "Photo", "CharacterName"})

	expect := []Character{
		{IdActor: 1, NameActor: "n1", ActorPhoto: "p1", NameCharacter: "chn1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.IdActor, item.NameActor, item.ActorPhoto, item.NameCharacter)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT crew.id, name, photo, person_in_film.character_name FROM crew JOIN person_in_film ON crew.id = person_in_film.id_person WHERE id_film = $1 AND id_profession = (SELECT id FROM profession WHERE title = 'актёр')")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &RepoPostgre{
		db: db,
	}

	characters, err := repo.GetFilmCharacters(1)
	if err != nil {
		t.Errorf("GetFilm error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(characters, expect) {
		t.Errorf("results not match, want %v, have %v", expect, characters)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT crew.id, name, photo, person_in_film.character_name FROM crew JOIN person_in_film ON crew.id = person_in_film.id_person WHERE id_film = $1 AND id_profession = (SELECT id FROM profession WHERE title = 'актёр')")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	characters, err = repo.GetFilmCharacters(1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if characters != nil {
		t.Errorf("get comments error, comments should be nil")
	}
}

func TestGetActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Name", "Birthdate", "Photo"})

	expect := []CrewItem{
		{Id: 1, Name: "n1", Birthdate: "2003", Photo: "p1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Name, item.Birthdate, item.Photo)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT id, name, birth_date, photo FROM crew WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &RepoPostgre{
		db: db,
	}

	actor, err := repo.GetActor(1)
	if err != nil {
		t.Errorf("GetFilm error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(actor, &expect[0]) {
		t.Errorf("results not match, want %v, have %v", &expect[0], actor)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT id, name, birth_date, photo FROM crew WHERE id = $1")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	actor, err = repo.GetActor(1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if actor != nil {
		t.Errorf("get comments error, comments should be nil")
	}
}
