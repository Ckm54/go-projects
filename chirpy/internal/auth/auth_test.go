package auth_test

import (
	"testing"

	"github.com/ckm54/go-projects/chirpy/internal/auth"
)

// Test that hashing a password returns a non-empty hash and no error.
func TestHashPassword(t *testing.T) {
	password := "super-secret"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}
}

// Test that a correct password matches its hash.
func TestCheckPasswordHash_Valid(t *testing.T) {
	password := "mypassword123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	valid, err := auth.CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("unexpected error checking password hash: %v", err)
	}

	if !valid {
		t.Fatalf("expected password to be valid")
	}
}

// Test that an incorrect password does not match the hash.
func TestCheckPasswordHash_Invalid(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	valid, err := auth.CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("unexpected error checking password hash: %v", err)
	}

	if valid {
		t.Fatalf("expected password to be invalid")
	}
}

// Test that a malformed hash returns an error.
func TestCheckPasswordHash_MalformedHash(t *testing.T) {
	password := "test123"
	badHash := "not-a-valid-hash"

	_, err := auth.CheckPasswordHash(password, badHash)
	if err == nil {
		t.Fatalf("expected error when checking malformed hash")
	}
}
