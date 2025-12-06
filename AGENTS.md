# Repository Guidelines

## Project Structure & Modules
- `main.go` wires the CLI loop and initializes the pet lifecycle.
- `pet.go` holds core state, stat decay, and serialization (save file `tamagotchi_save.json`).
- `minigames.go`, `absurd.go`, and `endgame.go` provide optional side modes and late-game content.
- `mooc/` implements the mesh networking/identity protocol used by experimental features.
- Tests live alongside sources as `*_test.go`; assets are generated at runtime rather than stored in the repo.

## Build, Test, and Development Commands
- `go run .` — start the simulator from source (uses Go 1.25.3 module settings).
- `go build -o tamagotchi` — produce the release binary in the repo root.
- `go test ./...` — run all unit and integration tests across modules.
- `go test ./... -run TestName` — focus on a single scenario while iterating.

## Coding Style & Naming Conventions
- Go defaults: tabs for indentation, `gofmt` required before sending changes.
- Keep exports minimal; prefer package-private helpers unless consumed by other packages (notably `mooc`).
- Use clear, imperative function names for actions (`feed`, `clean`, `play`) and noun-based structs (`Pet`, `Identity`, `Network`).
- Avoid global mutable state beyond the existing save-path constants; pass dependencies explicitly.

## Testing Guidelines
- Framework: standard library `testing` only; table-driven tests encouraged.
- Name tests with the unit and behavior (`TestPetStatDecay`, `TestNetworkBroadcast`).
- When adding features, cover both happy-path gameplay and failure/timeout cases in `mooc` networking.
- Run `go test ./...` before opening a PR; add regression tests for any bugfix.

## Commit & Pull Request Guidelines
- Commits should be concise, imperative subject lines (`Add mesh gossip`, `Fix save timestamp`), mirroring current history.
- Keep changes focused; include brief context in the body when behavior shifts or data formats change.
- Pull requests should describe the user-facing impact, mention touched commands or flags, and link any tracking issue.
- If UI/UX behavior changes (prompts, timing, save format), include a short repro note or terminal screenshot.

## Security & Configuration Tips
- Saved state is JSON in the repo root; avoid checking in personal playthroughs. Delete `tamagotchi_save.json` before publishing.
- The experimental mesh features open local listeners; prefer running offline during development unless explicitly testing gossip.
- UI modes: set `TAMAGOTCHI_REDUCED_MOTION=1` or `TAMAGOTCHI_SCREEN_READER=1` for low- or no-animation output; `TAMAGOTCHI_HIGH_CONTRAST=1`/`TAMAGOTCHI_COLORBLIND=1` for safer palettes.
