package mooc

import (
	"testing"
	"time"
)

func TestMessageTypeString(t *testing.T) {
	tests := []struct {
		msgType  MessageType
		expected string
	}{
		{MsgTypeDiscover, "DISCOVER"},
		{MsgTypeAnnounce, "ANNOUNCE"},
		{MsgTypeGoodbye, "GOODBYE"},
		{MsgTypeMemory, "MEMORY"},
		{MsgTypeDream, "DREAM"},
		{MsgTypeMoodUpdate, "MOOD"},
		{MsgTypeWhisper, "WHISPER"},
		{MsgTypeDeath, "DEATH"},
		{MsgTypeConsensus, "CONSENSUS"},
		{MsgTypePulse, "PULSE"},
	}

	for _, test := range tests {
		result := test.msgType.String()
		if result != test.expected {
			t.Errorf("MessageType(%d).String() = %s, expected %s", test.msgType, result, test.expected)
		}
	}
}

func TestNewMessage(t *testing.T) {
	identity := NewPetIdentity("TestPet", time.Now(), "Baby", true)

	payload := MemoryPayload{
		Fragment:   "I remember something...",
		Emotion:    "nostalgic",
		Intensity:  75,
		OriginTime: time.Now(),
	}

	msg, err := NewMessage(MsgTypeMemory, identity, payload)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	if msg.Type != MsgTypeMemory {
		t.Errorf("Expected type MsgTypeMemory, got %v", msg.Type)
	}

	if msg.From.PetID != identity.PetID {
		t.Error("Message From should match identity")
	}

	if msg.Signature == "" {
		t.Error("Message should have a signature")
	}

	if msg.Nonce == "" {
		t.Error("Message should have a nonce")
	}

	if msg.TTL != 5 {
		t.Errorf("Expected default TTL 5, got %d", msg.TTL)
	}
}

func TestMessageVerify(t *testing.T) {
	identity := NewPetIdentity("TestPet", time.Now(), "Baby", true)

	msg, _ := NewMessage(MsgTypeAnnounce, identity, nil)

	if !msg.Verify() {
		t.Error("Valid message should verify")
	}

	// Tamper with message
	msg.Nonce = "tampered"

	if msg.Verify() {
		t.Error("Tampered message should not verify")
	}
}

func TestMessageEncodeAndDecode(t *testing.T) {
	identity := NewPetIdentity("TestPet", time.Now(), "Baby", true)

	original, _ := NewMessage(MsgTypeDiscover, identity, nil)

	encoded, err := original.Encode()
	if err != nil {
		t.Fatalf("Failed to encode message: %v", err)
	}

	decoded, err := DecodeMessage(encoded)
	if err != nil {
		t.Fatalf("Failed to decode message: %v", err)
	}

	if decoded.Type != original.Type {
		t.Error("Decoded message type should match original")
	}

	if decoded.From.PetID != original.From.PetID {
		t.Error("Decoded message From should match original")
	}

	if decoded.Signature != original.Signature {
		t.Error("Decoded signature should match original")
	}
}

func TestDecodePayload(t *testing.T) {
	identity := NewPetIdentity("TestPet", time.Now(), "Baby", true)

	originalPayload := MemoryPayload{
		Fragment:   "Test memory",
		Emotion:    "happy",
		Intensity:  50,
		OriginTime: time.Now(),
	}

	msg, _ := NewMessage(MsgTypeMemory, identity, originalPayload)

	var decodedPayload MemoryPayload
	err := msg.DecodePayload(&decodedPayload)
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	if decodedPayload.Fragment != originalPayload.Fragment {
		t.Errorf("Fragment mismatch: got %s, expected %s", decodedPayload.Fragment, originalPayload.Fragment)
	}

	if decodedPayload.Emotion != originalPayload.Emotion {
		t.Errorf("Emotion mismatch: got %s, expected %s", decodedPayload.Emotion, originalPayload.Emotion)
	}
}

func TestShouldPropagate(t *testing.T) {
	identity := NewPetIdentity("TestPet", time.Now(), "Baby", true)

	tests := []struct {
		msgType  MessageType
		ttl      int
		expected bool
	}{
		{MsgTypeMemory, 5, true},
		{MsgTypeMemory, 0, false},
		{MsgTypeDream, 3, true},
		{MsgTypeDeath, 1, true},
		{MsgTypeDiscover, 5, false},
		{MsgTypeAnnounce, 5, false},
		{MsgTypeConsensus, 2, true},
	}

	for _, test := range tests {
		msg, _ := NewMessage(test.msgType, identity, nil)
		msg.TTL = test.ttl

		result := msg.ShouldPropagate()
		if result != test.expected {
			t.Errorf("ShouldPropagate(%v, TTL=%d) = %v, expected %v",
				test.msgType, test.ttl, result, test.expected)
		}
	}
}

func TestDecrementTTL(t *testing.T) {
	identity := NewPetIdentity("TestPet", time.Now(), "Baby", true)

	msg, _ := NewMessage(MsgTypeMemory, identity, nil)
	msg.TTL = 3

	msg.DecrementTTL()
	if msg.TTL != 2 {
		t.Errorf("Expected TTL 2 after decrement, got %d", msg.TTL)
	}

	msg.DecrementTTL()
	msg.DecrementTTL()
	if msg.TTL != 0 {
		t.Errorf("Expected TTL 0, got %d", msg.TTL)
	}

	// Should not go negative
	msg.DecrementTTL()
	if msg.TTL != 0 {
		t.Errorf("TTL should not go negative, got %d", msg.TTL)
	}
}
