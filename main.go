package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	saveFile = "tamagotchi_save.json"
)

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
  help   - Show this menu 📖
  quit   - Save and exit 👋
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
}

// displayPet shows the pet and its current status
func displayPet(pet *Pet) {
	clearScreen()
	printTitle()
	showPetAnimation(pet)
	fmt.Println(pet.GetStatus())
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
func gameLoop(pet *Pet, reader *bufio.Reader) {
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

	for {
		displayPet(pet)
		printMenu()

		fmt.Print("Enter command: ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(strings.ToLower(command))

		var message string

		switch command {
		case "feed", "f":
			pet.Update()
			message = pet.Feed()

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

		case "quit", "q", "exit":
			fmt.Println("\n💾 Saving your pet...")
			pet.Update()
			if err := pet.Save(); err != nil {
				fmt.Printf("❌ Error saving: %v\n", err)
			} else {
				fmt.Println("✅ Saved successfully!")
			}
			fmt.Println("👋 Goodbye! See you next time!")
			return

		default:
			message = "❓ Unknown command. Type 'help' to see available commands."
		}

		if message != "" {
			fmt.Printf("\n%s\n", message)
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
		}

		// Save after each action
		pet.Save()

		// Check if pet died
		if pet.Stage == Dead {
			displayPet(pet)
			fmt.Println("\n💀 Your pet has passed away due to neglect...")
			fmt.Println("😢 Game Over")
			fmt.Print("\nPress Enter to exit...")
			reader.ReadString('\n')
			return
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

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

	// Start game loop
	gameLoop(pet, reader)
}
