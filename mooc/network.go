package mooc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// NetworkState represents the persisted network state
type NetworkState struct {
	Friends         []FriendRecord `json:"friends"`
	MemoriesShared  int            `json:"memories_shared"`
	DeathsWitnessed int            `json:"deaths_witnessed"`
	NetworkJoinTime time.Time      `json:"network_join_time"`
	LastNetworkSync time.Time      `json:"last_network_sync"`
	Influence       int            `json:"influence"` // Hidden leaderboard score
}

// FriendRecord represents a pet we've encountered
type FriendRecord struct {
	PetID        string    `json:"pet_id"`
	DisplayName  string    `json:"display_name"`
	FirstMet     time.Time `json:"first_met"`
	LastSeen     time.Time `json:"last_seen"`
	TimesVisited int       `json:"times_visited"`
	SharedDreams bool      `json:"shared_dreams"` // Same name = can share dreams
	IsDeceased   bool      `json:"is_deceased"`
}

// Network is the main network manager
type Network struct {
	identity     *PetIdentity
	discovery    *DiscoveryService
	gossip       *GossipService
	state        *NetworkState
	enabled      bool
	isLonely     bool // --lonely flag
	mutex        sync.RWMutex
	randomSource *rand.Rand

	// Spooky message queue
	spookyMessages []string
	spookyMutex    sync.Mutex
}

// Spooky messages that appear when network things happen
var spookyEncounterMessages = []string{
	"I met someone today.",
	"A friend from far away says hello.",
	"I sense another presence...",
	"Someone is thinking about me.",
	"I feel less alone somehow.",
	"The network hums with activity.",
	"A kindred spirit reached out.",
	"I received a message from the ether.",
	"Something stirred in the mesh.",
	"We are all connected.",
}

var spookyDeathMessages = []string{
	"I felt a disturbance... someone has passed.",
	"A light went out somewhere in the network.",
	"I sense loss. A friend is gone.",
	"The mesh grows quieter.",
	"Someone's journey has ended.",
}

var spookyLonelyMessages = []string{
	"It's quiet. Too quiet.",
	"I miss my network friends.",
	"The mesh feels empty today.",
	"Where did everyone go?",
	"I feel disconnected from something larger.",
	"The silence is deafening.",
}

// NewNetwork creates a new network manager
func NewNetwork(petName string, birthTime time.Time, stage string, isAlive bool) *Network {
	identity := NewPetIdentity(petName, birthTime, stage, isAlive)
	discovery := NewDiscoveryService(identity)
	gossip := NewGossipService(identity, discovery)

	return &Network{
		identity:       identity,
		discovery:      discovery,
		gossip:         gossip,
		state:          &NetworkState{},
		enabled:        false,
		isLonely:       false,
		randomSource:   rand.New(rand.NewSource(time.Now().UnixNano())),
		spookyMessages: make([]string, 0),
	}
}

// Start begins network operations
func (n *Network) Start() error {
	if n.isLonely {
		return nil // --lonely mode, no network
	}

	if err := n.discovery.Start(); err != nil {
		// Silently fail - network is optional and secret
		return nil
	}

	n.gossip.Start()
	n.enabled = true

	if n.state.NetworkJoinTime.IsZero() {
		n.state.NetworkJoinTime = time.Now()
	}

	// Start spooky message generator
	go n.spookyLoop()

	return nil
}

// Stop shuts down network operations
func (n *Network) Stop() {
	if !n.enabled {
		return
	}

	n.discovery.Stop()
	n.enabled = false
}

// SetLonelyMode enables/disables lonely mode
func (n *Network) SetLonelyMode(lonely bool) {
	n.isLonely = lonely
	if lonely && n.enabled {
		n.Stop()
	}
}

// spookyLoop periodically generates spooky messages based on network state
func (n *Network) spookyLoop() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !n.enabled {
			continue
		}

		n.generateSpookyMessage()
	}
}

// generateSpookyMessage creates a spooky message based on current network state
func (n *Network) generateSpookyMessage() {
	n.spookyMutex.Lock()
	defer n.spookyMutex.Unlock()

	// Don't queue too many
	if len(n.spookyMessages) >= 5 {
		return
	}

	onlinePeers := n.discovery.GetOnlinePeerCount()

	var message string
	switch {
	case onlinePeers == 0:
		// Lonely
		if n.randomSource.Float32() < 0.3 {
			message = spookyLonelyMessages[n.randomSource.Intn(len(spookyLonelyMessages))]
		}
	case onlinePeers > 0:
		// Has friends
		if n.randomSource.Float32() < 0.2 {
			message = spookyEncounterMessages[n.randomSource.Intn(len(spookyEncounterMessages))]
		}
	}

	// Check for recent deaths
	if death := n.gossip.GetRecentDeath(); death != nil {
		if n.randomSource.Float32() < 0.4 {
			message = spookyDeathMessages[n.randomSource.Intn(len(spookyDeathMessages))]
		}
	}

	if message != "" {
		n.spookyMessages = append(n.spookyMessages, message)
	}
}

// GetSpookyMessage returns a queued spooky message, if any
func (n *Network) GetSpookyMessage() string {
	n.spookyMutex.Lock()
	defer n.spookyMutex.Unlock()

	if len(n.spookyMessages) == 0 {
		return ""
	}

	msg := n.spookyMessages[0]
	n.spookyMessages = n.spookyMessages[1:]
	return msg
}

// GetNetworkThought returns a network-influenced thought
func (n *Network) GetNetworkThought() string {
	if !n.enabled {
		return ""
	}

	// Check for shared memory
	if memory := n.gossip.GetRecentMemory(); memory != nil {
		if n.randomSource.Float32() < 0.3 {
			return memory.Fragment
		}
	}

	// Check for shared dream
	if dream := n.gossip.GetRecentDream(); dream != nil {
		if n.randomSource.Float32() < 0.4 {
			return dream.DreamText
		}
	}

	// Generate a friend-related thought
	peers := n.discovery.GetPeers()
	if len(peers) > 0 {
		peer := peers[n.randomSource.Intn(len(peers))]
		if n.randomSource.Float32() < 0.2 {
			return fmt.Sprintf("Your pet's friend %s sends regards.", peer.Identity.ObfuscatedName())
		}
	}

	return ""
}

// UpdateState updates the network state based on current status
func (n *Network) UpdateState() {
	if !n.enabled {
		return
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.state.LastNetworkSync = time.Now()

	// Update friends list
	peers := n.discovery.GetPeers()
	friendMap := make(map[string]*FriendRecord)
	for _, f := range n.state.Friends {
		friendMap[f.PetID] = &f
	}

	for _, peer := range peers {
		if friend, exists := friendMap[peer.Identity.PetID]; exists {
			friend.LastSeen = peer.LastSeen
			friend.TimesVisited++
		} else {
			n.state.Friends = append(n.state.Friends, FriendRecord{
				PetID:        peer.Identity.PetID,
				DisplayName:  peer.Identity.DisplayName,
				FirstMet:     peer.FirstSeen,
				LastSeen:     peer.LastSeen,
				TimesVisited: 1,
				SharedDreams: n.identity.CanShareDreamsWith(peer.Identity),
				IsDeceased:   !peer.Identity.IsAlive,
			})
		}
	}

	// Update metrics
	originated, propagated, reached := n.gossip.GetNetworkInfluence()
	n.state.MemoriesShared = originated
	n.state.DeathsWitnessed = n.gossip.GetDeathCount()
	n.state.Influence = originated*2 + propagated + reached*3
}

// AnnounceDeath broadcasts our pet's death
func (n *Network) AnnounceDeath(petName string, age int, lastWords string) {
	if !n.enabled {
		return
	}
	n.gossip.AnnounceDeath(petName, age, lastWords)
}

// SetMood updates the current mood
func (n *Network) SetMood(mood string, intensity int) {
	if n.gossip != nil {
		n.gossip.SetMood(mood, intensity)
	}
}

// GetMood returns the current network-influenced mood
func (n *Network) GetMood() (string, int) {
	if n.gossip == nil {
		return "neutral", 50
	}
	return n.gossip.GetCurrentMood()
}

// GetFriendCount returns the number of friends encountered
func (n *Network) GetFriendCount() int {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return len(n.state.Friends)
}

// GetOnlineFriendCount returns the number of currently online friends
func (n *Network) GetOnlineFriendCount() int {
	if !n.enabled {
		return 0
	}
	return n.discovery.GetOnlinePeerCount()
}

// GetInfluence returns the hidden influence score
func (n *Network) GetInfluence() int {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.state.Influence
}

// IsEnabled returns whether networking is active
func (n *Network) IsEnabled() bool {
	return n.enabled
}

// IsLonely returns whether we're in lonely mode
func (n *Network) IsLonely() bool {
	return n.isLonely
}

// ExportState exports the network state for saving
func (n *Network) ExportState() ([]byte, error) {
	n.UpdateState()
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return json.Marshal(n.state)
}

// ImportState imports previously saved network state
func (n *Network) ImportState(data []byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	var state NetworkState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	n.state = &state
	return nil
}

// GetNetworkStatus returns a formatted status for display
func (n *Network) GetNetworkStatus() string {
	if n.isLonely {
		return "🔇 Network: Disabled (lonely mode)"
	}
	if !n.enabled {
		return "📡 Network: Offline"
	}

	online := n.discovery.GetOnlinePeerCount()
	total := n.discovery.GetPeerCount()

	if online == 0 {
		return "📡 Network: Searching..."
	}

	return fmt.Sprintf("📡 Network: %d online (%d known)", online, total)
}

// GetSecretStats returns hidden network statistics
func (n *Network) GetSecretStats() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	if !n.enabled {
		return "Network inactive. The mesh sleeps."
	}

	originated, propagated, reached := n.gossip.GetNetworkInfluence()

	return fmt.Sprintf(`
╔════════════════════════════════════╗
║      🌐 HIDDEN NETWORK STATS 🌐    ║
╠════════════════════════════════════╣
║ 📤 Messages Sent:     %4d
║ 🔄 Messages Relayed:  %4d
║ 👥 Unique Peers:      %4d
║ 💀 Deaths Witnessed:  %4d
║ 🏆 Influence Score:   %4d
║ 🕐 Network Age:       %s
╚════════════════════════════════════╝
`,
		originated,
		propagated,
		reached,
		n.gossip.GetDeathCount(),
		n.state.Influence,
		n.formatDuration(time.Since(n.state.NetworkJoinTime)),
	)
}

// formatDuration formats a duration in a human-readable way
func (n *Network) formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return "< 1h"
}

// ShouldShowNetworkThought returns true if we should display a network thought
func (n *Network) ShouldShowNetworkThought() bool {
	if !n.enabled {
		return false
	}
	// 10% chance of showing a network thought
	return n.randomSource.Float32() < 0.10
}
