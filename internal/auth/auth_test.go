package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "super-secret-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == password {
		t.Fatal("hash should not equal the raw password")
	}

	ok, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if !ok {
		t.Fatal("expected password to match hash")
	}
}

func TestCheckPasswordHashWrongPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	ok, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if ok {
		t.Fatal("expected password mismatch, got match")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if gotID != userID {
		t.Fatalf("expected userID %v, got %v", userID, gotID)
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "correct-secret", time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error when validating token with wrong secret")
	}
}

func TestValidateJWTExpired(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	// Token already expired
	token, err := MakeJWT(userID, secret, -time.Second)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}
