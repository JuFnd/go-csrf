package delivery

import (
	"log/slog"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/csrf"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/repository/session"
)

func TestCreateAndCheckCsrf(t *testing.T) {
	logFile, _ := os.Create("log.log")
	lg := slog.New(slog.NewJSONHandler(logFile, nil))

	csrf, err := csrf.GetCsrfRepo(lg)

	if err != nil {
		t.Errorf("Csrf repository is not responding")
	}

	testCore := Core{
		csrfTokens: *csrf,
	}

	sid, err := testCore.CreateCsrfToken()

	if err != nil {
		t.Error("failed create csrf token")
	}

	isFound, err := testCore.CheckCsrfToken(sid)
	if !isFound {
		t.Errorf("csrf not found")
	}
}

func TestCreateAndKillSession(t *testing.T) {
	logFile, _ := os.Create("log.log")
	lg := slog.New(slog.NewJSONHandler(logFile, nil))

	newSession, err := session.GetSessionRepo(lg)

	if err != nil {
		t.Errorf("Session repository is not responding")
	}

	login := "testLogin"
	testCore := Core{
		sessions: *newSession,
	}

	sid, _, _ := testCore.CreateSession(login)
	isFound, _ := testCore.FindActiveSession(sid)
	if !isFound {
		t.Errorf("session not found")
	}

	err = testCore.KillSession(sid)
	if err != nil {
		t.Errorf("failed to kill session")
	}

	isFound, _ = testCore.FindActiveSession(sid)
	if isFound {
		t.Errorf("found killed session")
	}
}
