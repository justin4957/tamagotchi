package mooc

import (
	"testing"
	"time"
)

func TestNewNetwork(t *testing.T) {
	birthTime := time.Now()
	network := NewNetwork("TestPet", birthTime, "Baby", true)

	if network == nil {
		t.Fatal("NewNetwork should not return nil")
	}

	if network.identity == nil {
		t.Error("Network should have an identity")
	}

	if network.discovery == nil {
		t.Error("Network should have a discovery service")
	}

	if network.gossip == nil {
		t.Error("Network should have a gossip service")
	}

	if network.IsEnabled() {
		t.Error("Network should not be enabled before Start()")
	}
}

func TestLonelyMode(t *testing.T) {
	network := NewNetwork("LonelyPet", time.Now(), "Baby", true)

	network.SetLonelyMode(true)

	if !network.IsLonely() {
		t.Error("Network should be in lonely mode")
	}

	// Start should not enable network in lonely mode
	network.Start()

	if network.IsEnabled() {
		t.Error("Network should not be enabled in lonely mode")
	}
}

func TestGetNetworkStatus(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	// Before start - offline
	status := network.GetNetworkStatus()
	if status != "📡 Network: Offline" {
		t.Errorf("Expected offline status, got: %s", status)
	}

	// In lonely mode
	network.SetLonelyMode(true)
	status = network.GetNetworkStatus()
	if status != "🔇 Network: Disabled (lonely mode)" {
		t.Errorf("Expected lonely mode status, got: %s", status)
	}
}

func TestGetSpookyMessage(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	// Should return empty when no messages
	msg := network.GetSpookyMessage()
	if msg != "" {
		t.Error("Should return empty string when no spooky messages")
	}
}

func TestGetNetworkThought(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	// Should return empty when network not enabled
	thought := network.GetNetworkThought()
	if thought != "" {
		t.Error("Should return empty string when network not enabled")
	}
}

func TestExportImportState(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	// Set some state
	network.state.MemoriesShared = 10
	network.state.DeathsWitnessed = 3
	network.state.Influence = 42
	network.state.NetworkJoinTime = time.Now().Add(-24 * time.Hour)

	// Export
	data, err := network.ExportState()
	if err != nil {
		t.Fatalf("Failed to export state: %v", err)
	}

	// Create new network and import
	network2 := NewNetwork("TestPet2", time.Now(), "Child", true)
	err = network2.ImportState(data)
	if err != nil {
		t.Fatalf("Failed to import state: %v", err)
	}

	if network2.state.MemoriesShared != 10 {
		t.Errorf("Expected MemoriesShared 10, got %d", network2.state.MemoriesShared)
	}

	if network2.state.Influence != 42 {
		t.Errorf("Expected Influence 42, got %d", network2.state.Influence)
	}
}

func TestGetSecretStats(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	// When network is not enabled
	stats := network.GetSecretStats()
	if stats != "Network inactive. The mesh sleeps." {
		t.Errorf("Expected inactive message, got: %s", stats)
	}
}

func TestAnnounceDeath(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Adult", true)

	// Should not panic when network is not enabled
	network.AnnounceDeath("TestPet", 72, "Goodbye world")
}

func TestSetAndGetMood(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	network.SetMood("happy", 80)

	mood, intensity := network.GetMood()
	if mood != "happy" {
		t.Errorf("Expected mood 'happy', got '%s'", mood)
	}
	if intensity != 80 {
		t.Errorf("Expected intensity 80, got %d", intensity)
	}
}

func TestGetFriendCount(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	count := network.GetFriendCount()
	if count != 0 {
		t.Errorf("Expected 0 friends initially, got %d", count)
	}
}

func TestGetOnlineFriendCount(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	count := network.GetOnlineFriendCount()
	if count != 0 {
		t.Errorf("Expected 0 online friends when not enabled, got %d", count)
	}
}

func TestShouldShowNetworkThought(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	// When network is not enabled, should always return false
	for i := 0; i < 100; i++ {
		if network.ShouldShowNetworkThought() {
			t.Error("Should not show network thought when network is disabled")
			break
		}
	}
}

func TestFormatDuration(t *testing.T) {
	network := NewNetwork("TestPet", time.Now(), "Baby", true)

	tests := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Minute, "< 1h"},
		{2 * time.Hour, "2h"},
		{25 * time.Hour, "1d 1h"},
		{72 * time.Hour, "3d 0h"},
	}

	for _, test := range tests {
		result := network.formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %s, expected %s", test.duration, result, test.expected)
		}
	}
}
