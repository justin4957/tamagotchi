# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Build the application
go build -o tamagotchi

# Run the application
./tamagotchi

# Run all tests
go test -v

# Run a specific test
go test -v -run TestNewPet

# Format code
go fmt ./...
```

## Architecture

This is a terminal-based Tamagotchi virtual pet simulator written in pure Go with no external dependencies.

### Core Files

- **main.go**: Terminal UI, game loop, user input handling, ASCII art animations, and auto-save ticker (30-second interval)
- **pet.go**: Pet state management, stats system, life stage progression, save/load functionality using JSON

### Key Concepts

**Life Stages**: `Egg → Baby → Child → Teen → Adult → Dead` (defined as `LifeStage` iota in pet.go). Progression is time-based using real hours since birth.

**Stats System**: Four stats (Hunger, Happiness, Health, Cleanliness) range 0-100. Degradation rates increase with life stage (e.g., Adult degrades 2x faster than Child). The `Update()` method calculates time-based stat changes.

**Persistence**: JSON save file (`tamagotchi_save.json`) stores pet state. On load, `Update()` is called to apply time-passed degradation since last save.

**Game Loop Pattern**: Each action calls `pet.Update()` first to apply time-based changes, then performs the action, then saves. Auto-save runs in a goroutine.
