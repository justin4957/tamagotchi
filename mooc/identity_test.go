package mooc

import (
	"testing"
	"time"
)

func TestGeneratePetID(t *testing.T) {
	birthTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	id1 := GeneratePetID("Nibbles", birthTime)
	id2 := GeneratePetID("Nibbles", birthTime)
	id3 := GeneratePetID("Fluffy", birthTime)
	id4 := GeneratePetID("Nibbles", birthTime.Add(1*time.Second))

	// Same inputs should produce same ID
	if id1 != id2 {
		t.Error("Same name and birthTime should produce same ID")
	}

	// Different name should produce different ID
	if id1 == id3 {
		t.Error("Different names should produce different IDs")
	}

	// Same name but different time should produce different ID
	if id1 == id4 {
		t.Error("Same name with different birthTime should produce different IDs")
	}

	// ID should be 32 hex characters
	if len(id1) != 32 {
		t.Errorf("Expected ID length 32, got %d", len(id1))
	}
}

func TestGenerateNameHash(t *testing.T) {
	hash1 := GenerateNameHash("Nibbles")
	hash2 := GenerateNameHash("Nibbles")
	hash3 := GenerateNameHash("Fluffy")

	// Same name should produce same hash
	if hash1 != hash2 {
		t.Error("Same name should produce same hash")
	}

	// Different names should produce different hashes
	if hash1 == hash3 {
		t.Error("Different names should produce different hashes")
	}

	// Hash should be 16 hex characters
	if len(hash1) != 16 {
		t.Errorf("Expected hash length 16, got %d", len(hash1))
	}
}

func TestNewPetIdentity(t *testing.T) {
	birthTime := time.Now()
	identity := NewPetIdentity("TestPet", birthTime, "Baby", true)

	if identity.DisplayName != "TestPet" {
		t.Errorf("Expected DisplayName 'TestPet', got '%s'", identity.DisplayName)
	}

	if identity.Stage != "Baby" {
		t.Errorf("Expected Stage 'Baby', got '%s'", identity.Stage)
	}

	if !identity.IsAlive {
		t.Error("Expected IsAlive to be true")
	}

	if identity.PetID == "" {
		t.Error("PetID should not be empty")
	}

	if identity.PublicKey == "" {
		t.Error("PublicKey should not be empty")
	}
}

func TestShortID(t *testing.T) {
	identity := NewPetIdentity("TestPet", time.Now(), "Adult", true)
	shortID := identity.ShortID()

	if len(shortID) != 8 {
		t.Errorf("Expected ShortID length 8, got %d", len(shortID))
	}

	// ShortID should be prefix of full ID
	if identity.PetID[:8] != shortID {
		t.Error("ShortID should be prefix of PetID")
	}
}

func TestCanShareDreamsWith(t *testing.T) {
	now := time.Now()

	// Two pets with the same name
	pet1 := NewPetIdentity("Nibbles", now, "Baby", true)
	pet2 := NewPetIdentity("Nibbles", now.Add(1*time.Hour), "Child", true)

	// Pet with different name
	pet3 := NewPetIdentity("Fluffy", now, "Baby", true)

	if !pet1.CanShareDreamsWith(pet2) {
		t.Error("Pets with same name should be able to share dreams")
	}

	if pet1.CanShareDreamsWith(pet3) {
		t.Error("Pets with different names should not share dreams")
	}
}

func TestObfuscatedName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"Nibbles", "N*****s"},
		{"Ab", "???"},
		{"A", "???"},
		{"Cat", "C*t"},
		{"Tamagotchi", "T********i"},
	}

	for _, test := range tests {
		identity := NewPetIdentity(test.name, time.Now(), "Baby", true)
		result := identity.ObfuscatedName()

		if result != test.expected {
			t.Errorf("ObfuscatedName(%s) = %s, expected %s", test.name, result, test.expected)
		}
	}
}
