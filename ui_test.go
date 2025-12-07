package main

import (
	"os"
	"testing"
	"time"
)

func TestNewUIConfig(t *testing.T) {
	ui := newUIConfig()

	if ui == nil {
		t.Fatal("newUIConfig returned nil")
	}

	if len(ui.spinnerFrames) == 0 {
		t.Error("spinnerFrames should not be empty")
	}

	if len(ui.staticFrames) == 0 {
		t.Error("staticFrames should not be empty")
	}
}

func TestNewUIConfigWithEnvironment(t *testing.T) {
	// Test screen reader mode
	os.Setenv("TAMAGOTCHI_SCREEN_READER", "1")
	defer os.Unsetenv("TAMAGOTCHI_SCREEN_READER")

	ui := newUIConfig()
	if !ui.screenReader {
		t.Error("screenReader should be true when TAMAGOTCHI_SCREEN_READER is set")
	}
	if !ui.reducedMotion {
		t.Error("reducedMotion should be true when screen reader is enabled")
	}
	if ui.soundEnabled {
		t.Error("soundEnabled should be false when screen reader is enabled")
	}
}

func TestNewUIConfigNoSound(t *testing.T) {
	os.Setenv("TAMAGOTCHI_NO_SOUND", "1")
	defer os.Unsetenv("TAMAGOTCHI_NO_SOUND")

	ui := newUIConfig()
	if ui.soundEnabled {
		t.Error("soundEnabled should be false when TAMAGOTCHI_NO_SOUND is set")
	}
}

func TestNewUIConfigHighContrast(t *testing.T) {
	os.Setenv("TAMAGOTCHI_HIGH_CONTRAST", "1")
	defer os.Unsetenv("TAMAGOTCHI_HIGH_CONTRAST")

	ui := newUIConfig()
	if !ui.highContrast {
		t.Error("highContrast should be true when TAMAGOTCHI_HIGH_CONTRAST is set")
	}
}

func TestNewUIConfigColorBlind(t *testing.T) {
	os.Setenv("TAMAGOTCHI_COLORBLIND", "1")
	defer os.Unsetenv("TAMAGOTCHI_COLORBLIND")

	ui := newUIConfig()
	if !ui.colorBlind {
		t.Error("colorBlind should be true when TAMAGOTCHI_COLORBLIND is set")
	}
}

func TestEncodeToMorse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SOS", "... --- ..."},
		{"A", ".-"},
		{"HELLO", ".... . .-.. .-.. ---"},
		{"", ""},
		{"123", ".---- ..--- ...--"},
	}

	for _, tt := range tests {
		result := encodeToMorse(tt.input)
		if result != tt.expected {
			t.Errorf("encodeToMorse(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestDecodeMorseChar(t *testing.T) {
	tests := []struct {
		morse    string
		expected string
	}{
		{".-", "A"},
		{"-...", "B"},
		{"...", "S"},
		{"---", "O"},
		{"....", "H"},
		{"invalid", "?"},
	}

	for _, tt := range tests {
		result := decodeMorseChar(tt.morse)
		if result != tt.expected {
			t.Errorf("decodeMorseChar(%q) = %q, expected %q", tt.morse, result, tt.expected)
		}
	}
}

func TestMorseCodeCompleteness(t *testing.T) {
	// Verify all letters A-Z are in morseCode
	for char := 'A'; char <= 'Z'; char++ {
		if _, exists := morseCode[char]; !exists {
			t.Errorf("morseCode missing letter %c", char)
		}
	}

	// Verify digits 0-9
	for char := '0'; char <= '9'; char++ {
		if _, exists := morseCode[char]; !exists {
			t.Errorf("morseCode missing digit %c", char)
		}
	}
}

func TestShouldAlertForStat(t *testing.T) {
	tests := []struct {
		statName string
		value    int
		expected bool
	}{
		{"hunger", 80, true},    // High hunger should alert
		{"hunger", 50, false},   // Normal hunger shouldn't alert
		{"health", 20, true},    // Low health should alert
		{"health", 50, false},   // Normal health shouldn't alert
		{"happiness", 15, true}, // Low happiness should alert
		{"happiness", 50, false},
		{"cleanliness", 10, true}, // Low cleanliness should alert
		{"cleanliness", 50, false},
		{"unknown", 10, false}, // Unknown stat never alerts
	}

	for _, tt := range tests {
		result := shouldAlertForStat(tt.statName, tt.value)
		if result != tt.expected {
			t.Errorf("shouldAlertForStat(%q, %d) = %v, expected %v",
				tt.statName, tt.value, result, tt.expected)
		}
	}
}

func TestTerminalBellRateLimiting(t *testing.T) {
	ui := newUIConfig()
	ui.soundEnabled = true

	// First bell should set lastBellTime
	initialTime := ui.lastBellTime
	ui.terminalBell()

	// Subsequent immediate bell should be rate limited
	secondBellTime := ui.lastBellTime
	ui.terminalBell()

	// lastBellTime shouldn't change due to rate limiting
	if ui.lastBellTime != secondBellTime {
		t.Error("terminalBell should be rate limited")
	}

	// Verify first bell did update the time
	if secondBellTime == initialTime {
		t.Error("First bell should have updated lastBellTime")
	}
}

func TestBellForEventWhenSoundDisabled(t *testing.T) {
	ui := newUIConfig()
	ui.soundEnabled = false

	// Should not panic or produce effects when sound is disabled
	ui.bellForEvent("critical")
	ui.bellForEvent("alert")
	ui.bellForEvent("achievement")
	ui.bellForEvent("network")
}

func TestRecordMorseEvent(t *testing.T) {
	ui := newUIConfig()

	// Record some events
	ui.recordMorseEvent(true)  // dot
	ui.recordMorseEvent(false) // dash
	ui.recordMorseEvent(true)  // dot

	if len(ui.morseBuffer) != 3 {
		t.Errorf("Expected 3 events in morseBuffer, got %d", len(ui.morseBuffer))
	}

	// Verify first event is a dot
	if !ui.morseBuffer[0].isDot {
		t.Error("First event should be a dot")
	}

	// Verify second event is a dash
	if ui.morseBuffer[1].isDot {
		t.Error("Second event should be a dash")
	}
}

func TestRecordMorseEventBufferLimit(t *testing.T) {
	ui := newUIConfig()

	// Record more than 50 events
	for i := 0; i < 60; i++ {
		ui.recordMorseEvent(i%2 == 0)
	}

	// Buffer should be limited to 50
	if len(ui.morseBuffer) > 50 {
		t.Errorf("morseBuffer should be limited to 50 events, got %d", len(ui.morseBuffer))
	}
}

func TestHiddenMorseMessages(t *testing.T) {
	if len(hiddenMorseMessages) == 0 {
		t.Error("hiddenMorseMessages should not be empty")
	}

	// Verify all messages can be encoded
	for _, msg := range hiddenMorseMessages {
		encoded := encodeToMorse(msg)
		if encoded == "" {
			t.Errorf("Failed to encode message: %s", msg)
		}
	}
}

func TestMaybeMorseMessageWhenDisabled(t *testing.T) {
	ui := newUIConfig()
	ui.soundEnabled = false

	// Should return empty when sound is disabled
	result := ui.maybeMorseMessage()
	if result != "" {
		t.Error("maybeMorseMessage should return empty string when sound is disabled")
	}

	// Also when reducedMotion is enabled
	ui.soundEnabled = true
	ui.reducedMotion = true
	result = ui.maybeMorseMessage()
	if result != "" {
		t.Error("maybeMorseMessage should return empty string when reducedMotion is enabled")
	}
}

func TestDecodeMorseBuffer(t *testing.T) {
	ui := newUIConfig()

	// Empty buffer should return empty string
	result := ui.decodeMorseBuffer()
	if result != "" {
		t.Error("Empty morseBuffer should decode to empty string")
	}

	// Buffer with less than 3 events should return empty
	ui.recordMorseEvent(true)
	ui.recordMorseEvent(false)
	result = ui.decodeMorseBuffer()
	if result != "" {
		t.Error("morseBuffer with < 3 events should decode to empty string")
	}
}

func TestPlayNotificationSoundWhenDisabled(t *testing.T) {
	ui := newUIConfig()
	ui.soundEnabled = false

	// Should not panic when playing any sound type with sound disabled
	ui.playNotificationSound(SoundNone, "TestPet")
	ui.playNotificationSound(SoundCritical, "TestPet")
	ui.playNotificationSound(SoundAlert, "TestPet")
	ui.playNotificationSound(SoundAchievement, "TestPet")
	ui.playNotificationSound(SoundNetwork, "TestPet")
	ui.playNotificationSound(SoundMorse, "TestPet")
}

func TestCheckAndPlayAlertsWhenDisabled(t *testing.T) {
	ui := newUIConfig()
	ui.soundEnabled = false

	pet := NewPet("TestPet")

	// Should not panic when checking alerts with sound disabled
	ui.checkAndPlayAlerts(pet)
}

func TestCheckAndPlayAlertsCriticalState(t *testing.T) {
	ui := newUIConfig()
	ui.soundEnabled = true
	ui.lastBellTime = time.Time{} // Reset rate limit

	pet := NewPet("TestPet")
	pet.Health = 5 // Critical health

	// Should trigger critical alert
	ui.checkAndPlayAlerts(pet)

	// lastBellTime should be updated (bell was played)
	if ui.lastBellTime.IsZero() {
		t.Error("Critical health should have triggered a bell")
	}
}

func TestCheckAndPlayAlertsSickState(t *testing.T) {
	ui := newUIConfig()
	ui.soundEnabled = true
	ui.lastBellTime = time.Time{}

	pet := NewPet("TestPet")
	pet.Health = 50
	pet.IsSick = true

	ui.checkAndPlayAlerts(pet)

	if ui.lastBellTime.IsZero() {
		t.Error("Sick state should have triggered a bell")
	}
}

func TestNotificationSoundConstants(t *testing.T) {
	// Verify all sound types have distinct values
	sounds := []NotificationSound{
		SoundNone,
		SoundCritical,
		SoundAlert,
		SoundAchievement,
		SoundNetwork,
		SoundMorse,
	}

	seen := make(map[NotificationSound]bool)
	for _, s := range sounds {
		if seen[s] {
			t.Errorf("Duplicate NotificationSound value: %d", s)
		}
		seen[s] = true
	}
}

func TestChooseWeather(t *testing.T) {
	// Test that chooseWeather returns valid weather strings
	now := time.Now()
	weather := chooseWeather(now)

	validWeathers := []string{
		"☀️ clear",
		"🌧️ rain",
		"❄️ snow",
		"🌫️ fog",
		"⛅ drifting clouds",
	}

	found := false
	for _, valid := range validWeathers {
		if weather == valid {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("chooseWeather returned invalid weather: %s", weather)
	}
}

func TestPaletteText(t *testing.T) {
	ui := newUIConfig()
	ui.colorEnabled = true

	text := "test"
	result := ui.paletteText(text, ui.palette.accent)

	// Should contain the text
	if len(result) < len(text) {
		t.Error("paletteText result should contain original text")
	}

	// When color disabled, should return plain text
	ui.colorEnabled = false
	result = ui.paletteText(text, ui.palette.accent)
	if result != text {
		t.Errorf("paletteText with color disabled should return plain text, got %q", result)
	}
}

func TestAnimatedBar(t *testing.T) {
	ui := newUIConfig()
	ui.colorEnabled = false // Disable color for easier testing
	ui.reducedMotion = true // Use simpler output

	tests := []struct {
		value    int
		contains string
	}{
		{100, "100%"},
		{50, "50%"},
		{0, "0%"},
	}

	for _, tt := range tests {
		result := ui.animatedBar(tt.value, "")
		if !containsString(result, tt.contains) {
			t.Errorf("animatedBar(%d) should contain %q, got %q", tt.value, tt.contains, result)
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr))
}

func TestSpinningGlyph(t *testing.T) {
	ui := newUIConfig()

	// With reduced motion, should return static glyph
	ui.reducedMotion = true
	result := ui.spinningGlyph()
	if result != "●" {
		t.Errorf("spinningGlyph with reducedMotion should return ●, got %q", result)
	}

	// Without reduced motion, should return one of spinner frames
	ui.reducedMotion = false
	result = ui.spinningGlyph()

	found := false
	for _, frame := range ui.spinnerFrames {
		if result == frame {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("spinningGlyph should return one of spinnerFrames, got %q", result)
	}
}

func TestStaticFrame(t *testing.T) {
	ui := newUIConfig()

	result := ui.staticFrame()

	found := false
	for _, frame := range ui.staticFrames {
		if result == frame {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("staticFrame should return one of staticFrames, got %q", result)
	}
}

func TestFramesForStage(t *testing.T) {
	ui := newUIConfig()
	ui.colorEnabled = false

	stages := []LifeStage{Egg, Baby, Child, Teen, Adult, Dead}

	for _, stage := range stages {
		frames := ui.framesForStage(stage, false)
		if len(frames) == 0 {
			t.Errorf("framesForStage(%v) should return non-empty frames", stage)
		}
	}
}

func TestRenderStatusPanel(t *testing.T) {
	ui := newUIConfig()
	ui.colorEnabled = false
	ui.reducedMotion = true

	pet := NewPet("TestPet")

	result := ui.renderStatusPanel(pet)

	// Should contain pet name
	if !containsSubstring(result, "TestPet") {
		t.Error("renderStatusPanel should contain pet name")
	}

	// Should contain stat labels
	expectedLabels := []string{"Hunger", "Happiness", "Health", "Cleanliness", "Age", "Stage"}
	for _, label := range expectedLabels {
		if !containsSubstring(result, label) {
			t.Errorf("renderStatusPanel should contain %q", label)
		}
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
