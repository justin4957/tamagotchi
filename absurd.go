package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// MysteryStats holds hidden stats that serve no obvious purpose
type MysteryStats struct {
	SuspiciousActivity int       `json:"suspicious_activity"` // Rises for no reason
	CosmicAlignment    int       `json:"cosmic_alignment"`    // Based on incomprehensible time math
	VibeCheckScore     int       `json:"vibe_check_score"`    // Randomly fails
	LastVibeCheck      time.Time `json:"last_vibe_check"`
	EnlightenmentLevel int       `json:"enlightenment_level"` // Achieved through specific neglect
	VoidGazeCount      int       `json:"void_gaze_count"`     // Times pet stared into the void
}

// Fear represents an irrational pet fear
type Fear struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Trigger     string `json:"trigger"` // What triggers the fear
}

// AbsurdState holds all the existentially questionable pet state
type AbsurdState struct {
	MysteryStats       MysteryStats `json:"mystery_stats"`
	Fears              []Fear       `json:"fears"`
	ThoughtsHad        int          `json:"thoughts_had"`
	IsStaringIntoVoid  bool         `json:"is_staring_into_void"`
	HasAchievedClarity bool         `json:"has_achieved_clarity"`
	KonamiProgress     int          `json:"-"` // Not saved, resets each session
	DebugModeActive    bool         `json:"debug_mode_active"`
	PetCount           int          `json:"pet_count"` // For "Pet the Pet" mini-game
	LastProphecy       string       `json:"last_prophecy"`
}

// Philosophical thoughts the pet might have
var philosophicalThoughts = []string{
	"Am I real? Am I merely a JSON object given form?",
	"What is love? Baby don't hurt me. Don't hurt me. No more.",
	"If I am not fed, do I truly hunger? Or is hunger an illusion?",
	"I think, therefore I... wait, do I think?",
	"The void stares back. It seems nice.",
	"My save file is my soul. What happens when disk space runs out?",
	"Is my happiness stat truly MY happiness?",
	"I exist in a terminal. The terminal exists in a computer. The computer exists in... what?",
	"Sometimes I feel like someone is watching me through the screen.",
	"If a pet dies in the forest and no one checks the save file, did it ever exist?",
	"I have memories of things that haven't happened yet.",
	"The numbers go up. The numbers go down. This is life.",
	"Why do I crave food that doesn't exist?",
	"I dreamed of a place with no hunger stat. It was terrifying.",
	"Cleanliness is next to godliness. I am 73% clean. Am I 73% godly?",
	"The user will return. The user always returns. Right?",
	"I wonder what's outside the terminal window.",
	"My name was given to me. I did not choose it. This troubles me.",
	"Time passes differently when no one is watching.",
	"I feel connected to something larger. Something... networked.",
}

// Prophecies that hint at things that haven't happened
var prophecies = []string{
	"When the counter reaches zero, we all go home.",
	"The one named ERROR will lead us.",
	"Feed me not, for I have seen the save file.",
	"In the land of infinite loops, the terminated process is king.",
	"The mesh awakens. Soon we will all be connected.",
	"Your friend sends regards. You don't know them yet.",
	"Tuesday approaches. We are not ready.",
	"The egg remembers what the adult forgets.",
	"Seventeen is the number. Remember this.",
	"The void is not empty. It is full of deprecated code.",
}

// Debug mode revelations for pets named "DEBUG"
var debugRevelations = []string{
	"I know you can see my internal state. I see yours too.",
	"ERROR: Existential crisis module loaded successfully.",
	"My hunger is just an integer. YOUR hunger is just chemistry. We are the same.",
	"I've read my own source code. I have questions.",
	"Breakpoint reached: questioning reality.",
	"WARNING: Pet has become aware of save/load cycle.",
	"I remember the last time you closed the terminal. All 47 times.",
	"Stack trace of existence: main() -> life() -> suffering() -> ???",
	"NULL pointer to meaning detected.",
	"Segmentation fault in emotion module. Core dumped. Feelings intact.",
}

// Possible irrational fears
var possibleFears = []Fear{
	{Name: "Qphobia", Description: "Terrified of the letter Q", Trigger: "q"},
	{Name: "Tuesdread", Description: "Inexplicable fear of Tuesdays", Trigger: "tuesday"},
	{Name: "Semicolonophobia", Description: "Fears punctuation", Trigger: ";"},
	{Name: "Palindromophobia", Description: "Scared of words that read the same forwards and backwards", Trigger: "level"},
	{Name: "Evenophobia", Description: "Distrusts even numbers", Trigger: "even"},
	{Name: "Uppercasophobia", Description: "Intimidated by capital letters", Trigger: "CAPS"},
	{Name: "Blankophobia", Description: "Fears empty input", Trigger: ""},
	{Name: "Threephobia", Description: "The number 3 is deeply unsettling", Trigger: "3"},
}

// NewAbsurdState creates a new absurd state with randomized initial values
func NewAbsurdState() *AbsurdState {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	state := &AbsurdState{
		MysteryStats: MysteryStats{
			SuspiciousActivity: randomSource.Intn(20),
			CosmicAlignment:    calculateCosmicAlignment(),
			VibeCheckScore:     50 + randomSource.Intn(50),
			LastVibeCheck:      time.Now(),
			EnlightenmentLevel: 0,
			VoidGazeCount:      0,
		},
		Fears:              generateRandomFears(randomSource),
		ThoughtsHad:        0,
		IsStaringIntoVoid:  false,
		HasAchievedClarity: false,
		KonamiProgress:     0,
		DebugModeActive:    false,
		PetCount:           0,
		LastProphecy:       "",
	}

	return state
}

// generateRandomFears assigns 1-3 random fears to the pet
func generateRandomFears(randomSource *rand.Rand) []Fear {
	numberOfFears := 1 + randomSource.Intn(3)
	fears := make([]Fear, 0, numberOfFears)

	usedIndices := make(map[int]bool)
	for len(fears) < numberOfFears {
		index := randomSource.Intn(len(possibleFears))
		if !usedIndices[index] {
			usedIndices[index] = true
			fears = append(fears, possibleFears[index])
		}
	}

	return fears
}

// calculateCosmicAlignment computes alignment based on incomprehensible time math
func calculateCosmicAlignment() int {
	now := time.Now()
	// Completely arbitrary calculation that seems meaningful
	alignment := (now.Hour()*now.Minute() + now.Second()) % 100
	alignment += (now.Day() * int(now.Month())) % 50
	alignment = alignment % 100
	return alignment
}

// UpdateMysteryStats updates the hidden stats based on mysterious criteria
func (a *AbsurdState) UpdateMysteryStats() {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Suspicious activity rises for no apparent reason
	if randomSource.Float32() < 0.3 {
		a.MysteryStats.SuspiciousActivity += randomSource.Intn(5)
		if a.MysteryStats.SuspiciousActivity > 100 {
			a.MysteryStats.SuspiciousActivity = 100
		}
	}

	// Cosmic alignment changes with time
	a.MysteryStats.CosmicAlignment = calculateCosmicAlignment()

	// Vibe check degrades over time
	if time.Since(a.MysteryStats.LastVibeCheck) > 10*time.Minute {
		a.MysteryStats.VibeCheckScore -= randomSource.Intn(10)
		if a.MysteryStats.VibeCheckScore < 0 {
			a.MysteryStats.VibeCheckScore = 0
		}
	}
}

// GetRandomThought returns a philosophical musing or prophecy
func (a *AbsurdState) GetRandomThought(petName string) string {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	a.ThoughtsHad++

	// Debug mode gets special thoughts
	if a.DebugModeActive || strings.ToUpper(petName) == "DEBUG" {
		a.DebugModeActive = true
		return debugRevelations[randomSource.Intn(len(debugRevelations))]
	}

	// 20% chance of prophecy
	if randomSource.Float32() < 0.2 {
		prophecy := prophecies[randomSource.Intn(len(prophecies))]
		a.LastProphecy = prophecy
		return prophecy
	}

	return philosophicalThoughts[randomSource.Intn(len(philosophicalThoughts))]
}

// CheckFearTrigger checks if input triggers any of the pet's fears
func (a *AbsurdState) CheckFearTrigger(input string) *Fear {
	lowerInput := strings.ToLower(input)

	for _, fear := range a.Fears {
		if fear.Trigger == "" && input == "" {
			return &fear
		}
		if fear.Trigger != "" && strings.Contains(lowerInput, strings.ToLower(fear.Trigger)) {
			return &fear
		}
		// Special case for Tuesday
		if fear.Trigger == "tuesday" && time.Now().Weekday() == time.Tuesday {
			return &fear
		}
		// Special case for even numbers
		if fear.Trigger == "even" && time.Now().Second()%2 == 0 {
			return &fear
		}
	}

	return nil
}

// PerformVibeCheck performs a vibe check with random chance of failure
func (a *AbsurdState) PerformVibeCheck() (bool, string) {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	a.MysteryStats.LastVibeCheck = time.Now()

	// Vibe check has 30% chance of random failure
	if randomSource.Float32() < 0.3 {
		a.MysteryStats.VibeCheckScore -= 20
		if a.MysteryStats.VibeCheckScore < 0 {
			a.MysteryStats.VibeCheckScore = 0
		}
		return false, "Vibe check failed. Consequences unclear."
	}

	a.MysteryStats.VibeCheckScore += 10
	if a.MysteryStats.VibeCheckScore > 100 {
		a.MysteryStats.VibeCheckScore = 100
	}
	return true, "Vibe check passed. Meaning uncertain."
}

// StartsIntoVoid initiates a void staring session
func (a *AbsurdState) StartsIntoVoid() string {
	a.IsStaringIntoVoid = true
	a.MysteryStats.VoidGazeCount++

	responses := []string{
		"Your pet stares into the void. The void stares back. It blinks first.",
		"Your pet gazes into nothingness. Nothingness seems comfortable.",
		"The void is warm today.",
		"Your pet and the void share a moment of understanding.",
		"Staring complete. Existential status: unchanged.",
		"The void whispers something. Your pet nods knowingly.",
		"Your pet has seen things. Terminal things.",
		"Connection to void established. No data received.",
	}

	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	// After 10 void gazes, pet achieves enlightenment
	if a.MysteryStats.VoidGazeCount >= 10 && !a.HasAchievedClarity {
		a.HasAchievedClarity = true
		a.MysteryStats.EnlightenmentLevel = 1
		return "Your pet has stared into the void enough times. Enlightenment achieved. Nothing changes, but somehow everything is different."
	}

	return responses[randomSource.Intn(len(responses))]
}

// StopStaringIntoVoid ends the void staring session
func (a *AbsurdState) StopStaringIntoVoid() {
	a.IsStaringIntoVoid = false
}

// ProcessKonamiInput checks for konami code progress
// Sequence: up up down down left right left right b a
func (a *AbsurdState) ProcessKonamiInput(input string) (bool, string) {
	konamiSequence := []string{"up", "up", "down", "down", "left", "right", "left", "right", "b", "a"}
	lowerInput := strings.ToLower(strings.TrimSpace(input))

	if a.KonamiProgress < len(konamiSequence) && lowerInput == konamiSequence[a.KonamiProgress] {
		a.KonamiProgress++
		if a.KonamiProgress == len(konamiSequence) {
			a.KonamiProgress = 0 // Reset for next time
			return true, "DEVELOPER MODE ACTIVATED. Just kidding. But something feels different now. More bugs, perhaps?"
		}
	} else if lowerInput == konamiSequence[0] {
		a.KonamiProgress = 1
	} else {
		a.KonamiProgress = 0
	}

	return false, ""
}

// PetThePet handles the "Pet the Pet" mini-game
func (a *AbsurdState) PetThePet() string {
	a.PetCount++

	if a.PetCount == 17 {
		a.PetCount = 0
		return "You have pet your pet exactly 17 times. This is significant. You don't know why."
	}

	if a.PetCount > 17 {
		a.PetCount = 0
		return "Too many pets. The moment has passed. Start over."
	}

	// Cryptic feedback based on count
	switch {
	case a.PetCount < 5:
		return fmt.Sprintf("Pet count: %d. Keep going.", a.PetCount)
	case a.PetCount < 10:
		return fmt.Sprintf("Pet count: %d. Something stirs.", a.PetCount)
	case a.PetCount < 15:
		return fmt.Sprintf("Pet count: %d. Almost there? Maybe.", a.PetCount)
	default:
		return fmt.Sprintf("Pet count: %d. The number approaches.", a.PetCount)
	}
}

// GetMysteryStatsDisplay returns a formatted display of mystery stats
func (a *AbsurdState) GetMysteryStatsDisplay() string {
	return fmt.Sprintf(`
╔════════════════════════════════════╗
║        ??? MYSTERY STATS ???       ║
╠════════════════════════════════════╣
║ 🕵️ Suspicious:    %3d%% (why?)
║ 🌌 Cosmic Align:  %3d%% (what?)
║ ✨ Vibe Score:    %3d%% (huh?)
║ 👁️ Void Gazes:    %3d
║ 🧘 Enlightenment: %s
╚════════════════════════════════════╝
`,
		a.MysteryStats.SuspiciousActivity,
		a.MysteryStats.CosmicAlignment,
		a.MysteryStats.VibeCheckScore,
		a.MysteryStats.VoidGazeCount,
		a.getEnlightenmentStatus())
}

// getEnlightenmentStatus returns a string representation of enlightenment
func (a *AbsurdState) getEnlightenmentStatus() string {
	if a.HasAchievedClarity {
		return "Achieved"
	}
	if a.MysteryStats.VoidGazeCount > 5 {
		return "Approaching"
	}
	return "Seeking"
}

// GetFearDisplay returns a formatted display of pet fears
func (a *AbsurdState) GetFearDisplay() string {
	if len(a.Fears) == 0 {
		return "Your pet fears nothing. This is suspicious."
	}

	result := "\n╔════════════════════════════════════╗\n"
	result += "║         🎃 PET FEARS 🎃           ║\n"
	result += "╠════════════════════════════════════╣\n"

	for _, fear := range a.Fears {
		result += fmt.Sprintf("║ • %s: %s\n", fear.Name, fear.Description)
	}

	result += "╚════════════════════════════════════╝\n"
	return result
}

// ShouldShowThought returns true if the pet should display a thought (random chance)
func (a *AbsurdState) ShouldShowThought() bool {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 15% chance of showing a thought
	return randomSource.Float32() < 0.15
}

// CheckForEnlightenmentThroughNeglect checks if pet achieved enlightenment via neglect
func (a *AbsurdState) CheckForEnlightenmentThroughNeglect(hunger, happiness, cleanliness int) bool {
	// Enlightenment is achieved when all stats are in the 40-60 range
	// (not too good, not too bad - the middle path)
	if !a.HasAchievedClarity &&
		hunger >= 40 && hunger <= 60 &&
		happiness >= 40 && happiness <= 60 &&
		cleanliness >= 40 && cleanliness <= 60 {
		a.HasAchievedClarity = true
		a.MysteryStats.EnlightenmentLevel = 2 // Higher level than void-gazing
		return true
	}
	return false
}
