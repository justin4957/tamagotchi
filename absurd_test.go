package main

import (
	"strings"
	"testing"
)

func TestNewAbsurdState(t *testing.T) {
	state := NewAbsurdState()

	if state == nil {
		t.Fatal("Expected non-nil AbsurdState")
	}

	if state.MysteryStats.SuspiciousActivity < 0 || state.MysteryStats.SuspiciousActivity > 100 {
		t.Errorf("SuspiciousActivity should be 0-100, got %d", state.MysteryStats.SuspiciousActivity)
	}

	if state.MysteryStats.CosmicAlignment < 0 || state.MysteryStats.CosmicAlignment > 100 {
		t.Errorf("CosmicAlignment should be 0-100, got %d", state.MysteryStats.CosmicAlignment)
	}

	if len(state.Fears) < 1 || len(state.Fears) > 3 {
		t.Errorf("Expected 1-3 fears, got %d", len(state.Fears))
	}

	if state.HasAchievedClarity {
		t.Error("New state should not have achieved clarity")
	}
}

func TestGetRandomThought(t *testing.T) {
	state := NewAbsurdState()

	thought := state.GetRandomThought("TestPet")

	if thought == "" {
		t.Error("Expected non-empty thought")
	}

	if state.ThoughtsHad != 1 {
		t.Errorf("Expected ThoughtsHad to be 1, got %d", state.ThoughtsHad)
	}
}

func TestDebugModeThoughts(t *testing.T) {
	state := NewAbsurdState()

	thought := state.GetRandomThought("DEBUG")

	if !state.DebugModeActive {
		t.Error("DEBUG name should activate debug mode")
	}

	// Debug thoughts should contain certain keywords
	foundDebugThought := false
	for _, debugThought := range debugRevelations {
		if thought == debugThought {
			foundDebugThought = true
			break
		}
	}

	if !foundDebugThought {
		t.Error("DEBUG pet should receive debug revelations")
	}
}

func TestCheckFearTrigger(t *testing.T) {
	state := NewAbsurdState()

	// Create a known fear
	state.Fears = []Fear{
		{Name: "Qphobia", Description: "Terrified of the letter Q", Trigger: "q"},
	}

	fear := state.CheckFearTrigger("question")

	if fear == nil {
		t.Error("Expected fear to be triggered by 'question' containing 'q'")
	}

	if fear != nil && fear.Name != "Qphobia" {
		t.Errorf("Expected Qphobia, got %s", fear.Name)
	}

	// Test no trigger
	noFear := state.CheckFearTrigger("hello")
	if noFear != nil {
		t.Error("Expected no fear trigger for 'hello'")
	}
}

func TestPerformVibeCheck(t *testing.T) {
	state := NewAbsurdState()
	initialScore := state.MysteryStats.VibeCheckScore

	// Run multiple vibe checks to test randomness
	passCount := 0
	failCount := 0

	for i := 0; i < 100; i++ {
		state.MysteryStats.VibeCheckScore = 50 // Reset for each test
		passed, message := state.PerformVibeCheck()
		if passed {
			passCount++
		} else {
			failCount++
		}
		if message == "" {
			t.Error("Expected non-empty vibe check message")
		}
	}

	// With 30% fail rate, we should see both passes and fails in 100 attempts
	if passCount == 0 {
		t.Error("Expected at least some vibe checks to pass")
	}
	if failCount == 0 {
		t.Error("Expected at least some vibe checks to fail (30% chance)")
	}

	// Verify initial score was set
	if initialScore < 50 || initialScore > 100 {
		t.Errorf("Initial vibe check score should be 50-100, got %d", initialScore)
	}
}

func TestStareIntoVoid(t *testing.T) {
	state := NewAbsurdState()

	message := state.StartsIntoVoid()

	if message == "" {
		t.Error("Expected non-empty void staring message")
	}

	if state.MysteryStats.VoidGazeCount != 1 {
		t.Errorf("Expected VoidGazeCount to be 1, got %d", state.MysteryStats.VoidGazeCount)
	}

	// Test enlightenment after 10 gazes
	for i := 0; i < 10; i++ {
		state.StartsIntoVoid()
	}

	if !state.HasAchievedClarity {
		t.Error("Expected enlightenment after 10 void gazes")
	}

	if state.MysteryStats.EnlightenmentLevel != 1 {
		t.Errorf("Expected EnlightenmentLevel 1, got %d", state.MysteryStats.EnlightenmentLevel)
	}
}

func TestKonamiCode(t *testing.T) {
	state := NewAbsurdState()

	// Test correct sequence
	konamiSequence := []string{"up", "up", "down", "down", "left", "right", "left", "right", "b", "a"}

	for i, input := range konamiSequence[:9] {
		activated, _ := state.ProcessKonamiInput(input)
		if activated {
			t.Errorf("Konami code should not activate before complete sequence (step %d)", i)
		}
	}

	activated, message := state.ProcessKonamiInput("a")
	if !activated {
		t.Error("Konami code should activate after complete sequence")
	}
	if message == "" {
		t.Error("Expected Konami activation message")
	}

	// Test reset on wrong input
	state.ProcessKonamiInput("up")
	state.ProcessKonamiInput("up")
	state.ProcessKonamiInput("wrong")

	if state.KonamiProgress != 0 {
		t.Errorf("Konami progress should reset on wrong input, got %d", state.KonamiProgress)
	}
}

func TestPetThePet(t *testing.T) {
	state := NewAbsurdState()

	// Pet 16 times
	for i := 1; i <= 16; i++ {
		message := state.PetThePet()
		if message == "" {
			t.Errorf("Expected message on pet %d", i)
		}
	}

	// The 17th pet should be special
	message := state.PetThePet()
	if !strings.Contains(message, "17") {
		t.Errorf("Expected 17th pet to be significant, got: %s", message)
	}

	if state.PetCount != 0 {
		t.Errorf("Expected PetCount to reset after 17, got %d", state.PetCount)
	}
}

func TestEnlightenmentThroughNeglect(t *testing.T) {
	state := NewAbsurdState()

	// Not in middle path - should not achieve enlightenment
	achieved := state.CheckForEnlightenmentThroughNeglect(20, 80, 90)
	if achieved {
		t.Error("Should not achieve enlightenment when not on middle path")
	}

	// Middle path - all stats 40-60
	achieved = state.CheckForEnlightenmentThroughNeglect(50, 50, 50)
	if !achieved {
		t.Error("Should achieve enlightenment on middle path (all stats 40-60)")
	}

	if !state.HasAchievedClarity {
		t.Error("HasAchievedClarity should be true after middle path enlightenment")
	}

	if state.MysteryStats.EnlightenmentLevel != 2 {
		t.Errorf("Middle path enlightenment should be level 2, got %d", state.MysteryStats.EnlightenmentLevel)
	}
}

func TestUpdateMysteryStats(t *testing.T) {
	state := NewAbsurdState()
	initialAlignment := state.MysteryStats.CosmicAlignment

	state.UpdateMysteryStats()

	// Cosmic alignment should potentially change (based on time)
	// We can't guarantee it changes, but we can verify it stays in range
	if state.MysteryStats.CosmicAlignment < 0 || state.MysteryStats.CosmicAlignment > 100 {
		t.Errorf("CosmicAlignment out of range: %d", state.MysteryStats.CosmicAlignment)
	}

	// SuspiciousActivity should stay in range
	if state.MysteryStats.SuspiciousActivity < 0 || state.MysteryStats.SuspiciousActivity > 100 {
		t.Errorf("SuspiciousActivity out of range: %d", state.MysteryStats.SuspiciousActivity)
	}

	// Just verify the function ran without error
	_ = initialAlignment
}

func TestGetMysteryStatsDisplay(t *testing.T) {
	state := NewAbsurdState()

	display := state.GetMysteryStatsDisplay()

	if display == "" {
		t.Error("Expected non-empty mystery stats display")
	}

	if !strings.Contains(display, "MYSTERY STATS") {
		t.Error("Display should contain 'MYSTERY STATS'")
	}

	if !strings.Contains(display, "Suspicious") {
		t.Error("Display should show Suspicious stat")
	}
}

func TestGetFearDisplay(t *testing.T) {
	state := NewAbsurdState()

	display := state.GetFearDisplay()

	if display == "" {
		t.Error("Expected non-empty fear display")
	}

	if !strings.Contains(display, "FEARS") {
		t.Error("Display should contain 'FEARS'")
	}
}

func TestShouldShowThought(t *testing.T) {
	state := NewAbsurdState()

	// Run multiple times to test probability
	shownCount := 0
	for i := 0; i < 1000; i++ {
		if state.ShouldShowThought() {
			shownCount++
		}
	}

	// With 15% probability, expect roughly 150 in 1000 (allow some variance)
	if shownCount < 50 || shownCount > 250 {
		t.Errorf("ShouldShowThought probability seems off: %d/1000", shownCount)
	}
}
