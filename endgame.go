package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// EndgameState holds all the absurd endgame progression data
type EndgameState struct {
	// Prestige System
	PrestigeLevel    int    `json:"prestige_level"`
	PrestigeEggColor string `json:"prestige_egg_color"`
	TimesPrestiged   int    `json:"times_prestiged"`

	// Arbitrary Currency
	TamaCoins      int       `json:"tama_coins"` // Can't be spent
	LastLoginBonus time.Time `json:"last_login_bonus"`
	LoginStreak    int       `json:"login_streak"`

	// Achievements
	UnlockedAchievements []string       `json:"unlocked_achievements"`
	AchievementProgress  map[string]int `json:"achievement_progress"`

	// Gacha/Inventory
	InvisibleAccessories []string `json:"invisible_accessories"`
	GachaPulls           int      `json:"gacha_pulls"`

	// Guild
	GuildName   string    `json:"guild_name"`
	GuildRank   string    `json:"guild_rank"`
	GuildJoined time.Time `json:"guild_joined"`

	// Quests
	ActiveQuest     *Quest `json:"active_quest"`
	QuestsCompleted int    `json:"quests_completed"`

	// ARG
	ARGProgress     int       `json:"arg_progress"`
	DiscoveredCodes []string  `json:"discovered_codes"`
	CountdownStart  time.Time `json:"countdown_start"`

	// Social
	FriendCode string `json:"friend_code"`
	ShareCount int    `json:"share_count"`

	// Meta Stats
	TotalPlayTime     time.Duration `json:"total_play_time"`
	SessionStart      time.Time     `json:"-"`
	CommandsEntered   int           `json:"commands_entered"`
	TimesCheckedStats int           `json:"times_checked_stats"`

	// New Game+
	NewGamePlusLevel int  `json:"new_game_plus_level"`
	SpeakInRiddles   bool `json:"speak_in_riddles"`
}

// Quest represents a procedurally generated quest
type Quest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Target      int       `json:"target"`
	Progress    int       `json:"progress"`
	StartTime   time.Time `json:"start_time"`
	Reward      string    `json:"reward"`
}

// Achievement represents an achievement (most are impossible)
type Achievement struct {
	ID          string
	Name        string
	Description string
	Secret      bool
	Impossible  bool
}

// Prestige egg colors
var prestigeColors = []string{
	"Slightly Off-White",
	"Beige (But Different)",
	"Eggshell (Ironic)",
	"Cream Adjacent",
	"Pale Yellow-ish",
	"Almost White",
	"Definitely Not White",
	"Suspicious Ivory",
	"Questionable Ecru",
	"Vaguely Champagne",
}

// Guild name generators
var guildPrefixes = []string{
	"The Order of", "The Brotherhood of", "The Society of",
	"The Collective of", "The Assembly of", "The Union of",
	"The Fellowship of", "The League of", "The Council of",
}

var guildSuffixes = []string{
	"Forgotten Snacks", "Misplaced Keys", "Unread Emails",
	"Awkward Silences", "Lost Socks", "Expired Coupons",
	"Unnecessary Meetings", "Unfinished Projects", "Broken Promises",
	"Abandoned Hobbies", "Missed Connections", "Vague Intentions",
}

// Quest templates
var questTemplates = []struct {
	Name   string
	Desc   string
	Type   string
	Target int
}{
	{"The Waiting Game", "Wait for %d seconds", "wait", 60},
	{"Patience is a Virtue", "Do nothing for %d minutes", "wait", 120},
	{"The Long Pause", "Stare at the screen for %d seconds", "wait", 30},
	{"Contemplative Rest", "Let %d seconds pass in silence", "wait", 90},
	{"The Art of Stillness", "Exist for %d more seconds", "wait", 45},
	{"Temporal Meditation", "Allow %d seconds to flow by", "wait", 75},
	{"The Void Beckons", "Spend %d seconds in contemplation", "wait", 100},
}

// Achievements (including impossible ones)
var allAchievements = []Achievement{
	// Possible achievements
	{ID: "first_feed", Name: "First Meal", Description: "Feed your pet for the first time", Secret: false, Impossible: false},
	{ID: "play_10", Name: "Playful", Description: "Play with your pet 10 times", Secret: false, Impossible: false},
	{ID: "survive_day", Name: "Day One", Description: "Keep your pet alive for 24 hours", Secret: false, Impossible: false},
	{ID: "prestige_1", Name: "Fresh Start", Description: "Prestige for the first time", Secret: false, Impossible: false},
	{ID: "void_gaze", Name: "Void Gazer", Description: "Stare into the void", Secret: false, Impossible: false},
	{ID: "enlightened", Name: "Enlightened One", Description: "Achieve enlightenment", Secret: false, Impossible: false},
	{ID: "guild_join", Name: "Guild Member", Description: "Join a guild", Secret: false, Impossible: false},
	{ID: "quest_complete", Name: "Quest Champion", Description: "Complete a quest", Secret: false, Impossible: false},

	// Secret achievements
	{ID: "debug_mode", Name: "???", Description: "Discover debug mode", Secret: true, Impossible: false},
	{ID: "konami", Name: "Old School", Description: "Enter the code", Secret: true, Impossible: false},
	{ID: "pet_17", Name: "The Number", Description: "Pet your pet exactly 17 times", Secret: true, Impossible: false},
	{ID: "touch_grass", Name: "Touched Grass", Description: "Received the touch grass reminder", Secret: true, Impossible: false},

	// Impossible achievements
	{ID: "impossible_1", Name: "Divide by Zero", Description: "Divide your TamaCoins by zero", Secret: false, Impossible: true},
	{ID: "impossible_2", Name: "Time Traveler", Description: "Play the game yesterday", Secret: false, Impossible: true},
	{ID: "impossible_3", Name: "The Chosen One", Description: "Be selected as Pet of the Day", Secret: false, Impossible: true},
	{ID: "impossible_4", Name: "Infinite Wealth", Description: "Spend your TamaCoins", Secret: false, Impossible: true},
	{ID: "impossible_5", Name: "Social Butterfly", Description: "Have someone actually read your shared pet status", Secret: false, Impossible: true},
	{ID: "impossible_6", Name: "Visible Fashion", Description: "See your invisible accessories", Secret: false, Impossible: true},
	{ID: "impossible_7", Name: "Win the Battle", Description: "Actually win a pet battle", Secret: false, Impossible: true},
	{ID: "impossible_8", Name: "Meaningful Trade", Description: "Trade for something real", Secret: false, Impossible: true},
	{ID: "impossible_9", Name: "Premium User", Description: "Purchase premium features", Secret: false, Impossible: true},
	{ID: "impossible_10", Name: "The End", Description: "Reach the end of the countdown", Secret: false, Impossible: true},
}

// Invisible accessories
var invisibleAccessories = []string{
	"Invisible Top Hat", "Transparent Monocle", "See-Through Cape",
	"Clear Bow Tie", "Invisible Crown", "Transparent Sunglasses",
	"Non-Visible Scarf", "Absent Necklace", "Unseen Earrings",
	"Missing Watch", "Void Bracelet", "Null Ring",
	"Empty Backpack", "Invisible Sword", "Transparent Shield",
}

// NewEndgameState creates a new endgame state
func NewEndgameState() *EndgameState {
	return &EndgameState{
		PrestigeLevel:        0,
		PrestigeEggColor:     "Classic White",
		TamaCoins:            0,
		UnlockedAchievements: make([]string, 0),
		AchievementProgress:  make(map[string]int),
		InvisibleAccessories: make([]string, 0),
		DiscoveredCodes:      make([]string, 0),
		FriendCode:           generateFriendCode(),
		SessionStart:         time.Now(),
		CountdownStart:       time.Now(),
	}
}

// generateFriendCode creates a 47-character friend code
func generateFriendCode() string {
	data := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
	hash := sha256.Sum256([]byte(data))
	code := hex.EncodeToString(hash[:])
	// Format as groups for extra absurdity
	return fmt.Sprintf("%s-%s-%s-%s-%s-%s-%s-%s",
		code[0:6], code[6:12], code[12:18], code[18:24],
		code[24:30], code[30:36], code[36:42], code[42:48])
}

// CheckDailyBonus checks and awards daily login bonus
func (e *EndgameState) CheckDailyBonus() (bool, string) {
	now := time.Now()
	lastBonus := e.LastLoginBonus

	// Check if it's a new day
	if lastBonus.Year() == now.Year() &&
		lastBonus.YearDay() == now.YearDay() {
		return false, ""
	}

	// Check streak
	yesterday := now.AddDate(0, 0, -1)
	if lastBonus.Year() == yesterday.Year() &&
		lastBonus.YearDay() == yesterday.YearDay() {
		e.LoginStreak++
	} else {
		e.LoginStreak = 1
	}

	e.LastLoginBonus = now
	e.TamaCoins++

	return true, fmt.Sprintf(`
╔════════════════════════════════════╗
║      🎁 DAILY LOGIN BONUS! 🎁     ║
╠════════════════════════════════════╣
║ +1 TamaCoin                        ║
║ (Total: %d TamaCoins)              ║
║                                    ║
║ Login Streak: %d days              ║
║                                    ║
║ Note: TamaCoins cannot be spent.   ║
║ They simply exist, like you.       ║
╚════════════════════════════════════╝
`, e.TamaCoins, e.LoginStreak)
}

// GenerateGuildName creates an absurd guild name
func GenerateGuildName() string {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	prefix := guildPrefixes[randomSource.Intn(len(guildPrefixes))]
	suffix := guildSuffixes[randomSource.Intn(len(guildSuffixes))]
	return prefix + " " + suffix
}

// JoinGuild joins a randomly named guild
func (e *EndgameState) JoinGuild() string {
	if e.GuildName != "" {
		return fmt.Sprintf("You're already a member of '%s'.\nYour rank: %s\nLeaving guilds is not implemented.", e.GuildName, e.GuildRank)
	}

	e.GuildName = GenerateGuildName()
	e.GuildRank = "Confused Initiate"
	e.GuildJoined = time.Now()

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      🏰 GUILD JOINED! 🏰          ║
╠════════════════════════════════════╣
║ Welcome to:                        ║
║ "%s"
║                                    ║
║ Your Rank: %s
║                                    ║
║ Guild Benefits:                    ║
║ • None                             ║
║ • Absolutely nothing               ║
║ • A sense of belonging (fake)      ║
╚════════════════════════════════════╝
`, e.GuildName, e.GuildRank)
}

// GenerateQuest creates a new procedural quest
func (e *EndgameState) GenerateQuest() string {
	if e.ActiveQuest != nil {
		return fmt.Sprintf("You already have an active quest:\n%s\nProgress: %d/%d",
			e.ActiveQuest.Name, e.ActiveQuest.Progress, e.ActiveQuest.Target)
	}

	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	template := questTemplates[randomSource.Intn(len(questTemplates))]

	e.ActiveQuest = &Quest{
		Name:        template.Name,
		Description: fmt.Sprintf(template.Desc, template.Target),
		Type:        template.Type,
		Target:      template.Target,
		Progress:    0,
		StartTime:   time.Now(),
		Reward:      "1 TamaCoin (non-spendable)",
	}

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      📜 NEW QUEST! 📜              ║
╠════════════════════════════════════╣
║ %s
║                                    ║
║ Objective:                         ║
║ %s
║                                    ║
║ Reward: %s
╚════════════════════════════════════╝
`, e.ActiveQuest.Name, e.ActiveQuest.Description, e.ActiveQuest.Reward)
}

// UpdateQuest updates quest progress
func (e *EndgameState) UpdateQuest() string {
	if e.ActiveQuest == nil {
		return ""
	}

	elapsed := int(time.Since(e.ActiveQuest.StartTime).Seconds())
	e.ActiveQuest.Progress = elapsed

	if e.ActiveQuest.Progress >= e.ActiveQuest.Target {
		e.QuestsCompleted++
		e.TamaCoins++
		questName := e.ActiveQuest.Name
		e.ActiveQuest = nil

		return fmt.Sprintf(`
╔════════════════════════════════════╗
║      ✅ QUEST COMPLETE! ✅         ║
╠════════════════════════════════════╣
║ "%s" finished!
║                                    ║
║ Reward: +1 TamaCoin                ║
║ (Still can't spend them)           ║
║                                    ║
║ Total Quests Completed: %d         ║
╚════════════════════════════════════╝
`, questName, e.QuestsCompleted)
	}

	return ""
}

// PullGacha does a gacha pull for invisible accessories
func (e *EndgameState) PullGacha() string {
	e.GachaPulls++

	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	accessory := invisibleAccessories[randomSource.Intn(len(invisibleAccessories))]

	// Check for duplicate
	for _, owned := range e.InvisibleAccessories {
		if owned == accessory {
			return fmt.Sprintf(`
╔════════════════════════════════════╗
║      🎰 GACHA RESULT 🎰           ║
╠════════════════════════════════════╣
║ You got: %s
║                                    ║
║ ⚠️ DUPLICATE!                      ║
║ You already own this item.         ║
║ You cannot see it twice.           ║
║                                    ║
║ Total Pulls: %d                    ║
╚════════════════════════════════════╝
`, accessory, e.GachaPulls)
		}
	}

	e.InvisibleAccessories = append(e.InvisibleAccessories, accessory)

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      🎰 GACHA RESULT 🎰           ║
╠════════════════════════════════════╣
║ ✨ NEW ITEM! ✨                    ║
║                                    ║
║ You got: %s
║                                    ║
║ Note: This item is invisible.      ║
║ Your pet is now wearing it.        ║
║ You cannot see it.                 ║
║ But it's there. Trust us.          ║
║                                    ║
║ Total Pulls: %d                    ║
║ Collection: %d/%d                  ║
╚════════════════════════════════════╝
`, accessory, e.GachaPulls, len(e.InvisibleAccessories), len(invisibleAccessories))
}

// StartBattle initiates a pet battle where nothing happens
func (e *EndgameState) StartBattle() string {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	opponentNames := []string{
		"G****y", "F****y", "N*****s", "B***y", "S****w",
		"M****r", "P****t", "C****e", "W*****r", "D***y",
	}
	opponent := opponentNames[randomSource.Intn(len(opponentNames))]

	battleMessages := []string{
		"Both pets stare at each other.",
		"Nothing happens.",
		"The tension is palpable. Or is it?",
		"Your pet blinks. The opponent blinks.",
		"A tumbleweed rolls by.",
		"Both pets declare victory simultaneously.",
		"The battle ends before it begins.",
		"Everyone wins. Everyone also loses.",
	}

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      ⚔️ PET BATTLE! ⚔️            ║
╠════════════════════════════════════╣
║ VS: %s
║                                    ║
║ Battle Log:                        ║
║ > %s
║ > %s
║ > %s
║                                    ║
║ RESULT: TIE (as always)            ║
║                                    ║
║ Both pets have won!                ║
║ Both pets have also lost!          ║
║ Schrödinger would be proud.        ║
╚════════════════════════════════════╝
`,
		opponent,
		battleMessages[randomSource.Intn(len(battleMessages))],
		battleMessages[randomSource.Intn(len(battleMessages))],
		battleMessages[randomSource.Intn(len(battleMessages))],
	)
}

// AttemptTrade tries to trade items that don't exist
func (e *EndgameState) AttemptTrade() string {
	fakeItems := []string{
		"Nothing", "Void Essence", "Empty Promise",
		"Broken Dream", "Lost Potential", "Forgotten Memory",
	}

	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	yourItem := fakeItems[randomSource.Intn(len(fakeItems))]
	theirItem := fakeItems[randomSource.Intn(len(fakeItems))]

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      🔄 TRADE SYSTEM 🔄           ║
╠════════════════════════════════════╣
║ You offer: %s
║ They offer: %s
║                                    ║
║ Trade Status: PENDING              ║
║                                    ║
║ Note: Neither item exists.         ║
║ The trade will never complete.     ║
║ This is by design.                 ║
║                                    ║
║ Thank you for participating in     ║
║ the illusion of economy.           ║
╚════════════════════════════════════╝
`, yourItem, theirItem)
}

// GetCountdownStatus returns the status of the mysterious countdown
func (e *EndgameState) GetCountdownStatus() string {
	// Countdown to... nothing. It resets when it hits zero.
	elapsed := time.Since(e.CountdownStart)
	totalDuration := 7 * 24 * time.Hour // 7 days
	remaining := totalDuration - elapsed

	if remaining <= 0 {
		e.CountdownStart = time.Now()
		remaining = totalDuration
	}

	days := int(remaining.Hours()) / 24
	hours := int(remaining.Hours()) % 24
	minutes := int(remaining.Minutes()) % 60
	seconds := int(remaining.Seconds()) % 60

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      ⏰ THE COUNTDOWN ⏰           ║
╠════════════════════════════════════╣
║                                    ║
║   %dd %02dh %02dm %02ds
║                                    ║
║ Something is coming.               ║
║ Or is it?                          ║
║ No one knows.                      ║
║ (Not even us.)                     ║
║                                    ║
║ When it reaches zero:              ║
║ ???                                ║
╚════════════════════════════════════╝
`, days, hours, minutes, seconds)
}

// GetARGClue generates a cryptic ARG clue
func (e *EndgameState) GetARGClue() string {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate fake coordinates
	lat := 40.0 + randomSource.Float64()*10
	lon := -74.0 + randomSource.Float64()*10

	// Generate base64 message
	messages := []string{
		"THE MESH REMEMBERS",
		"SEVENTEEN IS THE KEY",
		"LOOK BEHIND THE SAVE FILE",
		"THE VOID SPEAKS TRUTH",
		"NOT ALL EGGS ARE EQUAL",
	}
	message := messages[randomSource.Intn(len(messages))]
	encoded := base64.StdEncoding.EncodeToString([]byte(message))

	e.ARGProgress++

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      🔮 MYSTERIOUS CLUE 🔮        ║
╠════════════════════════════════════╣
║                                    ║
║ Coordinates: %.4f, %.4f
║                                    ║
║ Encoded Message:                   ║
║ %s
║                                    ║
║ What does it mean?                 ║
║ We don't know either.              ║
║                                    ║
║ ARG Progress: %d/∞                 ║
╚════════════════════════════════════╝
`, lat, lon, encoded, e.ARGProgress)
}

// GenerateShareText creates absurdly long shareable text
func (e *EndgameState) GenerateShareText(petName string, petStage string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05 MST")

	return fmt.Sprintf(`
🎮 TAMAGOTCHI STATUS UPDATE 🎮
================================
Pet: %s
Stage: %s
TamaCoins: %d (non-spendable)
Prestige: Level %d (%s egg)
Guild: %s
Invisible Accessories: %d equipped (you can't see them)
Quests Completed: %d
Gacha Pulls: %d
ARG Progress: %d/∞
Friend Code: %s

📊 Arbitrary Statistics:
Commands Entered: %d
Times Checked Stats: %d
Login Streak: %d days

🕐 Timestamp: %s

#Tamagotchi #VirtualPet #TamaCoins #Gaming #Meaningless #TimeWasted

If you're reading this, I'm impressed and also concerned.
This was auto-generated and serves no purpose.
Thank you for your attention.
================================
`, petName, petStage, e.TamaCoins, e.PrestigeLevel, e.PrestigeEggColor,
		e.GuildName, len(e.InvisibleAccessories), e.QuestsCompleted,
		e.GachaPulls, e.ARGProgress, e.FriendCode,
		e.CommandsEntered, e.TimesCheckedStats, e.LoginStreak, timestamp)
}

// ShowPremiumOffer shows the fake premium content
func ShowPremiumOffer() string {
	return `
╔════════════════════════════════════╗
║      💎 PREMIUM CONTENT 💎        ║
╠════════════════════════════════════╣
║                                    ║
║  TAMAGOTCHI PREMIUM™               ║
║  Price: N/A                        ║
║                                    ║
║  Features:                         ║
║  • Nothing additional              ║
║  • Same experience as free         ║
║  • A sense of superiority (fake)   ║
║  • Golden TamaCoins (still useless)║
║                                    ║
║  "Premium is a state of mind."     ║
║        - Ancient Proverb           ║
║                                    ║
║  Purchase Options:                 ║
║  [Not Implemented]                 ║
║                                    ║
║  This message was brought to you   ║
║  by the concept of capitalism.     ║
║                                    ║
╚════════════════════════════════════╝

A Brief Essay on Digital Ownership:

In the age of digital goods, what does it mean
to "own" something you cannot touch? These
invisible accessories you've collected - are
they truly yours? Or are they merely entries
in a JSON file, ephemeral as morning dew?

The TamaCoins you've accumulated cannot be
spent. This is not a bug, but a feature - a
meditation on the nature of value itself.
What is currency without exchange? What is
wealth without spending?

Perhaps the real premium content was the
time we wasted along the way.

Thank you for attending this TED talk.
`
}

// ShowFakeAd shows a fake advertisement
func ShowFakeAd() string {
	return `
╔════════════════════════════════════╗
║      📺 ADVERTISEMENT 📺          ║
╠════════════════════════════════════╣
║                                    ║
║  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  ║
║  ░                              ░  ║
║  ░   BUY NOTHING TODAY!         ░  ║
║  ░                              ░  ║
║  ░   Limited Time: Forever      ░  ║
║  ░   Price: $0.00               ░  ║
║  ░   Value: Priceless           ░  ║
║  ░                              ░  ║
║  ░   Click Here: [No Link]      ░  ║
║  ░                              ░  ║
║  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  ║
║                                    ║
║  Thank you for watching!           ║
║  Reward: Satisfaction of waiting   ║
║                                    ║
╚════════════════════════════════════╝
`
}

// GetMetaStats returns absurd meta statistics
func (e *EndgameState) GetMetaStats() string {
	e.TimesCheckedStats++

	sessionDuration := time.Since(e.SessionStart)
	totalTime := e.TotalPlayTime + sessionDuration

	hours := int(totalTime.Hours())
	minutes := int(totalTime.Minutes()) % 60

	// Estimate time wasted
	wastedPercentage := 100.0 // All of it

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      📊 META STATISTICS 📊        ║
╠════════════════════════════════════╣
║                                    ║
║ Time Investment:                   ║
║ • This Session: %s
║ • Total Playtime: %dh %dm
║ • Time Wasted: %.1f%%
║                                    ║
║ Engagement Metrics:                ║
║ • Commands Entered: %d
║ • Stats Checked: %d times
║ • Achievements: %d/%d
║ • Quests Done: %d
║ • Gacha Pulls: %d
║                                    ║
║ Economic Status:                   ║
║ • TamaCoins: %d
║ • Spending Power: $0.00
║                                    ║
║ Existential Status:                ║
║ • Meaning Found: No                ║
║ • Regrets: Calculating...          ║
║                                    ║
╚════════════════════════════════════╝
`,
		formatDuration(sessionDuration),
		hours, minutes,
		wastedPercentage,
		e.CommandsEntered,
		e.TimesCheckedStats,
		len(e.UnlockedAchievements), len(allAchievements),
		e.QuestsCompleted,
		e.GachaPulls,
		e.TamaCoins,
	)
}

// formatDuration formats a duration nicely
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// CheckTouchGrass checks if user should be reminded to touch grass
func (e *EndgameState) CheckTouchGrass() (bool, string) {
	sessionDuration := time.Since(e.SessionStart)

	if sessionDuration >= 4*time.Hour {
		return true, `
╔════════════════════════════════════╗
║      🌿 GENTLE REMINDER 🌿        ║
╠════════════════════════════════════╣
║                                    ║
║  You've been playing for over      ║
║  4 hours.                          ║
║                                    ║
║  Have you considered:              ║
║  • Going outside                   ║
║  • Touching grass                  ║
║  • Feeling the sun                 ║
║  • Questioning your choices        ║
║                                    ║
║  Your pet is concerned.            ║
║  We are also concerned.            ║
║                                    ║
║  (This message unlocks the         ║
║   "Touched Grass" achievement      ║
║   ironically)                      ║
║                                    ║
╚════════════════════════════════════╝
`
	}

	return false, ""
}

// UnlockAchievement unlocks an achievement
func (e *EndgameState) UnlockAchievement(id string) (bool, string) {
	// Check if already unlocked
	for _, achieved := range e.UnlockedAchievements {
		if achieved == id {
			return false, ""
		}
	}

	// Find achievement
	for _, ach := range allAchievements {
		if ach.ID == id {
			if ach.Impossible {
				return false, "" // Can't unlock impossible achievements
			}

			e.UnlockedAchievements = append(e.UnlockedAchievements, id)
			return true, fmt.Sprintf(`
╔════════════════════════════════════╗
║      🏆 ACHIEVEMENT UNLOCKED! 🏆  ║
╠════════════════════════════════════╣
║                                    ║
║  %s
║  "%s"
║                                    ║
║  Progress: %d/%d achievements
╚════════════════════════════════════╝
`, ach.Name, ach.Description, len(e.UnlockedAchievements), len(allAchievements))
		}
	}

	return false, ""
}

// ShowAchievements displays all achievements
func (e *EndgameState) ShowAchievements() string {
	var builder strings.Builder

	builder.WriteString("\n╔════════════════════════════════════╗\n")
	builder.WriteString("║      🏆 ACHIEVEMENTS 🏆           ║\n")
	builder.WriteString("╠════════════════════════════════════╣\n")

	unlocked := make(map[string]bool)
	for _, id := range e.UnlockedAchievements {
		unlocked[id] = true
	}

	for _, ach := range allAchievements {
		status := "❌"
		if unlocked[ach.ID] {
			status = "✅"
		}

		name := ach.Name
		desc := ach.Description
		if ach.Secret && !unlocked[ach.ID] {
			name = "???"
			desc = "Secret achievement"
		}
		if ach.Impossible {
			desc += " (IMPOSSIBLE)"
		}

		builder.WriteString(fmt.Sprintf("║ %s %s\n", status, name))
		builder.WriteString(fmt.Sprintf("║    %s\n", desc))
	}

	builder.WriteString(fmt.Sprintf("║\n║ Total: %d/%d\n", len(e.UnlockedAchievements), len(allAchievements)))
	builder.WriteString("╚════════════════════════════════════╝\n")

	return builder.String()
}

// ShowLeaderboard shows a fake leaderboard
func (e *EndgameState) ShowLeaderboard() string {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	metrics := []string{
		"TamaCoins Hoarded", "Invisible Items Worn", "Void Gazes",
		"Meaningless Clicks", "Existential Crises", "Time Wasted (seconds)",
		"Arbitrary Points", "Cosmic Alignment", "Vibe Score",
	}

	metric := metrics[randomSource.Intn(len(metrics))]

	// Generate fake players
	fakeNames := []string{
		"xX_VoidMaster_Xx", "TamaPro2024", "EggLord420",
		"PetWhisperer", "Definitely_Not_A_Bot", "GrindNeverStops",
	}

	var builder strings.Builder
	builder.WriteString("\n╔════════════════════════════════════╗\n")
	builder.WriteString("║      🏅 LEADERBOARD 🏅            ║\n")
	builder.WriteString(fmt.Sprintf("║  Today's Metric: %s\n", metric))
	builder.WriteString("╠════════════════════════════════════╣\n")

	for i := 0; i < 5; i++ {
		score := 10000 - (i * 1000) + randomSource.Intn(500)
		name := fakeNames[i]
		builder.WriteString(fmt.Sprintf("║ #%d %s: %d\n", i+1, name, score))
	}

	// Player is always #6
	builder.WriteString(fmt.Sprintf("║ ...\n"))
	builder.WriteString(fmt.Sprintf("║ #6 You: %d\n", e.TamaCoins))
	builder.WriteString("║\n")
	builder.WriteString("║ Note: Leaderboard metric changes\n")
	builder.WriteString("║ daily for no reason.\n")
	builder.WriteString("╚════════════════════════════════════╝\n")

	return builder.String()
}

// IncrementCommand tracks command usage
func (e *EndgameState) IncrementCommand() {
	e.CommandsEntered++
}

// UpdatePlayTime updates the total play time
func (e *EndgameState) UpdatePlayTime() {
	e.TotalPlayTime += time.Since(e.SessionStart)
	e.SessionStart = time.Now()
}
