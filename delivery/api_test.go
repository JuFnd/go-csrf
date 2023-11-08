package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilmsPost(t *testing.T) {
	h := httptest.NewRequest(http.MethodPost, "/api/v1/films", nil)
	w := httptest.NewRecorder()

	api := API{}
	api.Films(w, h)
	var response Response

	body, _ := io.ReadAll(w.Body)
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Error("cant unmarshal jsone")
	}

	if response.Status != http.StatusMethodNotAllowed {
		t.Errorf("got incorrect status")
	}
}

func TestSignupGet(t *testing.T) {
	h := httptest.NewRequest(http.MethodGet, "/signup", nil)
	w := httptest.NewRecorder()

	api := API{}
	api.Signup(w, h)
	var response Response

	body, _ := io.ReadAll(w.Body)
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Error("cant unmarshal jsone")
	}

	if response.Status != http.StatusMethodNotAllowed {
		t.Errorf("got incorrect status")
	}
}

func TestSigninGet(t *testing.T) {
	h := httptest.NewRequest(http.MethodGet, "/signin", nil)
	w := httptest.NewRecorder()

	api := API{}
	api.Signin(w, h)
	var response Response

	body, _ := io.ReadAll(w.Body)
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Error("cant unmarshal jsone")
	}

	if response.Status != http.StatusMethodNotAllowed {
		t.Errorf("got incorrect status")
	}
}
