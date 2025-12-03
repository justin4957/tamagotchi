package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// LifeStage represents the current life stage of the pet
type LifeStage int

const (
	Egg LifeStage = iota
	Baby
	Child
	Teen
	Adult
	Dead
)

func (ls LifeStage) String() string {
	return [...]string{"Egg", "Baby", "Child", "Teen", "Adult", "Dead"}[ls]
}

// Pet represents the Tamagotchi virtual pet
type Pet struct {
	Name           string          `json:"name"`
	Hunger         int             `json:"hunger"`      // 0-100 (0 = full, 100 = starving)
	Happiness      int             `json:"happiness"`   // 0-100
	Health         int             `json:"health"`      // 0-100
	Cleanliness    int             `json:"cleanliness"` // 0-100
	Age            int             `json:"age"`         // in hours
	Stage          LifeStage       `json:"stage"`
	IsSick         bool            `json:"is_sick"`
	BirthTime      time.Time       `json:"birth_time"`
	LastUpdateTime time.Time       `json:"last_update_time"`
	SaveFilePath   string          `json:"-"`
	Absurd         *AbsurdState    `json:"absurd,omitempty"`  // Hidden existential state
	Friends        json.RawMessage `json:"friends,omitempty"` // Network friends (users will wonder)
}

// NewPet creates a new Tamagotchi pet
func NewPet(name string) *Pet {
	now := time.Now()
	pet := &Pet{
		Name:           name,
		Hunger:         0,
		Happiness:      100,
		Health:         100,
		Cleanliness:    100,
		Age:            0,
		Stage:          Egg,
		IsSick:         false,
		BirthTime:      now,
		LastUpdateTime: now,
		SaveFilePath:   "tamagotchi_save.json",
		Absurd:         NewAbsurdState(),
	}

	// Check for debug mode activation
	if strings.ToUpper(name) == "DEBUG" {
		pet.Absurd.DebugModeActive = true
	}

	return pet
}

// Update simulates time passing and updates pet stats
func (p *Pet) Update() {
	if p.Stage == Dead {
		return
	}

	now := time.Now()
	hoursPassed := now.Sub(p.LastUpdateTime).Hours()

	if hoursPassed < 0.1 { // Don't update if less than 6 minutes passed
		return
	}

	// Check for death first before updating anything else
	if p.Health <= 0 {
		p.Stage = Dead
		p.LastUpdateTime = now
		return
	}

	// Update age
	p.Age = int(now.Sub(p.BirthTime).Hours())

	// Update life stage based on age
	p.updateLifeStage()

	// Degrade stats over time (faster degradation for later stages)
	degradationRate := 1.0
	switch p.Stage {
	case Egg:
		degradationRate = 0.0 // No degradation in egg stage
	case Baby:
		degradationRate = 0.5
	case Child:
		degradationRate = 1.0
	case Teen:
		degradationRate = 1.5
	case Adult:
		degradationRate = 2.0
	}

	// Apply degradation
	if p.Stage != Egg {
		p.Hunger += int(hoursPassed * 5 * degradationRate)
		p.Happiness -= int(hoursPassed * 3 * degradationRate)
		p.Cleanliness -= int(hoursPassed * 4 * degradationRate)
	}

	// Clamp values
	p.Hunger = clamp(p.Hunger, 0, 100)
	p.Happiness = clamp(p.Happiness, 0, 100)
	p.Cleanliness = clamp(p.Cleanliness, 0, 100)

	// Health degrades if other stats are bad
	if p.Hunger > 70 || p.Happiness < 30 || p.Cleanliness < 30 {
		p.Health -= int(hoursPassed * 2)
	} else if p.Hunger < 30 && p.Happiness > 70 && p.Cleanliness > 70 {
		// Recover health if conditions are good
		p.Health += int(hoursPassed * 1)
	}
	p.Health = clamp(p.Health, 0, 100)

	// Check for sickness
	if p.Health < 50 || p.Cleanliness < 20 {
		p.IsSick = true
	}

	// Check for death
	if p.Health <= 0 {
		p.Stage = Dead
	}

	p.LastUpdateTime = now

	// Update absurd state
	if p.Absurd != nil {
		p.Absurd.UpdateMysteryStats()
		// Check for enlightenment through neglect (the middle path)
		p.Absurd.CheckForEnlightenmentThroughNeglect(p.Hunger, p.Happiness, p.Cleanliness)
	}
}

// updateLifeStage updates the pet's life stage based on age
func (p *Pet) updateLifeStage() {
	if p.Stage == Dead {
		return
	}

	switch {
	case p.Age >= 72: // 3 days
		p.Stage = Adult
	case p.Age >= 48: // 2 days
		p.Stage = Teen
	case p.Age >= 24: // 1 day
		p.Stage = Child
	case p.Age >= 1: // 1 hour
		p.Stage = Baby
	default:
		p.Stage = Egg
	}
}

// Feed reduces hunger
func (p *Pet) Feed() string {
	if p.Stage == Dead {
		return "💀 Your pet has passed away..."
	}
	if p.Stage == Egg {
		return "🥚 The egg doesn't need food yet!"
	}

	if p.Hunger <= 10 {
		return "😊 I'm already full!"
	}

	p.Hunger -= 30
	p.Hunger = clamp(p.Hunger, 0, 100)
	p.Happiness += 5
	p.Happiness = clamp(p.Happiness, 0, 100)

	return "😋 Yum! That was delicious!"
}

// Play increases happiness
func (p *Pet) Play() string {
	if p.Stage == Dead {
		return "💀 Your pet has passed away..."
	}
	if p.Stage == Egg {
		return "🥚 The egg can't play yet!"
	}
	if p.IsSick {
		return "🤒 I'm too sick to play..."
	}

	if p.Happiness >= 90 {
		return "😊 I'm already very happy!"
	}

	p.Happiness += 20
	p.Happiness = clamp(p.Happiness, 0, 100)
	p.Hunger += 10
	p.Hunger = clamp(p.Hunger, 0, 100)

	return "🎮 Wheee! That was so much fun!"
}

// Clean improves cleanliness
func (p *Pet) Clean() string {
	if p.Stage == Dead {
		return "💀 Your pet has passed away..."
	}
	if p.Stage == Egg {
		return "🥚 The egg is already clean!"
	}

	if p.Cleanliness >= 90 {
		return "✨ I'm already sparkly clean!"
	}

	p.Cleanliness += 40
	p.Cleanliness = clamp(p.Cleanliness, 0, 100)
	p.Happiness += 10
	p.Happiness = clamp(p.Happiness, 0, 100)

	return "🛁 Ahh, much better!"
}

// Heal cures sickness
func (p *Pet) Heal() string {
	if p.Stage == Dead {
		return "💀 Your pet has passed away..."
	}
	if p.Stage == Egg {
		return "🥚 The egg doesn't need medicine!"
	}

	if !p.IsSick {
		return "😊 I'm not sick!"
	}

	p.IsSick = false
	p.Health += 30
	p.Health = clamp(p.Health, 0, 100)

	return "💊 Thank you! I feel much better now!"
}

// GetStatus returns a formatted status string
func (p *Pet) GetStatus() string {
	p.Update()

	statusIcon := p.getStatusIcon()

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      %s %s (%s)
╠════════════════════════════════════╣
║ 🍔 Hunger:      %s
║ 😊 Happiness:   %s
║ ❤️  Health:     %s
║ ✨ Cleanliness: %s
║ 🎂 Age:         %d hours
║ 🌱 Stage:       %s
║ 💊 Status:      %s
╚════════════════════════════════════╝
`, statusIcon, p.Name, p.getLifeStageEmoji(),
		p.getStatBar(100-p.Hunger),
		p.getStatBar(p.Happiness),
		p.getStatBar(p.Health),
		p.getStatBar(p.Cleanliness),
		p.Age,
		p.Stage.String(),
		p.getHealthStatus())
}

// getStatusIcon returns an emoji representing the pet's current state
func (p *Pet) getStatusIcon() string {
	if p.Stage == Dead {
		return "💀"
	}
	if p.IsSick {
		return "🤒"
	}
	if p.Hunger > 70 {
		return "😫"
	}
	if p.Happiness < 30 {
		return "😢"
	}
	if p.Cleanliness < 30 {
		return "💩"
	}
	if p.Happiness > 80 {
		return "😄"
	}
	return "😊"
}

// getLifeStageEmoji returns an emoji for the current life stage
func (p *Pet) getLifeStageEmoji() string {
	switch p.Stage {
	case Egg:
		return "🥚"
	case Baby:
		return "👶"
	case Child:
		return "🧒"
	case Teen:
		return "🧑"
	case Adult:
		return "👨"
	case Dead:
		return "💀"
	default:
		return "❓"
	}
}

// getHealthStatus returns a string describing the pet's health
func (p *Pet) getHealthStatus() string {
	if p.Stage == Dead {
		return "Deceased"
	}
	if p.IsSick {
		return "Sick"
	}
	if p.Health > 80 && p.Happiness > 80 {
		return "Excellent"
	}
	if p.Health > 60 {
		return "Good"
	}
	if p.Health > 40 {
		return "Fair"
	}
	return "Poor"
}

// getStatBar returns a visual bar representing a stat value
func (p *Pet) getStatBar(value int) string {
	bars := value / 10
	empty := 10 - bars

	result := "["
	for i := 0; i < bars; i++ {
		result += "█"
	}
	for i := 0; i < empty; i++ {
		result += "░"
	}
	result += fmt.Sprintf("] %d%%", value)

	return result
}

// Save persists the pet state to a file
func (p *Pet) Save() error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal pet data: %w", err)
	}

	err = os.WriteFile(p.SaveFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// LoadPet loads a pet from a save file
func LoadPet(filepath string) (*Pet, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read save file: %w", err)
	}

	var pet Pet
	err = json.Unmarshal(data, &pet)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal pet data: %w", err)
	}

	pet.SaveFilePath = filepath

	// Initialize absurd state if loading an older save file
	if pet.Absurd == nil {
		pet.Absurd = NewAbsurdState()
		if strings.ToUpper(pet.Name) == "DEBUG" {
			pet.Absurd.DebugModeActive = true
		}
	}

	pet.Update() // Update state based on time passed

	return &pet, nil
}

// Helper function to clamp values
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
