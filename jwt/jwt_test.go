package jwt

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("SIGNING_KEY", "test")
	Init()
}

func TestUserTokens(t *testing.T) {
	id := "test-user"
	token, err := NewUser(id)

	if err != nil {
		t.Errorf("new user: %v", err)
		return
	}

	parsed_id, err := ParseUser(token)

	if err != nil {
		t.Errorf("parse user: %v", err)
		return
	}

	if *parsed_id != id {
		t.Errorf("parsed_id not equal id: \"%s\" != \"%s\"", *parsed_id, id)
	}
}
