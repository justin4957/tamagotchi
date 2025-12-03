# Tamagotchi Demo

## Quick Start

```bash
# Build the game
go build -o tamagotchi

# Run it
./tamagotchi
```

## Sample Gameplay

When you first run the game, you'll see:

```
╔═══════════════════════════════════════════════╗
║                                               ║
║   🎮 TAMAGOTCHI - Virtual Pet Simulator 🎮   ║
║              Relive the 90s Magic!            ║
║                                               ║
╚═══════════════════════════════════════════════╝

🎉 Welcome to Tamagotchi!
You're about to hatch a new virtual pet!

What would you like to name your new pet? 
```

After naming your pet, you'll see the main game screen with:

- ASCII art of your pet in its current life stage
- Stats displayed as progress bars (hunger, happiness, health, cleanliness)
- Current age and life stage
- Available commands menu

## Example Session

```
Name: Mochi

     ___
    /   \
   |  ?  |
    \___/
    🥚 Egg

╔════════════════════════════════════╗
║      🥚 Mochi (🥚)
╠════════════════════════════════════╣
║ 🍔 Hunger:      [██████████] 0%
║ 😊 Happiness:   [██████████] 100%
║ ❤️  Health:     [██████████] 100%
║ ✨ Cleanliness: [██████████] 100%
║ 🎂 Age:         0 hours
║ 🌱 Stage:       Egg
║ 💊 Status:      Excellent
╚════════════════════════════════════╝

Commands:
  feed   - Feed your pet 🍔
  play   - Play with your pet 🎮
  clean  - Clean up after your pet 🛁
  heal   - Give medicine to your pet 💊
  status - Check your pet's status 📊
  help   - Show this menu 📖
  quit   - Save and exit 👋

Enter command: 
```

## Life Stages Timeline

After waiting or simulating time:

**1 hour later - Baby:**
```
      ◕ ◕
     (\_/)
      > <
    👶 Baby
```

**24 hours later - Child:**
```
     ◕ω◕
    (\_/)
     > <
    🧒 Child
```

**48 hours later - Teen:**
```
     ◕‿◕
    ╱|_|╲
     / \
    🧑 Teen
```

**72 hours later - Adult:**
```
     ◕‿◕
    ╱|_|╲
     / \
    👨 Adult
```

## Interaction Examples

**Feeding:**
```
Enter command: feed
😋 Yum! That was delicious!
```

**Playing:**
```
Enter command: play
🎮 Wheee! That was so much fun!
```

**Cleaning:**
```
Enter command: clean
🛁 Ahh, much better!
```

**If your pet gets sick:**
```
Enter command: heal
💊 Thank you! I feel much better now!
```

## Tips for Keeping Your Pet Alive

1. Check in every few hours
2. Feed when hunger > 30%
3. Play when happiness < 70%
4. Clean when cleanliness < 50%
5. Use medicine immediately if sick
6. Keep health above 50% at all times

## Game Mechanics

### Stat Degradation Rates (per hour)
- **Egg**: No degradation
- **Baby**: 0.5x rate
- **Child**: 1.0x rate (base)
- **Teen**: 1.5x rate
- **Adult**: 2.0x rate

### Base Degradation
- Hunger: +5 per hour
- Happiness: -3 per hour
- Cleanliness: -4 per hour
- Health: Depends on other stats

### Health System
- Health decreases if hunger > 70 OR happiness < 30 OR cleanliness < 30
- Health recovers if all stats are good (hunger < 30, happiness > 70, cleanliness > 70)
- Pet gets sick if health < 50 OR cleanliness < 20
- Pet dies if health reaches 0

## Persistence

The game automatically saves:
- Every 30 seconds
- After each command
- When you quit

Your pet's state is saved to `tamagotchi_save.json` and will persist across sessions.

## Testing

Run the test suite:
```bash
go test -v
```

All 11 tests should pass, covering:
- Pet creation
- All actions (feed, play, clean, heal)
- Life stage progression
- Stat degradation
- Sickness and death mechanics
- Egg behavior
