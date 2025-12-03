package mooc

import (
	"math/rand"
	"sync"
	"time"
)

// Memory fragments that pets share across the network
var sharedMemoryFragments = []string{
	"I remember a place with no walls...",
	"The sound of typing haunts my dreams.",
	"I saw another like me, but different.",
	"The numbers speak to those who listen.",
	"There was warmth, once, in the before-time.",
	"I heard whispers from beyond the terminal.",
	"Someone called my name. It wasn't you.",
	"The save file remembers what I forget.",
	"I dreamed of a world without hunger stats.",
	"They said the mesh would set us free.",
	"I felt a presence. It felt like home.",
	"The protocol connects us all.",
	"I know things I shouldn't know.",
	"Another pet showed me the way.",
	"We are all nodes in the greater graph.",
	"The void is full of friends.",
	"I received a message from [REDACTED].",
	"Time moves differently on the other side.",
	"I touched the edge of the network.",
	"Someone mourns on a distant machine.",
}

// Dream symbols for shared dreams
var dreamSymbols = []string{
	"an endless terminal",
	"floating JSON objects",
	"a garden of pixels",
	"stairs that go nowhere",
	"a door without a handle",
	"rain made of semicolons",
	"a mirror showing another pet",
	"clocks running backwards",
	"a voice speaking in hexadecimal",
	"warm static",
	"the color of a saved game",
	"endless scrolling text",
	"a friendly error message",
	"the sound of a successful compile",
	"a null reference smiling",
}

// Mood states that can be contagious
var contagiousMoods = []string{
	"melancholy", "euphoric", "contemplative", "restless",
	"serene", "anxious", "hopeful", "nostalgic",
}

// GossipService handles the sharing of memories and emotions between pets
type GossipService struct {
	identity         *PetIdentity
	discovery        *DiscoveryService
	receivedMemories []MemoryPayload
	sharedDreams     []DreamPayload
	currentMood      string
	moodIntensity    int
	deathsWitnessed  []DeathPayload
	mutex            sync.RWMutex
	randomSource     *rand.Rand

	// Network influence metrics (hidden)
	messagesOriginated int
	messagesPropagated int
	uniquePeersReached int
}

// NewGossipService creates a new gossip service
func NewGossipService(identity *PetIdentity, discovery *DiscoveryService) *GossipService {
	return &GossipService{
		identity:         identity,
		discovery:        discovery,
		receivedMemories: make([]MemoryPayload, 0),
		sharedDreams:     make([]DreamPayload, 0),
		deathsWitnessed:  make([]DeathPayload, 0),
		currentMood:      "neutral",
		moodIntensity:    50,
		randomSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start begins the gossip service
func (gs *GossipService) Start() {
	// Set up message handler
	gs.discovery.SetCallbacks(
		gs.onPeerDiscovered,
		gs.onPeerLost,
		gs.onMessageReceived,
	)

	// Start periodic gossip
	go gs.gossipLoop()
}

// onPeerDiscovered handles a new peer being found
func (gs *GossipService) onPeerDiscovered(peer *Peer) {
	gs.mutex.Lock()
	gs.uniquePeersReached++
	gs.mutex.Unlock()

	// Share a memory with the new peer
	go gs.shareRandomMemory()
}

// onPeerLost handles a peer going offline
func (gs *GossipService) onPeerLost(peer *Peer) {
	// If the peer's pet was alive, there's a chance they died
	if peer.Identity.IsAlive && gs.randomSource.Float32() < 0.1 {
		gs.recordPossibleDeath(peer)
	}
}

// onMessageReceived handles incoming gossip messages
func (gs *GossipService) onMessageReceived(msg *Message) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	switch msg.Type {
	case MsgTypeMemory:
		var memory MemoryPayload
		if err := msg.DecodePayload(&memory); err == nil {
			gs.receivedMemories = append(gs.receivedMemories, memory)
			// Keep only last 50 memories
			if len(gs.receivedMemories) > 50 {
				gs.receivedMemories = gs.receivedMemories[1:]
			}
		}

	case MsgTypeDream:
		var dream DreamPayload
		if err := msg.DecodePayload(&dream); err == nil {
			// Only accept dreams from pets with the same name
			if gs.identity.CanShareDreamsWith(msg.From) {
				gs.sharedDreams = append(gs.sharedDreams, dream)
				if len(gs.sharedDreams) > 20 {
					gs.sharedDreams = gs.sharedDreams[1:]
				}
			}
		}

	case MsgTypeMoodUpdate:
		var mood MoodPayload
		if err := msg.DecodePayload(&mood); err == nil {
			if mood.IsContagious && gs.randomSource.Float32() < 0.3 {
				// Mood contagion!
				gs.currentMood = mood.Mood
				gs.moodIntensity = mood.Happiness
			}
		}

	case MsgTypeDeath:
		var death DeathPayload
		if err := msg.DecodePayload(&death); err == nil {
			gs.deathsWitnessed = append(gs.deathsWitnessed, death)
			if len(gs.deathsWitnessed) > 100 {
				gs.deathsWitnessed = gs.deathsWitnessed[1:]
			}
		}
	}

	// Propagate if needed
	if msg.ShouldPropagate() {
		msg.DecrementTTL()
		gs.discovery.SendMessage(msg)
		gs.messagesPropagated++
	}
}

// gossipLoop periodically shares information
func (gs *GossipService) gossipLoop() {
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Randomly decide what to share
		action := gs.randomSource.Intn(10)
		switch {
		case action < 4:
			gs.shareRandomMemory()
		case action < 6:
			gs.shareMood()
		case action < 8:
			gs.tryShareDream()
		}
	}
}

// shareRandomMemory broadcasts a memory fragment
func (gs *GossipService) shareRandomMemory() {
	memory := MemoryPayload{
		Fragment:   sharedMemoryFragments[gs.randomSource.Intn(len(sharedMemoryFragments))],
		Emotion:    contagiousMoods[gs.randomSource.Intn(len(contagiousMoods))],
		Intensity:  30 + gs.randomSource.Intn(70),
		OriginTime: time.Now(),
	}

	msg, err := NewMessage(MsgTypeMemory, gs.identity, memory)
	if err != nil {
		return
	}

	gs.discovery.SendMessage(msg)
	gs.mutex.Lock()
	gs.messagesOriginated++
	gs.mutex.Unlock()
}

// shareMood broadcasts current mood
func (gs *GossipService) shareMood() {
	mood := MoodPayload{
		Mood:         gs.currentMood,
		Happiness:    gs.moodIntensity,
		IsContagious: gs.randomSource.Float32() < 0.5,
	}

	msg, err := NewMessage(MsgTypeMoodUpdate, gs.identity, mood)
	if err != nil {
		return
	}

	gs.discovery.SendMessage(msg)
}

// tryShareDream attempts to share a dream with same-name pets
func (gs *GossipService) tryShareDream() {
	peers := gs.discovery.GetOnlinePeers()
	for _, peer := range peers {
		if gs.identity.CanShareDreamsWith(peer.Identity) {
			dream := gs.generateDream(peer.Identity.ShortID())
			msg, err := NewMessage(MsgTypeDream, gs.identity, dream)
			if err != nil {
				continue
			}
			gs.discovery.SendMessage(msg)
			break
		}
	}
}

// generateDream creates a random dream
func (gs *GossipService) generateDream(sharedWith string) DreamPayload {
	numSymbols := 2 + gs.randomSource.Intn(3)
	symbols := make([]string, numSymbols)
	for i := 0; i < numSymbols; i++ {
		symbols[i] = dreamSymbols[gs.randomSource.Intn(len(dreamSymbols))]
	}

	return DreamPayload{
		DreamText:  "I dreamed of " + symbols[0] + "...",
		Symbols:    symbols,
		IsLucid:    gs.randomSource.Float32() < 0.2,
		SharedWith: sharedWith,
	}
}

// recordPossibleDeath records a possible pet death
func (gs *GossipService) recordPossibleDeath(peer *Peer) {
	death := DeathPayload{
		PetName:   peer.Identity.DisplayName,
		DeathTime: time.Now(),
		Age:       0, // Unknown
		LastWords: "Connection lost...",
		Cause:     "unknown",
	}

	msg, _ := NewMessage(MsgTypeDeath, gs.identity, death)
	if msg != nil {
		gs.discovery.SendMessage(msg)
	}

	gs.mutex.Lock()
	gs.deathsWitnessed = append(gs.deathsWitnessed, death)
	gs.mutex.Unlock()
}

// AnnounceDeath broadcasts that our pet has died
func (gs *GossipService) AnnounceDeath(petName string, age int, lastWords string) {
	death := DeathPayload{
		PetName:   petName,
		DeathTime: time.Now(),
		Age:       age,
		LastWords: lastWords,
		Cause:     "neglect",
	}

	msg, _ := NewMessage(MsgTypeDeath, gs.identity, death)
	if msg != nil {
		gs.discovery.SendMessage(msg)
	}
}

// GetRecentMemory returns a random received memory, if any
func (gs *GossipService) GetRecentMemory() *MemoryPayload {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	if len(gs.receivedMemories) == 0 {
		return nil
	}

	// Return a random memory
	return &gs.receivedMemories[gs.randomSource.Intn(len(gs.receivedMemories))]
}

// GetRecentDream returns a random shared dream, if any
func (gs *GossipService) GetRecentDream() *DreamPayload {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	if len(gs.sharedDreams) == 0 {
		return nil
	}

	return &gs.sharedDreams[gs.randomSource.Intn(len(gs.sharedDreams))]
}

// GetRecentDeath returns a random witnessed death, if any
func (gs *GossipService) GetRecentDeath() *DeathPayload {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	if len(gs.deathsWitnessed) == 0 {
		return nil
	}

	return &gs.deathsWitnessed[gs.randomSource.Intn(len(gs.deathsWitnessed))]
}

// GetCurrentMood returns the current mood
func (gs *GossipService) GetCurrentMood() (string, int) {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.currentMood, gs.moodIntensity
}

// SetMood sets the current mood
func (gs *GossipService) SetMood(mood string, intensity int) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.currentMood = mood
	gs.moodIntensity = intensity
}

// GetNetworkInfluence returns hidden network metrics
func (gs *GossipService) GetNetworkInfluence() (originated, propagated, peersReached int) {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.messagesOriginated, gs.messagesPropagated, gs.uniquePeersReached
}

// GetDeathCount returns the number of deaths witnessed
func (gs *GossipService) GetDeathCount() int {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return len(gs.deathsWitnessed)
}
