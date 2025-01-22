package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	pswrd := "1234"
	hashedPswrd, err := HashPassword(pswrd)
	if err != nil {
		t.Fatalf(`HashedPassword("1234") returned error %v`, err)
	}
	if hashedPswrd == pswrd {
		t.Fatalf(`HashedPassword("1234") - passwords aren't different after hashing`)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	pswrd := "1234"
	hashedPswrd, err := HashPassword(pswrd)
	if err != nil {
		t.Fatalf(`HashedPassword("1234") returned error %v`, err)
	}

	err = CheckPasswordHash(pswrd, hashedPswrd)
	if err != nil {
		t.Fatalf(`CheckPasswordHash(%s, %s): expected no error, got %v`, pswrd, hashedPswrd, err)
	}

	err = CheckPasswordHash("random", hashedPswrd)
	if err == nil {
		t.Fatalf(`CheckPasswordHash(%s, %s): expected error, got no error`, pswrd, hashedPswrd)
	}
}
