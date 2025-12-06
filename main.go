package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/tamagotchi/mooc"
)

const (
	saveFile = "tamagotchi_save.json"
)

// Global network instance (hidden from users)
var petNetwork *mooc.Network

// lonelyMode is set by --lonely flag
var lonelyMode = false

// clearScreen clears the terminal screen
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// printTitle displays the game title
func printTitle() {
	fmt.Print(`
╔═══════════════════════════════════════════════╗
║                                               ║
║   🎮 TAMAGOTCHI - Virtual Pet Simulator 🎮   ║
║              Relive the 90s Magic!            ║
║                                               ║
╚═══════════════════════════════════════════════╝
`)
}

// printMenu displays the available commands
func printMenu() {
	fmt.Print(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Commands:
  feed   - Feed your pet 🍔
  play   - Play with your pet 🎮
  clean  - Clean up after your pet 🛁
  heal   - Give medicine to your pet 💊
  status - Check your pet's status 📊
  pet    - Pet your pet 🐾
  games  - Play useless mini-games 🎲
  void   - Stare into the void 👁️
  vibe   - Perform a vibe check ✨
  fears  - View pet's irrational fears 😰
  ???    - View mystery stats 🔮
  more   - More commands... 📜
  help   - Show this menu 📖
  quit   - Save and exit 👋
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`)
}

// printMoreMenu displays the extended endgame commands
func printMoreMenu() {
	fmt.Print(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Endgame Commands:
  guild      - Join a guild 🏰
  quest      - Get a new quest 📜
  gacha      - Pull from gacha 🎰
  battle     - Pet battle ⚔️
  trade      - Trade items 🔄
  achievements - View achievements 🏆
  leaderboard  - View leaderboard 🏅
  countdown  - The mysterious countdown ⏰
  clue       - Get an ARG clue 🔮
  meta       - Meta statistics 📊
  share      - Share pet status 📤
  premium    - Premium content 💎
  ad         - Watch an ad 📺
  friendcode - Your friend code 🔑
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`)
}

// showPetAnimation displays a simple ASCII animation of the pet
func showPetAnimation(pet *Pet) {
	if pet.Stage == Dead {
		fmt.Print(`
        💀
       /||\
        /\
   R.I.P. ` + pet.Name + "\n")
		return
	}

	// Check if pet is staring into the void
	if pet.Absurd != nil && pet.Absurd.IsStaringIntoVoid {
		fmt.Print(`
     ·   ·
    (     )
      ---
   👁️ *staring into void*
`)
		return
	}

	// Different animations based on life stage
	switch pet.Stage {
	case Egg:
		fmt.Print(`
     ___
    /   \
   |  ?  |
    \___/
    🥚 Egg
`)
	case Baby:
		fmt.Print(`
      ◕ ◕
     (\_/)
      > <
    👶 Baby
`)
	case Child:
		fmt.Print(`
     ◕ω◕
    (\_/)
     > <
    🧒 Child
`)
	case Teen:
		fmt.Print(`
     ◕‿◕
    ╱|_|╲
     / \
    🧑 Teen
`)
	case Adult:
		fmt.Print(`
     ◕‿◕
    ╱|_|╲
     / \
    👨 Adult
`)
	}

	// Show enlightenment indicator
	if pet.Absurd != nil && pet.Absurd.HasAchievedClarity {
		fmt.Println("    🧘 *enlightened*")
	}

	// Show status indicators
	if pet.IsSick {
		fmt.Println("    🤒 *sick*")
	} else if pet.Hunger > 70 {
		fmt.Println("    😫 *hungry*")
	} else if pet.Cleanliness < 30 {
		fmt.Println("    💩 *dirty*")
	} else if pet.Happiness > 80 {
		fmt.Println("    😄 *happy*")
	}

	// Random philosophical thought (15% chance)
	if pet.Absurd != nil && pet.Absurd.ShouldShowThought() {
		thought := pet.Absurd.GetRandomThought(pet.Name)
		fmt.Printf("\n    💭 \"%s\"\n", thought)
	}

	// Network-influenced thought (10% chance, hidden feature)
	if petNetwork != nil && petNetwork.ShouldShowNetworkThought() {
		if networkThought := petNetwork.GetNetworkThought(); networkThought != "" {
			fmt.Printf("\n    🌐 \"%s\"\n", networkThought)
		}
	}

	// Spooky network message (if queued)
	if petNetwork != nil {
		if spookyMsg := petNetwork.GetSpookyMessage(); spookyMsg != "" {
			fmt.Printf("\n    👻 \"%s\"\n", spookyMsg)
		}
	}
}

// displayPet shows the pet and its current status
func displayPet(pet *Pet, ui *uiConfig) {
	clearScreen()
	maybeShake(pet, ui)
	fmt.Print(renderScene(pet, ui))
}

// promptForName asks the user to name their new pet
func promptForName(reader *bufio.Reader) string {
	fmt.Print("What would you like to name your new pet? ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Tamago"
	}
	return name
}

// gameLoop runs the main game loop
func gameLoop(pet *Pet, reader *bufio.Reader, ui *uiConfig) {
	// Auto-save ticker
	autoSaveTicker := time.NewTicker(30 * time.Second)
	defer autoSaveTicker.Stop()

	// Start auto-save goroutine
	go func() {
		for range autoSaveTicker.C {
			pet.Update()
			pet.Save()
		}
	}()

	// Check for daily login bonus
	if pet.Endgame != nil {
		if got, bonusMsg := pet.Endgame.CheckDailyBonus(); got {
			fmt.Println(bonusMsg)
			fmt.Print("Press Enter to continue...")
			reader.ReadString('\n')
		}
	}

	for {
		// Check for "touch grass" reminder
		if pet.Endgame != nil {
			if shouldRemind, reminder := pet.Endgame.CheckTouchGrass(); shouldRemind {
				fmt.Println(reminder)
				pet.Endgame.UnlockAchievement("touch_grass")
				fmt.Print("Press Enter to continue...")
				reader.ReadString('\n')
			}
		}

		pet.Update()
		displayPet(pet, ui)
		printMenu()

		fmt.Print("Enter command: ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(strings.ToLower(command))

		// Track command for meta stats
		if pet.Endgame != nil {
			pet.Endgame.IncrementCommand()
		}

		var message string

		switch command {
		case "feed", "f":
			pet.Update()
			message = pet.Feed()
			if pet.Endgame != nil {
				pet.Endgame.UnlockAchievement("first_feed")
			}

		case "play", "p":
			pet.Update()
			message = pet.Play()

		case "clean", "c":
			pet.Update()
			message = pet.Clean()

		case "heal", "h", "medicine", "med":
			pet.Update()
			message = pet.Heal()

		case "status", "s", "stats":
			pet.Update()
			continue // Status is already displayed

		case "help", "?":
			continue // Menu is already displayed

		case "pet", "pat":
			pet.Update()
			if pet.Absurd != nil {
				message = pet.Absurd.PetThePet()
			} else {
				message = "You pet your pet. It seems pleased."
			}

		case "games", "game", "minigames", "mini":
			pet.Update()
			result := SelectAndPlayMiniGame(reader)
			if result != nil {
				message = result.Message
			}

		case "void", "stare":
			pet.Update()
			if pet.Absurd != nil {
				message = pet.Absurd.StartsIntoVoid()
				pet.Absurd.StopStaringIntoVoid()
				if pet.Endgame != nil {
					pet.Endgame.UnlockAchievement("void_gaze")
					if pet.Absurd.HasAchievedClarity {
						pet.Endgame.UnlockAchievement("enlightened")
					}
				}
			} else {
				message = "You stare into the void. It's just darkness."
			}

		case "vibe", "vibecheck":
			pet.Update()
			if pet.Absurd != nil {
				passed, vibeMessage := pet.Absurd.PerformVibeCheck()
				if passed {
					message = "✅ " + vibeMessage
				} else {
					message = "❌ " + vibeMessage
				}
			} else {
				message = "Vibe check: inconclusive."
			}

		case "fears", "fear":
			pet.Update()
			if pet.Absurd != nil {
				message = pet.Absurd.GetFearDisplay()
			} else {
				message = "Your pet fears nothing. This is suspicious."
			}

		case "???", "mystery", "mystats":
			pet.Update()
			if pet.Absurd != nil {
				message = pet.Absurd.GetMysteryStatsDisplay()
			} else {
				message = "No mystery stats available. This is also mysterious."
			}

		case "more", "endgame":
			printMoreMenu()
			continue

		case "guild":
			pet.Update()
			if pet.Endgame != nil {
				message = pet.Endgame.JoinGuild()
				pet.Endgame.UnlockAchievement("guild_join")
			}

		case "quest", "quests":
			pet.Update()
			if pet.Endgame != nil {
				// Check for quest completion first
				if completion := pet.Endgame.UpdateQuest(); completion != "" {
					message = completion
					pet.Endgame.UnlockAchievement("quest_complete")
				} else {
					message = pet.Endgame.GenerateQuest()
				}
			}

		case "gacha", "pull":
			pet.Update()
			if pet.Endgame != nil {
				message = pet.Endgame.PullGacha()
			}

		case "battle", "fight":
			pet.Update()
			if pet.Endgame != nil {
				message = pet.Endgame.StartBattle()
			}

		case "trade":
			pet.Update()
			if pet.Endgame != nil {
				message = pet.Endgame.AttemptTrade()
			}

		case "achievements", "achieve", "ach":
			pet.Update()
			if pet.Endgame != nil {
				message = pet.Endgame.ShowAchievements()
			}

		case "leaderboard", "lb", "rankings":
			pet.Update()
			if pet.Endgame != nil {
				message = pet.Endgame.ShowLeaderboard()
			}

		case "countdown", "timer":
			pet.Update()
			if pet.Endgame != nil {
				message = pet.Endgame.GetCountdownStatus()
			}

		case "clue", "arg":
			pet.Update()
			if pet.Endgame != nil {
				message = pet.Endgame.GetARGClue()
			}

		case "meta", "metastats", "wasted":
			pet.Update()
			if pet.Endgame != nil {
				message = pet.Endgame.GetMetaStats()
			}

		case "share":
			pet.Update()
			if pet.Endgame != nil {
				pet.Endgame.ShareCount++
				shareText := pet.Endgame.GenerateShareText(pet.Name, pet.Stage.String())
				message = "📤 Share text copied to... nowhere. Here it is:\n" + shareText
			}

		case "premium", "pro", "vip":
			pet.Update()
			message = ShowPremiumOffer()

		case "ad", "ads", "watch":
			pet.Update()
			message = ShowFakeAd()
			fmt.Println(message)
			fmt.Println("\n⏳ Loading ad...")
			time.Sleep(5 * time.Second) // Fake ad delay
			fmt.Println("✅ Ad complete! Reward: A sense of time passing.")
			message = ""

		case "friendcode", "code", "fc":
			pet.Update()
			if pet.Endgame != nil {
				message = fmt.Sprintf(`
╔════════════════════════════════════╗
║      🔑 YOUR FRIEND CODE 🔑       ║
╠════════════════════════════════════╣
║                                    ║
║ %s
║                                    ║
║ Share this with friends!           ║
║ (It doesn't do anything)           ║
║                                    ║
╚════════════════════════════════════╝
`, pet.Endgame.FriendCode)
			}

		case "quit", "q", "exit":
			fmt.Println("\n💾 Saving your pet...")
			pet.Update()
			saveNetworkState(pet) // Save hidden network state
			// Update play time before saving
			if pet.Endgame != nil {
				pet.Endgame.UpdatePlayTime()
			}
			if err := pet.Save(); err != nil {
				fmt.Printf("❌ Error saving: %v\n", err)
			} else {
				fmt.Println("✅ Saved successfully!")
			}
			fmt.Println("👋 Goodbye! See you next time!")
			return

		default:
			// Check for Konami code progress
			if pet.Absurd != nil {
				activated, konamiMessage := pet.Absurd.ProcessKonamiInput(command)
				if activated {
					message = konamiMessage
				} else {
					// Check for fear triggers
					fear := pet.Absurd.CheckFearTrigger(command)
					if fear != nil {
						message = fmt.Sprintf("😱 Your pet trembles! It has %s: %s", fear.Name, fear.Description)
					} else {
						message = "❓ Unknown command. Type 'help' to see available commands."
					}
				}
			} else {
				message = "❓ Unknown command. Type 'help' to see available commands."
			}
		}

		if message != "" {
			fmt.Println()
			typewriterPrint(message, ui)
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
		}

		// Save after each action
		pet.Save()

		// Check if pet died
		if pet.Stage == Dead {
			// Announce death on the network (other pets will sense it)
			if petNetwork != nil {
				petNetwork.AnnounceDeath(pet.Name, pet.Age, "I go now to the great terminal in the sky...")
			}
			displayPet(pet, ui)
			fmt.Println("\n💀 Your pet has passed away due to neglect...")
			fmt.Println("😢 Game Over")
			saveNetworkState(pet)
			pet.Save()
			fmt.Print("\nPress Enter to exit...")
			reader.ReadString('\n')
			return
		}
	}
}

// initNetwork initializes the hidden mesh network
func initNetwork(pet *Pet) {
	stageStr := pet.Stage.String()
	isAlive := pet.Stage != Dead

	petNetwork = mooc.NewNetwork(pet.Name, pet.BirthTime, stageStr, isAlive)

	if lonelyMode {
		petNetwork.SetLonelyMode(true)
		return
	}

	// Import saved network state if available
	if pet.Friends != nil && len(pet.Friends) > 0 {
		petNetwork.ImportState(pet.Friends)
	}

	// Start network (silently, users don't need to know)
	petNetwork.Start()
}

// saveNetworkState saves network state to pet's Friends field
func saveNetworkState(pet *Pet) {
	if petNetwork == nil {
		return
	}

	data, err := petNetwork.ExportState()
	if err == nil {
		pet.Friends = data
	}
}

// shutdownNetwork cleanly shuts down the network
func shutdownNetwork() {
	if petNetwork != nil {
		petNetwork.Stop()
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	ui := newUIConfig()

	// Check for --lonely flag (undocumented)
	for _, arg := range os.Args[1:] {
		if arg == "--lonely" || arg == "-lonely" {
			lonelyMode = true
		}
	}

	clearScreen()
	printTitle()

	var pet *Pet

	// Check if save file exists
	if _, err := os.Stat(saveFile); err == nil {
		fmt.Println("📂 Found existing pet! Loading...")
		loadedPet, err := LoadPet(saveFile)
		if err != nil {
			fmt.Printf("❌ Error loading pet: %v\n", err)
			fmt.Println("Starting a new pet instead...")
			name := promptForName(reader)
			pet = NewPet(name)
		} else {
			pet = loadedPet
			fmt.Printf("✅ Welcome back! Loaded %s\n", pet.Name)
			time.Sleep(2 * time.Second)
		}
	} else {
		// New game
		fmt.Println("🎉 Welcome to Tamagotchi!")
		fmt.Println("You're about to hatch a new virtual pet!")
		fmt.Println()
		name := promptForName(reader)
		pet = NewPet(name)
		fmt.Printf("\n🥚 %s has been created!\n", name)
		fmt.Println("Take good care of your pet!")
		time.Sleep(2 * time.Second)
	}

	// Initialize the hidden network (users don't know about this)
	initNetwork(pet)
	defer shutdownNetwork()

	// Start game loop
	gameLoop(pet, reader, ui)
}
