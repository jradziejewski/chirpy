package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "supersecret"
	expiresAt := time.Minute * 5

	tokenString, err := MakeJWT(userID, tokenSecret, expiresAt)
	if err != nil {
		t.Fatalf("MakeJWT(%v, %v, %v): expected no error, got %v", userID, tokenSecret, expiresAt, err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		t.Fatalf("MakeJWT(%v, %v, %v): error parsing jwt: %v", userID, tokenSecret, expiresAt, err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		t.Fatalf("MakeJWT(%v, %v, %v): token not valid", userID, tokenSecret, expiresAt)
	}

	if claims.Issuer != "chirpy" {
		t.Fatalf("MakeJWT(%v, %v, %v): expected issuer to be 'chirpy', got %s", userID, tokenSecret, expiresAt, claims.Issuer)
	}
	if claims.Subject != userID.String() {
		t.Fatalf("MakeJWT(%v, %v, %v): expected subject to be '%s', got %s", userID, tokenSecret, expiresAt, userID.String(), claims.Issuer)
	}
	if claims.ExpiresAt.Time.Sub(time.Now()) > expiresAt {
		t.Fatalf("MakeJWT(%v, %v, %v): expected expires_at to be '%v', got %v", userID, tokenSecret, expiresAt, expiresAt, claims.ExpiresAt.Time.Sub(time.Now()))
	}
	if claims.IssuedAt.Time.After(time.Now()) {
		t.Fatalf("MakeJWT(%v, %v, %v): expected issued_at to be '%v', got %v", userID, tokenSecret, expiresAt, time.Now(), claims.IssuedAt.Time.Sub(time.Now()))
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "supersecret"
	expiresAt := time.Minute * 5

	tokenString, err := MakeJWT(userID, tokenSecret, expiresAt)
	if err != nil {
		t.Fatalf("MakeJWT(%v, %v, %v): expected no error, got %v", userID, tokenSecret, expiresAt, err)
	}

	returnedUserID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT(%v, %v): expected no error, got %v", tokenString, tokenSecret, err)
	}

	if returnedUserID != userID {
		t.Fatalf("ValidateJWT(%v, %v): expected %v, got %v", tokenString, tokenSecret, userID, returnedUserID)
	}

	_, err = ValidateJWT(tokenString, "wrong secret")
	if err == nil {
		t.Fatalf("ValidateJWT('', 'whyamidoingthis'): expected error, got no error")
	}

	expiresAt = time.Minute * -5
	tokenString, err = MakeJWT(userID, tokenSecret, expiresAt)
	if err != nil {
		t.Fatalf("MakeJWT(%v, %v, %v): expected no error, got %v", userID, tokenSecret, expiresAt, err)
	}

	returnedUserID, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatalf("ValidateJWT(%v, %v): expected error, got no error", tokenString, tokenSecret)
	}
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "Bearer 1234")

	token, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf(`GetBearerToken("Authorization": "Bearer 1234"): expected no error, got %v`, err)
	}
	if token != "1234" {
		t.Fatalf(`GetBearerToken("Authorization": "Bearer 1234"): expected "1234", got %s`, token)
	}

	header = http.Header{}
	header.Add("UselessHeader", ":)")
	token, err = GetBearerToken(header)

	if err == nil {
		t.Fatalf(`GetBearerToken("UselssHeader": ":)"): expected error, got no error`)
	}

	header = http.Header{}
	header.Add("Authorization", "Bear :)")
	token, err = GetBearerToken(header)

	if err == nil {
		t.Fatalf(`GetBearerToken("Authorization": "Bear :)"): expected error, got no error`)
	}
}
