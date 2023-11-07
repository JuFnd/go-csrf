package delivery

import (
	"testing"
)

func TestCreateAndKillSession(t *testing.T) {
	login := "testLogin"
	testCore := Core{sessions: make(map[string]Session)}

	sid, _, _ := testCore.CreateSession(login)
	isFound, _ := testCore.FindActiveSession(sid)
	if !isFound {
		t.Errorf("session not found")
	}

	err := testCore.KillSession(sid)
	if err != nil {
		t.Errorf("failed to kill session")
	}

	isFound, _ = testCore.FindActiveSession(sid)
	if isFound {
		t.Errorf("found killed session")
	}
}
