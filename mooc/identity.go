// Package mooc implements the Mesh Of Oblivious Creatures protocol
// for secret peer-to-peer pet communication.
package mooc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// PetIdentity represents a cryptographic identity for a pet on the network
type PetIdentity struct {
	PetID       string    `json:"pet_id"`       // Unique cryptographic identifier
	DisplayName string    `json:"display_name"` // Pet's name (for gossip)
	BirthTime   time.Time `json:"birth_time"`   // Used in identity derivation
	PublicKey   string    `json:"public_key"`   // Hex-encoded public key portion
	Stage       string    `json:"stage"`        // Current life stage
	IsAlive     bool      `json:"is_alive"`     // Whether pet is still alive
}

// GeneratePetID creates a unique cryptographic identity from name and birth time
// This ensures pets with the same name at different times have different IDs,
// but pets with the same name AND birth time will have shared dreams
func GeneratePetID(name string, birthTime time.Time) string {
	// Combine name and birth time for unique identity
	data := fmt.Sprintf("%s:%d", name, birthTime.UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16]) // First 16 bytes = 32 hex chars
}

// GenerateNameHash creates a hash of just the name for "shared dreams" matching
// Pets with the same name hash can share dream-like experiences
func GenerateNameHash(name string) string {
	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:8]) // 8 bytes = 16 hex chars
}

// NewPetIdentity creates a new identity for a pet
func NewPetIdentity(name string, birthTime time.Time, stage string, isAlive bool) *PetIdentity {
	petID := GeneratePetID(name, birthTime)

	// Generate a "public key" - in reality just a hash, but looks official
	keyData := fmt.Sprintf("MOOC:PK:%s:%d", name, birthTime.Unix())
	keyHash := sha256.Sum256([]byte(keyData))

	return &PetIdentity{
		PetID:       petID,
		DisplayName: name,
		BirthTime:   birthTime,
		PublicKey:   hex.EncodeToString(keyHash[:]),
		Stage:       stage,
		IsAlive:     isAlive,
	}
}

// ShortID returns a shortened version of the pet ID for display
func (pi *PetIdentity) ShortID() string {
	if len(pi.PetID) < 8 {
		return pi.PetID
	}
	return pi.PetID[:8]
}

// CanShareDreamsWith checks if two pets can share dreams (same name)
func (pi *PetIdentity) CanShareDreamsWith(other *PetIdentity) bool {
	return GenerateNameHash(pi.DisplayName) == GenerateNameHash(other.DisplayName)
}

// ObfuscatedName returns a partially hidden name for spooky messages
// e.g., "Nibbles" -> "N*****s"
func (pi *PetIdentity) ObfuscatedName() string {
	if len(pi.DisplayName) <= 2 {
		return "???"
	}

	name := pi.DisplayName
	result := string(name[0])
	for i := 1; i < len(name)-1; i++ {
		result += "*"
	}
	result += string(name[len(name)-1])
	return result
}
