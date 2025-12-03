package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// MiniGameResult represents the outcome of a mini-game
type MiniGameResult struct {
	Message string
	Success bool
}

// PlayWatchPaintDry plays the "Watch Paint Dry" mini-game
// Literally just a timer with no reward
func PlayWatchPaintDry(reader *bufio.Reader) MiniGameResult {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║    🎨 WATCH PAINT DRY 🎨          ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║ Watch the paint dry for 10 seconds ║")
	fmt.Println("║ Press Enter to start...            ║")
	fmt.Println("╚════════════════════════════════════╝")

	reader.ReadString('\n')

	paintStages := []string{
		"The paint is wet. Very wet.",
		"The paint is still wet.",
		"Is it drying? Hard to tell.",
		"The paint glistens ominously.",
		"You think you see it drying.",
		"No, still wet.",
		"The paint mocks your patience.",
		"Drying... maybe...",
		"Almost there? Probably not.",
		"The paint is dry. Or is it?",
	}

	for i, stage := range paintStages {
		fmt.Printf("\r[%d/10] %s", i+1, stage)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n\n✅ Congratulations! You watched paint dry.")
	fmt.Println("🏆 Reward: None. What did you expect?")

	return MiniGameResult{
		Message: "You watched paint dry. Time you'll never get back.",
		Success: true, // Success is meaningless here
	}
}

// PlayStareContest plays the "Stare Contest" mini-game
// Press any key and you lose, don't press and nothing happens
func PlayStareContest(reader *bufio.Reader) MiniGameResult {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║    👁️ STARE CONTEST 👁️            ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║ Rules:                              ║")
	fmt.Println("║ - Don't press any key               ║")
	fmt.Println("║ - If you press a key, you lose      ║")
	fmt.Println("║ - If you don't press, nothing happens║")
	fmt.Println("║                                      ║")
	fmt.Println("║ The contest has already begun...    ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Println("\n       👁️     👁️")
	fmt.Println("         ___")
	fmt.Println("        \\   /")
	fmt.Println("         ---")
	fmt.Println("\n   Your pet stares at you.")
	fmt.Println("   You stare at your pet.")
	fmt.Println("   The universe holds its breath.")
	fmt.Println("\n   (Press any key to blink and lose)")

	reader.ReadString('\n')

	fmt.Println("\n❌ YOU BLINKED!")
	fmt.Println("Your pet wins. Your pet always wins.")
	fmt.Println("The staring contest was rigged from the start.")

	return MiniGameResult{
		Message: "You lost the stare contest. Inevitable.",
		Success: false,
	}
}

// PlayCountToThousand plays the "Count to 1000" mini-game
// Manual counting, loses progress if you mistype
func PlayCountToThousand(reader *bufio.Reader) MiniGameResult {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║    🔢 COUNT TO 1000 🔢             ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║ Rules:                              ║")
	fmt.Println("║ - Type numbers from 1 to 1000       ║")
	fmt.Println("║ - One wrong number resets everything║")
	fmt.Println("║ - Type 'quit' to give up            ║")
	fmt.Println("║                                      ║")
	fmt.Println("║ Good luck. You'll need it.          ║")
	fmt.Println("╚════════════════════════════════════╝")

	currentNumber := 1
	highestReached := 0

	for currentNumber <= 1000 {
		fmt.Printf("\nEnter %d: ", currentNumber)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "quit" {
			fmt.Printf("\n😔 You gave up at %d.", currentNumber)
			if highestReached > 0 {
				fmt.Printf(" Highest reached: %d", highestReached)
			}
			fmt.Println("\n🏆 Reward: The wisdom that some things aren't worth doing.")
			return MiniGameResult{
				Message: fmt.Sprintf("Gave up counting at %d. Wisdom gained.", currentNumber),
				Success: false,
			}
		}

		num, err := strconv.Atoi(input)
		if err != nil || num != currentNumber {
			fmt.Println("\n❌ WRONG!")
			fmt.Printf("You typed '%s' but needed '%d'\n", input, currentNumber)
			if currentNumber > highestReached {
				highestReached = currentNumber
			}
			fmt.Printf("Progress reset. Highest reached this session: %d\n", highestReached)
			fmt.Println("Starting over from 1...")
			currentNumber = 1
			continue
		}

		// Progress indicators
		switch currentNumber {
		case 10:
			fmt.Println("   ...only 990 to go.")
		case 50:
			fmt.Println("   ...you're really doing this, huh?")
		case 100:
			fmt.Println("   ...10% done. Are you okay?")
		case 250:
			fmt.Println("   ...25%. There's still time to quit.")
		case 500:
			fmt.Println("   ...halfway. No turning back now.")
		case 750:
			fmt.Println("   ...75%. The end is in sight.")
		case 900:
			fmt.Println("   ...so close. Don't mess up.")
		case 999:
			fmt.Println("   ...one more. Don't choke.")
		}

		currentNumber++
	}

	// If someone actually reaches 1000
	fmt.Println("\n🎉🎉🎉 YOU ACTUALLY DID IT 🎉🎉🎉")
	fmt.Println("You counted to 1000. Manually. One number at a time.")
	fmt.Println("🏆 Reward: A profound sense of... something. Not accomplishment.")
	fmt.Println("Maybe regret? It's hard to say.")
	fmt.Println("\nYour pet looks at you with what might be respect.")
	fmt.Println("Or concern. Probably concern.")

	return MiniGameResult{
		Message: "Counted to 1000. Why? Nobody knows.",
		Success: true,
	}
}

// PlayDoNothing plays the "Do Nothing" mini-game
// The game of doing absolutely nothing
func PlayDoNothing(reader *bufio.Reader) MiniGameResult {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║    🧘 DO NOTHING 🧘                ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║ Instructions:                       ║")
	fmt.Println("║ - Do nothing                        ║")
	fmt.Println("║ - Press Enter when done doing nothing║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Println("\n   Doing nothing...")
	fmt.Println("   ...")
	fmt.Println("   ...")
	fmt.Println("   (You're doing great at nothing)")
	fmt.Println("\n   Press Enter to stop doing nothing")

	reader.ReadString('\n')

	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	nothingTime := 1 + randomSource.Intn(60)

	fmt.Printf("\n✅ You did nothing for approximately %d seconds.\n", nothingTime)
	fmt.Println("🏆 Achievement Unlocked: Nothing")

	return MiniGameResult{
		Message: fmt.Sprintf("Did nothing for %d seconds. Impressive.", nothingTime),
		Success: true,
	}
}

// PlayGuessTheNumber plays a guess the number game where the number changes
func PlayGuessTheNumber(reader *bufio.Reader) MiniGameResult {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║    🎲 GUESS THE NUMBER 🎲          ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║ I'm thinking of a number 1-10      ║")
	fmt.Println("║ You have 3 guesses                 ║")
	fmt.Println("║ Type 'quit' to give up             ║")
	fmt.Println("╚════════════════════════════════════╝")

	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	for guess := 1; guess <= 3; guess++ {
		// The number changes each guess because the game is unfair
		targetNumber := 1 + randomSource.Intn(10)

		fmt.Printf("\nGuess %d/3: ", guess)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "quit" {
			fmt.Println("\n😔 You gave up.")
			fmt.Printf("The number was %d. Or was it? It kept changing.\n", targetNumber)
			return MiniGameResult{
				Message: "Gave up guessing. The game was rigged anyway.",
				Success: false,
			}
		}

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > 10 {
			fmt.Println("That's not a valid number between 1 and 10.")
			guess-- // Don't count this guess
			continue
		}

		if num == targetNumber {
			// This should rarely happen but it's possible
			fmt.Println("\n🎉 IMPOSSIBLE! You got it!")
			fmt.Println("The number was changing each guess, but you got lucky.")
			fmt.Println("🏆 Reward: Existential uncertainty about probability")
			return MiniGameResult{
				Message: "Won an unwinnable game. Reality questioned.",
				Success: true,
			}
		}

		if guess < 3 {
			if num < targetNumber {
				fmt.Println("Too low! (The number has also changed now)")
			} else {
				fmt.Println("Too high! (But the number shifted)")
			}
		}
	}

	fmt.Println("\n❌ Out of guesses!")
	fmt.Println("The number was... well, it kept changing.")
	fmt.Println("This game was never fair.")
	fmt.Println("🏆 Reward: Understanding that some games can't be won")

	return MiniGameResult{
		Message: "Lost guess the number. The game was rigged.",
		Success: false,
	}
}

// ShowMiniGameMenu displays available mini-games
func ShowMiniGameMenu() {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║     🎮 USELESS MINI-GAMES 🎮       ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║ 1. Watch Paint Dry                 ║")
	fmt.Println("║ 2. Stare Contest                   ║")
	fmt.Println("║ 3. Count to 1000                   ║")
	fmt.Println("║ 4. Do Nothing                      ║")
	fmt.Println("║ 5. Guess the Number                ║")
	fmt.Println("║                                    ║")
	fmt.Println("║ Type 'back' to return              ║")
	fmt.Println("╚════════════════════════════════════╝")
}

// SelectAndPlayMiniGame handles mini-game selection and playing
func SelectAndPlayMiniGame(reader *bufio.Reader) *MiniGameResult {
	ShowMiniGameMenu()

	for {
		fmt.Print("\nSelect a game (1-5): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "1", "paint", "watch":
			result := PlayWatchPaintDry(reader)
			return &result
		case "2", "stare", "contest":
			result := PlayStareContest(reader)
			return &result
		case "3", "count", "1000":
			result := PlayCountToThousand(reader)
			return &result
		case "4", "nothing", "do nothing":
			result := PlayDoNothing(reader)
			return &result
		case "5", "guess", "number":
			result := PlayGuessTheNumber(reader)
			return &result
		case "back", "quit", "exit":
			return nil
		default:
			fmt.Println("Unknown game. Try a number 1-5 or 'back'.")
		}
	}
}
