package comment

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetFilmRating(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Average", "Amount"})

	expectRating := 4.2
	expectAmount := uint64(3)

	rows = rows.AddRow(expectRating, expectAmount)

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT AVG(rating), COUNT(rating) FROM users_comment WHERE id_film")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &RepoPostgre{
		DB: db,
	}

	rating, number, err := repo.GetFilmRating(1)
	if err != nil {
		t.Errorf("GetFilm error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if rating != expectRating {
		t.Errorf("results not match, want %v, have %v", expectRating, rating)
		return
	}
	if number != expectAmount {
		t.Errorf("results not match, want %v, have %v", expectAmount, number)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT AVG(rating), COUNT(rating) FROM users_comment WHERE id_film")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	rating, number, err = repo.GetFilmRating(1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if rating != 0 {
		t.Errorf("expected rating 0, got %f", rating)
	}
	if number != 0 {
		t.Errorf("expected number 0, got %d", number)
	}
}

func TestGetFilmComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Login", "Rating", "Comment"})

	expect := []CommentItem{
		{Username: "l1", Rating: 4, Comment: "c1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Username, item.Rating, item.Comment)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT profile.login, rating, comment FROM users_comment JOIN profile ON users_comment.id_user = profile.id WHERE id_film = $1 OFFSET $2 LIMIT $3")).
		WithArgs(1, 0, 5).
		WillReturnRows(rows)

	repo := &RepoPostgre{
		DB: db,
	}

	comments, err := repo.GetFilmComments(1, 0, 5)
	if err != nil {
		t.Errorf("GetFilm error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(comments, expect) {
		t.Errorf("results not match, want %v, have %v", expect, comments)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT profile.login, rating, comment FROM users_comment JOIN profile ON users_comment.id_user = profile.id WHERE id_film = $1 OFFSET $2 LIMIT $3")).
		WithArgs(1, 0, 5).
		WillReturnError(fmt.Errorf("db_error"))

	comments, err = repo.GetFilmComments(1, 0, 5)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if comments != nil {
		t.Errorf("get comments error, comments should be nil")
	}
}
