package jwt

import (
	"os"
	"testing"
	"time"
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

	_, err = ParseGuest(token)

	if err == nil {
		t.Errorf("parse guest must fail on parse user token")
	}
}

func TestGuestTokens(t *testing.T) {
	id := "test-guest"
	token, err := NewGuest(id, time.Now().Add(time.Hour*time.Duration(1)))

	if err != nil {
		t.Errorf("new user: %v", err)
		return
	}

	parsed_id, err := ParseGuest(token)

	if err != nil {
		t.Errorf("parse guest: %v", err)
		return
	}

	if *parsed_id != id {
		t.Errorf("parsed_id not equal id: \"%s\" != \"%s\"", *parsed_id, id)
	}

	_, err = ParseUser(token)

	if err == nil {
		t.Errorf("parse user must fail on parse guest token")
	}
}
