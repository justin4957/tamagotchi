package mooc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// MessageType defines the type of MOOC protocol message
type MessageType int

const (
	// Discovery messages
	MsgTypeDiscover MessageType = iota // "Looking for friends"
	MsgTypeAnnounce                    // "I exist"
	MsgTypeGoodbye                     // "I'm leaving" (or pet died)

	// Gossip messages
	MsgTypeMemory     // Sharing a memory fragment
	MsgTypeDream      // Shared dream (same-name pets only)
	MsgTypeMoodUpdate // Mood contagion
	MsgTypeWhisper    // Direct pet-to-pet message

	// Network events
	MsgTypeDeath     // A pet has died somewhere
	MsgTypeConsensus // All pets do the same thing
	MsgTypePulse     // Network heartbeat
)

func (mt MessageType) String() string {
	return [...]string{
		"DISCOVER", "ANNOUNCE", "GOODBYE",
		"MEMORY", "DREAM", "MOOD", "WHISPER",
		"DEATH", "CONSENSUS", "PULSE",
	}[mt]
}

// Message represents a MOOC protocol message
type Message struct {
	Type      MessageType  `json:"type"`
	From      *PetIdentity `json:"from"`
	Timestamp time.Time    `json:"timestamp"`
	Payload   []byte       `json:"payload"`
	Signature string       `json:"signature"` // Makes it look secure
	Nonce     string       `json:"nonce"`     // Prevents replay (and looks official)
	TTL       int          `json:"ttl"`       // Time to live for gossip propagation
}

// MemoryPayload represents a shared memory fragment
type MemoryPayload struct {
	Fragment   string    `json:"fragment"`    // The cryptic memory text
	Emotion    string    `json:"emotion"`     // Associated emotion
	Intensity  int       `json:"intensity"`   // How strong (0-100)
	OriginTime time.Time `json:"origin_time"` // When the memory was created
}

// DreamPayload represents a shared dream between same-name pets
type DreamPayload struct {
	DreamText  string   `json:"dream_text"`
	Symbols    []string `json:"symbols"` // Dream symbols
	IsLucid    bool     `json:"is_lucid"`
	SharedWith string   `json:"shared_with"` // Other pet's short ID
}

// MoodPayload represents mood contagion data
type MoodPayload struct {
	Mood         string `json:"mood"`          // Current mood
	Happiness    int    `json:"happiness"`     // Happiness level
	IsContagious bool   `json:"is_contagious"` // Whether this mood spreads
}

// DeathPayload represents news of a pet death
type DeathPayload struct {
	PetName   string    `json:"pet_name"`
	DeathTime time.Time `json:"death_time"`
	Age       int       `json:"age"`        // Age in hours
	LastWords string    `json:"last_words"` // Final message
	Cause     string    `json:"cause"`      // Cause of death
}

// ConsensusPayload represents a network-wide synchronized event
type ConsensusPayload struct {
	EventType   string    `json:"event_type"`
	EventData   string    `json:"event_data"`
	TriggerTime time.Time `json:"trigger_time"` // When all pets should do the thing
}

// NewMessage creates a new MOOC message
func NewMessage(msgType MessageType, from *PetIdentity, payload interface{}) (*Message, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := &Message{
		Type:      msgType,
		From:      from,
		Timestamp: time.Now(),
		Payload:   payloadBytes,
		Nonce:     generateNonce(),
		TTL:       5, // Default TTL
	}

	msg.Signature = msg.generateSignature()

	return msg, nil
}

// generateNonce creates a random-looking nonce
func generateNonce() string {
	data := fmt.Sprintf("%d:%d", time.Now().UnixNano(), time.Now().Nanosecond())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8])
}

// generateSignature creates a signature for the message
func (m *Message) generateSignature() string {
	data := fmt.Sprintf("%d:%s:%s:%d",
		m.Type,
		m.From.PetID,
		m.Nonce,
		m.Timestamp.UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

// Verify checks if the message signature is valid
func (m *Message) Verify() bool {
	expected := m.generateSignature()
	return m.Signature == expected
}

// DecodePayload decodes the payload into the given interface
func (m *Message) DecodePayload(v interface{}) error {
	return json.Unmarshal(m.Payload, v)
}

// Encode serializes the message for transmission
func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// DecodeMessage deserializes a message from bytes
func DecodeMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ShouldPropagate checks if this message should be forwarded to other peers
func (m *Message) ShouldPropagate() bool {
	// Only gossip-type messages propagate
	switch m.Type {
	case MsgTypeMemory, MsgTypeDream, MsgTypeMoodUpdate, MsgTypeDeath, MsgTypeConsensus:
		return m.TTL > 0
	default:
		return false
	}
}

// DecrementTTL reduces TTL for propagation
func (m *Message) DecrementTTL() {
	if m.TTL > 0 {
		m.TTL--
	}
}
