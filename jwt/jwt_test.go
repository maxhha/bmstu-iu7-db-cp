package jwt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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

	if parsed_id != id {
		t.Errorf("parsed_id not equal id: \"%s\" != \"%s\"", parsed_id, id)
	}
}

func TestWrongAlg(t *testing.T) {
	// Signed with RS512
	token := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ0ZXN0LXVzZXIiLCJpYXQiOjE2NDk5MzM2MTcsInN1YiI6InVzZXIifQ.XjiJ733PI6G_CLUgIIpTnyCVbF-M0Tp_pJlNmgesHHja5poVOEubDs_b_2tN0S_nIAaKAVjC-_UNDPWqBxG71PZ86E2hW4XpHA3eOfGZAa6chVJr6aKMhvxOyUxLljSGmoGa_fWW57uVYtuJtMljDwlI-RsOJ_ChWg3XhDJEJEk"

	_, err := ParseUser(token)
	require.EqualErrorf(
		t,
		err,
		ErrWrongSigningMethod.Error(),
		"Error should be %v, got: %v",
		ErrWrongSigningMethod,
		err,
	)
}
