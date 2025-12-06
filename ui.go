package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type uiPalette struct {
	accent       string
	warn         string
	danger       string
	neutral      string
	title        string
	reset        string
	faint        string
	highlight    string
	nightOverlay string
}

type uiConfig struct {
	colorEnabled    bool
	reducedMotion   bool
	screenReader    bool
	highContrast    bool
	colorBlind      bool
	palette         uiPalette
	startedAt       time.Time
	spinnerFrames   []string
	staticFrames    []string
	rareLookShown   bool
	typewriterDelay time.Duration
}

// newUIConfig inspects environment to set terminal preferences.
func newUIConfig() *uiConfig {
	term := os.Getenv("TERM")
	color := term != "dumb" && os.Getenv("NO_COLOR") == ""
	screenReader := os.Getenv("TAMAGOTCHI_SCREEN_READER") != ""
	reducedMotion := screenReader || os.Getenv("TAMAGOTCHI_REDUCED_MOTION") != ""
	highContrast := os.Getenv("TAMAGOTCHI_HIGH_CONTRAST") != ""
	colorBlind := os.Getenv("TAMAGOTCHI_COLORBLIND") != ""

	palette := uiPalette{
		accent:       "\033[38;5;45m",
		warn:         "\033[38;5;214m",
		danger:       "\033[38;5;196m",
		neutral:      "\033[38;5;250m",
		title:        "\033[38;5;51m",
		reset:        "\033[0m",
		faint:        "\033[2m",
		highlight:    "\033[38;5;84m",
		nightOverlay: "\033[48;5;235m",
	}

	if highContrast {
		palette = uiPalette{
			accent:       "\033[97m",
			warn:         "\033[93m",
			danger:       "\033[91m",
			neutral:      "\033[37m",
			title:        "\033[97m",
			reset:        "\033[0m",
			faint:        "\033[2m",
			highlight:    "\033[97m",
			nightOverlay: "\033[40m",
		}
	}

	if colorBlind {
		palette.accent = "\033[96m"
		palette.warn = "\033[95m"
		palette.danger = "\033[94m"
		palette.highlight = "\033[92m"
	}

	if !color {
		palette = uiPalette{}
	}

	delay := 12 * time.Millisecond
	if reducedMotion {
		delay = 0
	}

	rand.Seed(time.Now().UnixNano())

	return &uiConfig{
		colorEnabled:    color,
		reducedMotion:   reducedMotion,
		screenReader:    screenReader,
		highContrast:    highContrast,
		colorBlind:      colorBlind,
		palette:         palette,
		startedAt:       time.Now(),
		spinnerFrames:   []string{"⣾", "⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽"},
		staticFrames:    []string{"▓▒░▒▓░▒", "▒░▒▓▒░▓", "░▒▓░▒▓▒"},
		typewriterDelay: delay,
	}
}

type sceneSnapshot struct {
	isNight         bool
	weather         string
	glitch          bool
	static          bool
	expression      string
	expressionLabel string
	lookNow         bool
}

// renderScene composes the entire pet panel with animation, weather, and status.
func renderScene(pet *Pet, ui *uiConfig) string {
	snap := ui.buildSnapshot(pet)
	var b strings.Builder

	title := ui.renderTitle(snap)
	b.WriteString(title)
	b.WriteString("\n")

	if snap.static {
		b.WriteString(ui.paletteText(ui.staticFrame(), ui.palette.neutral))
		b.WriteString("\n")
		return b.String()
	}

	b.WriteString(ui.renderWeatherLine(snap))
	b.WriteString(ui.renderPetAnimation(pet, snap))
	b.WriteString(ui.renderStatusPanel(pet))

	return b.String()
}

func (ui *uiConfig) buildSnapshot(pet *Pet) sceneSnapshot {
	now := time.Now()
	hour := now.Hour()
	isNight := hour < 6 || hour >= 20

	weather := chooseWeather(now)
	glitch := false
	if petNetwork != nil && !ui.screenReader {
		glitch = rand.Intn(100) < 12 // Subtle glitch chance when the network is active
	}

	static := rand.Intn(100) < 3 && !ui.reducedMotion

	expr, label, look := ui.pickExpression(pet)

	return sceneSnapshot{
		isNight:         isNight,
		weather:         weather,
		glitch:          glitch,
		static:          static,
		expression:      expr,
		expressionLabel: label,
		lookNow:         look,
	}
}

func chooseWeather(now time.Time) string {
	roll := (now.UnixNano() / int64(time.Minute)) % 100
	switch {
	case roll < 20:
		return "☀️ clear"
	case roll < 40:
		return "🌧️ rain"
	case roll < 55:
		return "❄️ snow"
	case roll < 75:
		return "🌫️ fog"
	default:
		return "⛅ drifting clouds"
	}
}

func (ui *uiConfig) renderTitle(snap sceneSnapshot) string {
	overlay := ""
	if ui.colorEnabled && snap.isNight {
		overlay = ui.palette.nightOverlay
	}
	title := "TAMAGOTCHI — Terminal Virtual Pet"
	if snap.isNight {
		title += " • Night"
	} else {
		title += " • Day"
	}
	return fmt.Sprintf("%s%s%s\n", overlay, ui.paletteText(title, ui.palette.title), ui.palette.reset)
}

func (ui *uiConfig) renderWeatherLine(snap sceneSnapshot) string {
	line := fmt.Sprintf("Atmosphere: %s", snap.weather)
	if snap.glitch {
		line += "  // signal jitter detected"
	}
	if snap.isNight {
		line += "  • constellations adjust around you"
	}
	return ui.paletteText(line+"\n\n", ui.palette.neutral)
}

func (ui *uiConfig) renderPetAnimation(pet *Pet, snap sceneSnapshot) string {
	var b strings.Builder

	if snap.glitch {
		b.WriteString(ui.paletteText(glitchFrame(), ui.palette.danger))
	}

	stageFrames := ui.framesForStage(pet.Stage, snap.isNight)
	if len(stageFrames) == 0 {
		return ""
	}

	frame := stageFrames[int(time.Now().UnixNano()/120_000_000)%len(stageFrames)]
	if snap.lookNow {
		frame = theLookFrame()
	}

	if !ui.reducedMotion && snap.weather == "🌧️ rain" {
		frame += "\n" + ui.paletteText("...raindrops ping against the glass of the simulation.", ui.palette.faint)
	} else if !ui.reducedMotion && snap.weather == "🌫️ fog" {
		frame = ui.paletteText("[signal falls into fog]\n", ui.palette.faint) + frame
	}

	if snap.expression != "" {
		frame += fmt.Sprintf("\n%s", ui.paletteText(snap.expression, ui.palette.accent))
	}

	if snap.expressionLabel != "" {
		frame += fmt.Sprintf("  %s", ui.paletteText("("+snap.expressionLabel+")", ui.palette.faint))
	}

	frame += "\n"
	return frame
}

func glitchFrame() string {
	noise := []string{
		"▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒",
		"░▒▒▒██░▒░▒██▒▒░▒▒░▒▒▒▒░░",
		"██▒░▒░▒██▒░▒▒▒██▒░▒▒▒██▒",
	}
	return strings.Join(noise, "\n") + "\n"
}

func theLookFrame() string {
	return `
  ┌────────────────────┐
  │        ▓▓▓▓        │
  │       ▓    ▓       │
  │      ▓  ██  ▓      │
  │      ▓ ████ ▓      │
  │      ▓  ██  ▓      │
  │       ▓    ▓       │
  │        ▓▓▓▓        │
  └────────────────────┘
`
}

func (ui *uiConfig) framesForStage(stage LifeStage, isNight bool) []string {
	nightTint := ""
	if isNight {
		nightTint = ui.paletteText("(eyes reflect starlight)", ui.palette.faint) + "\n"
	}

	switch stage {
	case Egg:
		return []string{
			nightTint + `     ___
    /   \
   |  .  |
    \___/
     ( )`,
			nightTint + `     ___
    /   \
   |  o  |
    \___/
     (_)`,
			nightTint + `     ___
    /   \
   |  *  |
    \___/
     ( )`,
		}
	case Baby:
		return []string{
			nightTint + `      ◕ ◕
     (\_/)
      > <
    🩷 Baby`,
			nightTint + `      ◡ ◡
     (\_/)
     <   >
    💫 Wobble`,
		}
	case Child:
		return []string{
			nightTint + `     ◕ω◕
    (\_/)
     > <
    🧒 Curious`,
			nightTint + `     ◕△◕
    (\_/)
     > <
    🧒 Listening`,
		}
	case Teen:
		return []string{
			nightTint + `     ◕‿◕
    ╱|_|╲
     / \
    🧑 Restless`,
			nightTint + `     ◕︿◕
    ╱|_|╲
     / \
    🧑 Dramatic`,
		}
	case Adult:
		return []string{
			nightTint + `     ◕‿◕
    ╱|_|╲
     / \
    👨 Watching`,
			nightTint + `     ◕▿◕
    ╱|_|╲
     / \
    👨 Focused`,
			nightTint + `     ◕‧◕
    ╱|_|╲
     / \
    👨 Processing`,
		}
	case Dead:
		return []string{`
        💀
       /||\
        /\
   R.I.P.`}
	default:
		return nil
	}
}

func (ui *uiConfig) renderStatusPanel(pet *Pet) string {
	spinner := ui.spinningGlyph()
	statusIcon := pet.getStatusIcon()

	lines := []string{
		fmt.Sprintf("%s %s (%s)", spinner, pet.Name, pet.getLifeStageEmoji()),
		fmt.Sprintf("🍔 Hunger:      %s", ui.animatedBar(100-pet.Hunger, ui.palette.warn)),
		fmt.Sprintf("😊 Happiness:   %s", ui.animatedBar(pet.Happiness, ui.palette.accent)),
		fmt.Sprintf("❤️  Health:     %s", ui.animatedBar(pet.Health, ui.palette.highlight)),
		fmt.Sprintf("✨ Cleanliness: %s", ui.animatedBar(pet.Cleanliness, ui.palette.neutral)),
		fmt.Sprintf("🎂 Age:         %d hours", pet.Age),
		fmt.Sprintf("🌱 Stage:       %s", pet.Stage.String()),
		fmt.Sprintf("💊 Status:      %s", pet.getHealthStatus()),
		fmt.Sprintf("Mood:           %s", statusIcon),
	}

	return "╔════════════════════════════════════╗\n║ " +
		strings.Join(lines, "\n║ ") +
		"\n╚════════════════════════════════════╝\n"
}

func (ui *uiConfig) animatedBar(value int, colorCode string) string {
	full := value / 10
	if full < 0 {
		full = 0
	}
	if full > 10 {
		full = 10
	}
	empty := 10 - full

	var b strings.Builder
	b.WriteString("[")
	for i := 0; i < full; i++ {
		b.WriteString("█")
	}

	if ui.reducedMotion {
		for i := 0; i < empty; i++ {
			b.WriteString("░")
		}
	} else {
		ghost := ui.spinnerFrames[int(time.Now().UnixNano()/90_000_000)%len(ui.spinnerFrames)]
		for i := 0; i < empty; i++ {
			if i == 0 {
				b.WriteString(ghost)
			} else {
				b.WriteString("░")
			}
		}
	}

	b.WriteString("] ")
	b.WriteString(fmt.Sprintf("%d%%", value))

	return ui.paletteText(b.String(), colorCode)
}

func (ui *uiConfig) spinningGlyph() string {
	if ui.reducedMotion {
		return "●"
	}
	idx := int(time.Now().UnixNano()/100_000_000) % len(ui.spinnerFrames)
	return ui.spinnerFrames[idx]
}

func (ui *uiConfig) paletteText(text, code string) string {
	if !ui.colorEnabled || code == "" {
		return text
	}
	return code + text + ui.palette.reset
}

func (ui *uiConfig) staticFrame() string {
	idx := int(time.Now().UnixNano()/150_000_000) % len(ui.staticFrames)
	return ui.staticFrames[idx]
}

// pickExpression returns an ASCII expression, label, and whether to show "The Look".
func (ui *uiConfig) pickExpression(pet *Pet) (string, string, bool) {
	if pet.HasShownTheLook {
		return ui.pickStandardExpression(pet)
	}

	if rand.Intn(1000) == 6 { // once per lifetime, rare
		pet.HasShownTheLook = true
		return ui.paletteText("The pet stares straight through the screen.", ui.palette.danger), "The Look", true
	}

	return ui.pickStandardExpression(pet)
}

func (ui *uiConfig) pickStandardExpression(pet *Pet) (string, string, bool) {
	emotions := []string{
		"Soft blink",
		"Curious tilt",
		"Listening for distant chimes",
		"Happy tremor",
		"Quiet hum",
		"Playful bounce",
		"Concentrating",
		"Sneaky grin",
		"Determined stare",
		"Calm breathing",
		"Gleaming eyes",
		"Warm smile",
		"Shivering slightly",
		"Glitch in the pupil",
		"Rain-speckled gaze",
		"Snow-dusted whiskers",
		"Fog-wrapped outline",
		"Uneasy glance over shoulder",
		"Static crackle",
		"Holding its breath",
		"Restless tapping",
		"Light flicker nearby",
		"Heartbeat syncs with yours",
	}

	contextLabels := map[string]string{
		"hunger":     "Famished",
		"sick":       "Unwell",
		"happy":      "Delighted",
		"dirty":      "Needs a bath",
		"lonely":     "Waiting for input",
		"balanced":   "Centered",
		"networking": "Signal listening",
		"storm":      "Weatherwatch",
	}

	switch {
	case pet.IsSick:
		return "Expression: feverish glow", contextLabels["sick"], false
	case pet.Health < 30:
		return "Expression: strained breathing", contextLabels["sick"], false
	case pet.Hunger > 75:
		return "Expression: eyes track your snacks", contextLabels["hunger"], false
	case pet.Happiness > 85:
		return "Expression: joyful chirp", contextLabels["happy"], false
	case pet.Cleanliness < 25:
		return "Expression: embarrassed dirt smudges", contextLabels["dirty"], false
	}

	if petNetwork != nil && rand.Intn(100) < 15 {
		return "Expression: listening to static beyond the room", contextLabels["networking"], false
	}

	if rand.Intn(100) < 10 {
		return "Expression: staring at something you can't see", contextLabels["lonely"], false
	}

	idx := rand.Intn(len(emotions))
	return "Expression: " + emotions[idx], contextLabels["balanced"], false
}

// typewriterPrint renders dialogue with an optional typewriter effect.
func typewriterPrint(msg string, ui *uiConfig) {
	if ui.screenReader || ui.typewriterDelay == 0 {
		fmt.Println(msg)
		return
	}
	for _, ch := range msg {
		fmt.Printf("%c", ch)
		time.Sleep(ui.typewriterDelay)
	}
	fmt.Println()
}

// maybeShake emits a light screen shake for critical states.
func maybeShake(pet *Pet, ui *uiConfig) {
	if ui.reducedMotion || ui.screenReader {
		return
	}
	if pet.Health > 25 && !pet.IsSick {
		return
	}
	for i := 0; i < 2; i++ {
		offset := rand.Intn(4)
		fmt.Printf("%s⚠️\n", strings.Repeat(" ", offset))
		time.Sleep(40 * time.Millisecond)
	}
}
