package external

import (
	"testing"
)

func TestBcryptPasswordHasher_Hash(t *testing.T) {
	hasher := NewBcryptPasswordHasher(10) // Lower cost for faster tests

	password := "testpassword123"
	hash, err := hasher.Hash(password)

	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	if hash == "" {
		t.Error("hash should not be empty")
	}

	if hash == password {
		t.Error("hash should not equal plaintext password")
	}
}

func TestBcryptPasswordHasher_Compare(t *testing.T) {
	hasher := NewBcryptPasswordHasher(10)
	password := "testpassword123"

	hash, _ := hasher.Hash(password)

	t.Run("correct password", func(t *testing.T) {
		err := hasher.Compare(password, hash)
		if err != nil {
			t.Errorf("Compare() with correct password should return nil, got %v", err)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		err := hasher.Compare("wrongpassword", hash)
		if err == nil {
			t.Error("Compare() with wrong password should return error")
		}
	})
}

func TestBcryptPasswordHasher_DifferentHashesForSamePassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher(10)
	password := "testpassword123"

	hash1, _ := hasher.Hash(password)
	hash2, _ := hasher.Hash(password)

	if hash1 == hash2 {
		t.Error("same password should produce different hashes (due to salt)")
	}

	// Both should still validate
	if err := hasher.Compare(password, hash1); err != nil {
		t.Error("hash1 should validate")
	}
	if err := hasher.Compare(password, hash2); err != nil {
		t.Error("hash2 should validate")
	}
}
