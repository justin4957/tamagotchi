package main

import (
	"testing"
	"time"
)

func TestNewPet(t *testing.T) {
	pet := NewPet("TestPet")

	if pet.Name != "TestPet" {
		t.Errorf("Expected name 'TestPet', got '%s'", pet.Name)
	}

	if pet.Hunger != 0 {
		t.Errorf("Expected hunger 0, got %d", pet.Hunger)
	}

	if pet.Happiness != 100 {
		t.Errorf("Expected happiness 100, got %d", pet.Happiness)
	}

	if pet.Health != 100 {
		t.Errorf("Expected health 100, got %d", pet.Health)
	}

	if pet.Stage != Egg {
		t.Errorf("Expected stage Egg, got %v", pet.Stage)
	}
}

func TestFeed(t *testing.T) {
	pet := NewPet("TestPet")
	pet.Stage = Baby // Change from egg to baby
	pet.Hunger = 50

	result := pet.Feed()

	if pet.Hunger > 50 {
		t.Errorf("Expected hunger to decrease, got %d", pet.Hunger)
	}

	if result == "" {
		t.Error("Expected feed result message")
	}
}

func TestPlay(t *testing.T) {
	pet := NewPet("TestPet")
	pet.Stage = Baby
	pet.Happiness = 50

	result := pet.Play()

	if pet.Happiness <= 50 {
		t.Errorf("Expected happiness to increase, got %d", pet.Happiness)
	}

	if result == "" {
		t.Error("Expected play result message")
	}
}

func TestClean(t *testing.T) {
	pet := NewPet("TestPet")
	pet.Stage = Baby
	pet.Cleanliness = 50

	result := pet.Clean()

	if pet.Cleanliness <= 50 {
		t.Errorf("Expected cleanliness to increase, got %d", pet.Cleanliness)
	}

	if result == "" {
		t.Error("Expected clean result message")
	}
}

func TestLifeStageProgression(t *testing.T) {
	pet := NewPet("TestPet")

	// Test egg stage
	if pet.Stage != Egg {
		t.Errorf("Expected initial stage Egg, got %v", pet.Stage)
	}

	// Simulate 2 hours passing (need to update both birth time and last update time)
	pet.BirthTime = time.Now().Add(-2 * time.Hour)
	pet.LastUpdateTime = time.Now().Add(-2 * time.Hour)
	pet.Update()

	if pet.Stage != Baby {
		t.Errorf("Expected stage Baby after 2 hours, got %v", pet.Stage)
	}

	// Simulate 25 hours passing
	pet.BirthTime = time.Now().Add(-25 * time.Hour)
	pet.LastUpdateTime = time.Now().Add(-1 * time.Hour)
	pet.Update()

	if pet.Stage != Child {
		t.Errorf("Expected stage Child after 25 hours, got %v", pet.Stage)
	}
}

func TestStatDegradation(t *testing.T) {
	pet := NewPet("TestPet")
	// Set birth time to make it a baby
	pet.BirthTime = time.Now().Add(-2 * time.Hour)
	pet.Stage = Baby

	initialHappiness := pet.Happiness
	initialCleanliness := pet.Cleanliness

	// Simulate 1 hour passing
	pet.LastUpdateTime = time.Now().Add(-1 * time.Hour)
	pet.Update()

	if pet.Hunger <= 0 {
		t.Error("Expected hunger to increase over time")
	}

	if pet.Happiness >= initialHappiness {
		t.Errorf("Expected happiness to decrease, was %d, now %d", initialHappiness, pet.Happiness)
	}

	if pet.Cleanliness >= initialCleanliness {
		t.Errorf("Expected cleanliness to decrease, was %d, now %d", initialCleanliness, pet.Cleanliness)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		value    int
		min      int
		max      int
		expected int
	}{
		{50, 0, 100, 50},
		{-10, 0, 100, 0},
		{150, 0, 100, 100},
		{0, 0, 100, 0},
		{100, 0, 100, 100},
	}

	for _, test := range tests {
		result := clamp(test.value, test.min, test.max)
		if result != test.expected {
			t.Errorf("clamp(%d, %d, %d) = %d, expected %d",
				test.value, test.min, test.max, result, test.expected)
		}
	}
}

func TestSickness(t *testing.T) {
	pet := NewPet("TestPet")
	pet.Stage = Baby
	pet.Health = 40
	pet.Cleanliness = 10
	pet.LastUpdateTime = time.Now().Add(-1 * time.Hour)

	pet.Update()

	if !pet.IsSick {
		t.Error("Expected pet to become sick with low health and cleanliness")
	}
}

func TestHeal(t *testing.T) {
	pet := NewPet("TestPet")
	pet.Stage = Baby
	pet.IsSick = true
	pet.Health = 50

	initialHealth := pet.Health
	result := pet.Heal()

	if pet.IsSick {
		t.Error("Expected pet to be cured after healing")
	}

	if pet.Health <= initialHealth {
		t.Errorf("Expected health to increase after healing, was %d, now %d", initialHealth, pet.Health)
	}

	if result == "" {
		t.Error("Expected heal result message")
	}
}

func TestDeath(t *testing.T) {
	pet := NewPet("TestPet")
	pet.BirthTime = time.Now().Add(-2 * time.Hour)
	pet.Stage = Baby
	pet.Health = 0
	pet.LastUpdateTime = time.Now().Add(-1 * time.Hour)

	pet.Update()

	if pet.Stage != Dead {
		t.Errorf("Expected pet to die with 0 health, stage is %v", pet.Stage)
	}
}

func TestEggBehavior(t *testing.T) {
	pet := NewPet("TestPet")

	// Egg shouldn't be able to do actions
	feedResult := pet.Feed()
	if feedResult != "🥚 The egg doesn't need food yet!" {
		t.Error("Expected egg to refuse food")
	}

	playResult := pet.Play()
	if playResult != "🥚 The egg can't play yet!" {
		t.Error("Expected egg to refuse play")
	}
}
