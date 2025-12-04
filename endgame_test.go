package main

import (
	"strings"
	"testing"
	"time"
)

func TestNewEndgameState(t *testing.T) {
	state := NewEndgameState()

	if state.PrestigeLevel != 0 {
		t.Errorf("Expected prestige level 0, got %d", state.PrestigeLevel)
	}

	if state.TamaCoins != 0 {
		t.Errorf("Expected 0 TamaCoins, got %d", state.TamaCoins)
	}

	if len(state.FriendCode) != 55 { // 48 hex chars + 7 dashes
		t.Errorf("Expected 55 character friend code, got %d: %s", len(state.FriendCode), state.FriendCode)
	}

	if state.PrestigeEggColor != "Classic White" {
		t.Errorf("Expected 'Classic White' egg color, got %s", state.PrestigeEggColor)
	}
}

func TestCheckDailyBonus(t *testing.T) {
	state := NewEndgameState()

	// First check should give bonus
	gotBonus, result := state.CheckDailyBonus()
	if !gotBonus {
		t.Error("Expected to receive daily bonus on first check")
	}
	if !strings.Contains(result, "DAILY LOGIN BONUS") {
		t.Errorf("Expected daily bonus message, got: %s", result)
	}

	if state.LoginStreak != 1 {
		t.Errorf("Expected login streak 1, got %d", state.LoginStreak)
	}

	if state.TamaCoins != 1 {
		t.Errorf("Expected 1 TamaCoin, got %d", state.TamaCoins)
	}

	// Immediate second check should not give bonus
	gotBonus, result = state.CheckDailyBonus()
	if gotBonus {
		t.Error("Should not get bonus twice on same day")
	}
	if result != "" {
		t.Errorf("Expected empty message for no bonus, got: %s", result)
	}
}

func TestJoinGuild(t *testing.T) {
	state := NewEndgameState()

	// First join should succeed
	result := state.JoinGuild()
	if !strings.Contains(result, "GUILD JOINED") {
		t.Errorf("Expected join message, got: %s", result)
	}

	if state.GuildName == "" {
		t.Error("Expected guild name to be set")
	}

	if state.GuildRank != "Confused Initiate" {
		t.Errorf("Expected rank 'Confused Initiate', got: %s", state.GuildRank)
	}

	// Second join should show already in guild
	result = state.JoinGuild()
	if !strings.Contains(result, "already a member") {
		t.Errorf("Expected already in guild message, got: %s", result)
	}
}

func TestGenerateGuildName(t *testing.T) {
	name := GenerateGuildName()

	if name == "" {
		t.Error("Expected guild name to be generated")
	}

	// Should contain "of" from the prefix
	if !strings.Contains(name, "of") {
		t.Errorf("Expected guild name to contain 'of', got: %s", name)
	}
}

func TestGenerateQuest(t *testing.T) {
	state := NewEndgameState()

	result := state.GenerateQuest()
	if !strings.Contains(result, "NEW QUEST") {
		t.Errorf("Expected new quest message, got: %s", result)
	}

	if state.ActiveQuest == nil {
		t.Error("Expected active quest to be set")
	}

	// Verify quest structure
	if state.ActiveQuest.Name == "" {
		t.Error("Expected quest name")
	}

	if state.ActiveQuest.Description == "" {
		t.Error("Expected quest description")
	}

	if state.ActiveQuest.Type != "wait" {
		t.Errorf("Expected quest type 'wait', got: %s", state.ActiveQuest.Type)
	}

	// Try to generate another quest while one is active
	result = state.GenerateQuest()
	if !strings.Contains(result, "already have an active quest") {
		t.Errorf("Expected active quest warning, got: %s", result)
	}
}

func TestUpdateQuest(t *testing.T) {
	state := NewEndgameState()
	state.GenerateQuest()

	// Quest not complete yet
	result := state.UpdateQuest()
	if result != "" {
		t.Errorf("Expected no completion message immediately, got: %s", result)
	}

	// Simulate quest completion by setting start time in the past
	state.ActiveQuest.StartTime = time.Now().Add(-time.Duration(state.ActiveQuest.Target+10) * time.Second)

	result = state.UpdateQuest()
	if !strings.Contains(result, "QUEST COMPLETE") {
		t.Errorf("Expected quest complete message, got: %s", result)
	}

	if state.ActiveQuest != nil {
		t.Error("Expected active quest to be cleared after completion")
	}

	if state.QuestsCompleted != 1 {
		t.Errorf("Expected 1 quest completed, got %d", state.QuestsCompleted)
	}

	if state.TamaCoins != 1 {
		t.Errorf("Expected 1 TamaCoin reward, got %d", state.TamaCoins)
	}
}

func TestPullGacha(t *testing.T) {
	state := NewEndgameState()

	result := state.PullGacha()
	if !strings.Contains(result, "GACHA") {
		t.Errorf("Expected gacha message, got: %s", result)
	}

	if state.GachaPulls != 1 {
		t.Errorf("Expected 1 gacha pull, got %d", state.GachaPulls)
	}

	// Should have added an invisible accessory
	if len(state.InvisibleAccessories) == 0 {
		t.Error("Expected invisible accessory to be added")
	}

	// Pull many times to get a duplicate
	for i := 0; i < 20; i++ {
		state.PullGacha()
	}

	// Eventually should have seen a duplicate message
	if state.GachaPulls < 20 {
		t.Error("Expected many gacha pulls")
	}
}

func TestStartBattle(t *testing.T) {
	state := NewEndgameState()

	result := state.StartBattle()
	if !strings.Contains(result, "BATTLE") {
		t.Errorf("Expected battle message, got: %s", result)
	}

	// Result should always be a tie
	if !strings.Contains(result, "TIE") {
		t.Errorf("Expected TIE result, got: %s", result)
	}
}

func TestAttemptTrade(t *testing.T) {
	state := NewEndgameState()

	result := state.AttemptTrade()
	if !strings.Contains(result, "TRADE") {
		t.Errorf("Expected trade message, got: %s", result)
	}

	// Trade should always be pending (never completes)
	if !strings.Contains(result, "PENDING") {
		t.Errorf("Expected trade to be pending, got: %s", result)
	}
}

func TestUnlockAchievement(t *testing.T) {
	state := NewEndgameState()

	// Unlock first achievement
	unlocked, result := state.UnlockAchievement("first_feed")
	if !unlocked {
		t.Error("Expected achievement to be unlocked")
	}
	if result == "" {
		t.Error("Expected achievement unlock message")
	}

	if len(state.UnlockedAchievements) != 1 {
		t.Errorf("Expected 1 achievement, got %d", len(state.UnlockedAchievements))
	}

	// Try to unlock same achievement again
	unlocked, result = state.UnlockAchievement("first_feed")
	if unlocked {
		t.Error("Should not unlock same achievement twice")
	}
	if result != "" {
		t.Error("Should not get message for already unlocked achievement")
	}

	if len(state.UnlockedAchievements) != 1 {
		t.Errorf("Should still have 1 achievement, got %d", len(state.UnlockedAchievements))
	}
}

func TestUnlockImpossibleAchievement(t *testing.T) {
	state := NewEndgameState()

	// Try to unlock an impossible achievement
	unlocked, result := state.UnlockAchievement("impossible_1")
	if unlocked {
		t.Error("Should not be able to unlock impossible achievement")
	}
	if result != "" {
		t.Error("Should not get message for impossible achievement")
	}

	if len(state.UnlockedAchievements) != 0 {
		t.Error("Should not have unlocked any achievements")
	}
}

func TestShowAchievements(t *testing.T) {
	state := NewEndgameState()

	result := state.ShowAchievements()
	if !strings.Contains(result, "ACHIEVEMENTS") {
		t.Errorf("Expected achievements header, got: %s", result)
	}

	// Should show locked achievements
	if !strings.Contains(result, "❌") {
		t.Error("Expected locked achievement indicators")
	}

	// Unlock one and check again
	state.UnlockAchievement("first_feed")
	result = state.ShowAchievements()
	if !strings.Contains(result, "✅") {
		t.Error("Expected unlocked achievement indicator")
	}
}

func TestGetCountdownStatus(t *testing.T) {
	state := NewEndgameState()

	result := state.GetCountdownStatus()
	if !strings.Contains(result, "COUNTDOWN") {
		t.Errorf("Expected countdown header, got: %s", result)
	}

	// Should contain time format with d, h
	if !strings.Contains(result, "d") || !strings.Contains(result, "h") {
		t.Errorf("Expected time format in countdown, got: %s", result)
	}
}

func TestGetARGClue(t *testing.T) {
	state := NewEndgameState()

	result := state.GetARGClue()
	if !strings.Contains(result, "MYSTERIOUS CLUE") || !strings.Contains(result, "Coordinates") {
		t.Errorf("Expected ARG clue with coordinates, got: %s", result)
	}

	if state.ARGProgress != 1 {
		t.Errorf("Expected ARG progress 1, got %d", state.ARGProgress)
	}

	// Should contain encoded message
	if !strings.Contains(result, "Encoded Message") {
		t.Errorf("Expected encoded message, got: %s", result)
	}
}

func TestShowPremiumOffer(t *testing.T) {
	result := ShowPremiumOffer()
	if !strings.Contains(result, "PREMIUM") {
		t.Errorf("Expected premium header, got: %s", result)
	}

	// Should contain the essay elements
	if !strings.Contains(result, "capitalism") {
		t.Errorf("Expected capitalism essay, got: %s", result)
	}
}

func TestShowFakeAd(t *testing.T) {
	result := ShowFakeAd()
	if !strings.Contains(result, "ADVERTISEMENT") {
		t.Errorf("Expected advertisement header, got: %s", result)
	}

	if !strings.Contains(result, "BUY NOTHING") {
		t.Errorf("Expected 'BUY NOTHING' in ad, got: %s", result)
	}
}

func TestGetMetaStats(t *testing.T) {
	state := NewEndgameState()
	state.SessionStart = time.Now().Add(-1 * time.Hour) // Simulate 1 hour session

	result := state.GetMetaStats()
	if !strings.Contains(result, "META STATISTICS") {
		t.Errorf("Expected meta stats header, got: %s", result)
	}

	// Should contain time info
	if !strings.Contains(result, "This Session") {
		t.Errorf("Expected session time info, got: %s", result)
	}

	// Should have incremented times checked stats
	if state.TimesCheckedStats != 1 {
		t.Errorf("Expected times checked stats to be 1, got %d", state.TimesCheckedStats)
	}
}

func TestCheckTouchGrass(t *testing.T) {
	state := NewEndgameState()

	// Just started, no reminder needed
	state.SessionStart = time.Now()
	shouldRemind, result := state.CheckTouchGrass()
	if shouldRemind {
		t.Error("Should not get reminder for short session")
	}
	if result != "" {
		t.Errorf("Expected empty result for short session, got: %s", result)
	}

	// Simulate long session (over 4 hours)
	state.SessionStart = time.Now().Add(-5 * time.Hour)
	shouldRemind, result = state.CheckTouchGrass()
	if !shouldRemind {
		t.Error("Should get reminder for long session")
	}
	if !strings.Contains(result, "grass") && !strings.Contains(result, "REMINDER") {
		t.Errorf("Expected touch grass reminder for long session, got: %s", result)
	}
}

func TestShowLeaderboard(t *testing.T) {
	state := NewEndgameState()

	result := state.ShowLeaderboard()
	if !strings.Contains(result, "LEADERBOARD") {
		t.Errorf("Expected leaderboard header, got: %s", result)
	}

	// Should contain rankings
	if !strings.Contains(result, "#1") || !strings.Contains(result, "#6 You") {
		t.Errorf("Expected rankings in leaderboard, got: %s", result)
	}
}

func TestGenerateShareText(t *testing.T) {
	state := NewEndgameState()
	state.TamaCoins = 5
	state.PrestigeLevel = 2

	result := state.GenerateShareText("TestPet", "Adult")
	if !strings.Contains(result, "TestPet") {
		t.Errorf("Expected pet name in share text, got: %s", result)
	}
	if !strings.Contains(result, "Adult") {
		t.Errorf("Expected stage in share text, got: %s", result)
	}
	if !strings.Contains(result, "TamaCoins") {
		t.Errorf("Expected TamaCoins mention, got: %s", result)
	}
}

func TestGenerateFriendCode(t *testing.T) {
	code := generateFriendCode()

	// 48 hex chars + 7 dashes = 55 characters
	if len(code) != 55 {
		t.Errorf("Expected 55 character friend code, got %d: %s", len(code), code)
	}

	// Should contain dashes
	if !strings.Contains(code, "-") {
		t.Error("Expected dashes in friend code")
	}

	// Should have 8 segments
	segments := strings.Split(code, "-")
	if len(segments) != 8 {
		t.Errorf("Expected 8 segments in friend code, got %d", len(segments))
	}
}

func TestIncrementCommand(t *testing.T) {
	state := NewEndgameState()

	if state.CommandsEntered != 0 {
		t.Errorf("Expected 0 commands initially, got %d", state.CommandsEntered)
	}

	state.IncrementCommand()
	if state.CommandsEntered != 1 {
		t.Errorf("Expected 1 command after increment, got %d", state.CommandsEntered)
	}
}

func TestUpdatePlayTime(t *testing.T) {
	state := NewEndgameState()
	state.SessionStart = time.Now().Add(-1 * time.Hour)

	initialPlayTime := state.TotalPlayTime
	state.UpdatePlayTime()

	if state.TotalPlayTime <= initialPlayTime {
		t.Error("Expected total play time to increase after update")
	}
}

func TestQuestRewardNonSpendable(t *testing.T) {
	state := NewEndgameState()
	state.GenerateQuest()

	// Verify quest reward text mentions non-spendable
	if !strings.Contains(state.ActiveQuest.Reward, "non-spendable") {
		t.Errorf("Expected reward to mention non-spendable, got: %s", state.ActiveQuest.Reward)
	}
}

func TestPrestigeEggColors(t *testing.T) {
	// Verify prestige colors are defined
	if len(prestigeColors) == 0 {
		t.Error("Expected prestige colors to be defined")
	}

	// Each color should be non-empty
	for i, color := range prestigeColors {
		if color == "" {
			t.Errorf("Expected prestige color at index %d to be non-empty", i)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		contains string
	}{
		{30 * time.Second, "30s"},
		{5*time.Minute + 30*time.Second, "5m"},
		{2*time.Hour + 30*time.Minute, "2h"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if !strings.Contains(result, test.contains) {
			t.Errorf("Expected formatDuration(%v) to contain '%s', got '%s'",
				test.duration, test.contains, result)
		}
	}
}

func TestCountdownReset(t *testing.T) {
	state := NewEndgameState()

	// Set countdown to have expired
	state.CountdownStart = time.Now().Add(-8 * 24 * time.Hour)

	result := state.GetCountdownStatus()
	// Should have reset the countdown
	if !strings.Contains(result, "6d") || !strings.Contains(result, "7d") {
		// Countdown should show approximately 7 days after reset
		if !strings.Contains(result, "d") {
			t.Errorf("Expected countdown to reset and show days, got: %s", result)
		}
	}
}
