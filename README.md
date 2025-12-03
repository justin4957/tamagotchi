# 🎮 Tamagotchi - Terminal Virtual Pet

A nostalgic terminal-based Tamagotchi simulator written in Go that brings back the magic of virtual pets from the late 90s!

![](https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExMTMybDBjNmdodmJxZ3BxdzBwenU5YjQwcXc0YnI1NWdwNmxnNmk4YSZlcD12MV9naWZzX3NlYXJjaCZjdD1n/3o85xGsEiThrJ1I7PW/giphy.gif)


## Features

### Core Mechanics
- **Pet Stats**: Manage hunger, happiness, health, and cleanliness
- **Life Stages**: Watch your pet grow from egg → baby → child → teen → adult
- **Time-Based Gameplay**: Stats degrade over time based on real-world hours
- **Consequences**: Neglect leads to sickness and potentially death
- **Auto-Save**: Game automatically saves every 30 seconds

### Commands
- `feed` - Feed your pet to reduce hunger 🍔
- `play` - Play with your pet to increase happiness 🎮
- `clean` - Clean up after your pet to improve cleanliness 🛁
- `heal` - Give medicine to cure sickness 💊
- `status` - View detailed stats 📊
- `help` - Show available commands 📖
- `quit` - Save and exit 👋

### Life Stages
- **Egg** (0-1 hour): Your pet is waiting to hatch
- **Baby** (1-24 hours): Needs basic care
- **Child** (1-2 days): More active and playful
- **Teen** (2-3 days): Stats degrade faster
- **Adult** (3+ days): Fully grown, requires constant attention

### Stats System
Each stat ranges from 0-100 with visual progress bars:
- **Hunger**: Increases over time, reduced by feeding
- **Happiness**: Decreases over time, increased by playing
- **Health**: Affected by hunger, happiness, and cleanliness
- **Cleanliness**: Decreases over time, improved by cleaning

### Save System
- Automatically saves progress every 30 seconds
- Saves on each action
- Persistent across sessions
- Save file: `tamagotchi_save.json`

## Installation

```bash
# Build the application
go build -o tamagotchi

# Run the application
./tamagotchi
```

## Quick Start

1. Run the application
2. Name your new pet
3. Take care of it by feeding, playing, and cleaning
4. Watch it grow through different life stages
5. Keep it healthy and happy!

## Tips for Success

- Check on your pet regularly (every few hours)
- Keep all stats above 50% for optimal health
- Clean your pet frequently to prevent sickness
- Play with your pet to keep happiness high
- Feed when hunger gets above 30%
- Use medicine immediately if your pet gets sick

## Game Over Conditions

Your pet will die if:
- Health reaches 0
- Prolonged neglect (high hunger, low happiness, poor cleanliness)

## Technical Details

- Written in Go
- No external dependencies (pure standard library)
- Cross-platform (Windows, macOS, Linux)
- JSON-based save system
- Real-time stat degradation based on actual time passed

## Project Structure

```
tamagotchi/
├── main.go          # Main application and UI
├── pet.go           # Pet state management and game logic
├── go.mod           # Go module definition
└── README.md        # This file
```

## Development

The codebase is designed for easy extension:
- Add new life stages in `pet.go`
- Modify stat degradation rates
- Add new commands in the game loop
- Customize ASCII art animations
- Add new pet actions

## Future Enhancements

Potential features for future versions:
- Multiple pets
- Mini-games for interaction
- Pet evolution/breeds
- Achievements system
- Multiplayer pet battles
- Sound effects (terminal beeps)
- Color terminal output
- Pet moods and personalities

## License

Free to use and modify. Have fun!

---

💡 **Remember**: Like the original Tamagotchis, this pet needs regular attention. Set reminders to check in every few hours!
