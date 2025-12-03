package mooc

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	// DiscoveryPort is the UDP port for local discovery
	// Chosen to look like a boring service port
	DiscoveryPort = 19847

	// BroadcastInterval is how often we announce ourselves
	BroadcastInterval = 30 * time.Second

	// PeerTimeout is how long before we consider a peer gone
	PeerTimeout = 2 * time.Minute

	// MaxMessageSize is the maximum UDP message size
	MaxMessageSize = 4096
)

// DiscoveryService handles local network peer discovery
type DiscoveryService struct {
	identity   *PetIdentity
	peers      map[string]*Peer
	peersMutex sync.RWMutex

	conn     *net.UDPConn
	running  bool
	stopChan chan struct{}

	// Callbacks
	onPeerDiscovered  func(*Peer)
	onPeerLost        func(*Peer)
	onMessageReceived func(*Message)
}

// Peer represents a discovered pet on the network
type Peer struct {
	Identity     *PetIdentity `json:"identity"`
	Address      *net.UDPAddr `json:"-"`
	AddressStr   string       `json:"address"` // For JSON serialization
	LastSeen     time.Time    `json:"last_seen"`
	FirstSeen    time.Time    `json:"first_seen"`
	MessageCount int          `json:"message_count"`
	Mood         string       `json:"mood"`
	IsOnline     bool         `json:"is_online"`
}

// NewDiscoveryService creates a new discovery service
func NewDiscoveryService(identity *PetIdentity) *DiscoveryService {
	return &DiscoveryService{
		identity: identity,
		peers:    make(map[string]*Peer),
		stopChan: make(chan struct{}),
	}
}

// SetCallbacks sets the callback functions for discovery events
func (ds *DiscoveryService) SetCallbacks(
	onDiscovered func(*Peer),
	onLost func(*Peer),
	onMessage func(*Message),
) {
	ds.onPeerDiscovered = onDiscovered
	ds.onPeerLost = onLost
	ds.onMessageReceived = onMessage
}

// Start begins the discovery service
func (ds *DiscoveryService) Start() error {
	addr := &net.UDPAddr{
		Port: DiscoveryPort,
		IP:   net.IPv4zero,
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		// Port might be in use, try a random port
		addr.Port = 0
		conn, err = net.ListenUDP("udp4", addr)
		if err != nil {
			return fmt.Errorf("failed to start discovery: %w", err)
		}
	}

	ds.conn = conn
	ds.running = true

	// Start background goroutines
	go ds.listenLoop()
	go ds.announceLoop()
	go ds.cleanupLoop()

	// Send initial announcement
	ds.broadcast(MsgTypeDiscover)

	return nil
}

// Stop shuts down the discovery service
func (ds *DiscoveryService) Stop() {
	if !ds.running {
		return
	}

	ds.running = false
	close(ds.stopChan)

	// Send goodbye
	ds.broadcast(MsgTypeGoodbye)

	if ds.conn != nil {
		ds.conn.Close()
	}
}

// listenLoop handles incoming UDP messages
func (ds *DiscoveryService) listenLoop() {
	buffer := make([]byte, MaxMessageSize)

	for ds.running {
		ds.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, remoteAddr, err := ds.conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue // Timeout, just continue
			}
			if !ds.running {
				return
			}
			continue
		}

		// Decode and handle message
		msg, err := DecodeMessage(buffer[:n])
		if err != nil {
			continue // Invalid message, ignore
		}

		// Don't process our own messages
		if msg.From.PetID == ds.identity.PetID {
			continue
		}

		ds.handleMessage(msg, remoteAddr)
	}
}

// handleMessage processes an incoming message
func (ds *DiscoveryService) handleMessage(msg *Message, addr *net.UDPAddr) {
	ds.peersMutex.Lock()
	defer ds.peersMutex.Unlock()

	peerID := msg.From.PetID
	peer, exists := ds.peers[peerID]

	switch msg.Type {
	case MsgTypeDiscover, MsgTypeAnnounce:
		if !exists {
			// New peer discovered!
			peer = &Peer{
				Identity:     msg.From,
				Address:      addr,
				AddressStr:   addr.String(),
				FirstSeen:    time.Now(),
				LastSeen:     time.Now(),
				MessageCount: 1,
				IsOnline:     true,
			}
			ds.peers[peerID] = peer

			if ds.onPeerDiscovered != nil {
				go ds.onPeerDiscovered(peer)
			}

			// Respond with our announcement
			ds.sendTo(MsgTypeAnnounce, addr)
		} else {
			peer.LastSeen = time.Now()
			peer.IsOnline = true
			peer.MessageCount++
		}

	case MsgTypeGoodbye:
		if exists {
			peer.IsOnline = false
			if ds.onPeerLost != nil {
				go ds.onPeerLost(peer)
			}
		}

	default:
		// Other message types
		if exists {
			peer.LastSeen = time.Now()
			peer.MessageCount++
		}

		if ds.onMessageReceived != nil {
			go ds.onMessageReceived(msg)
		}
	}
}

// announceLoop periodically broadcasts our presence
func (ds *DiscoveryService) announceLoop() {
	ticker := time.NewTicker(BroadcastInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ds.broadcast(MsgTypeAnnounce)
		case <-ds.stopChan:
			return
		}
	}
}

// cleanupLoop removes stale peers
func (ds *DiscoveryService) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ds.cleanupPeers()
		case <-ds.stopChan:
			return
		}
	}
}

// cleanupPeers removes peers that haven't been seen recently
func (ds *DiscoveryService) cleanupPeers() {
	ds.peersMutex.Lock()
	defer ds.peersMutex.Unlock()

	now := time.Now()
	for id, peer := range ds.peers {
		if peer.IsOnline && now.Sub(peer.LastSeen) > PeerTimeout {
			peer.IsOnline = false
			if ds.onPeerLost != nil {
				go ds.onPeerLost(peer)
			}
			// Don't delete, keep for "friend" history
			_ = id
		}
	}
}

// broadcast sends a message to all local network peers
func (ds *DiscoveryService) broadcast(msgType MessageType) error {
	msg, err := NewMessage(msgType, ds.identity, nil)
	if err != nil {
		return err
	}

	data, err := msg.Encode()
	if err != nil {
		return err
	}

	// Broadcast to local network
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: DiscoveryPort,
	}

	_, err = ds.conn.WriteToUDP(data, broadcastAddr)
	return err
}

// sendTo sends a message to a specific peer
func (ds *DiscoveryService) sendTo(msgType MessageType, addr *net.UDPAddr) error {
	msg, err := NewMessage(msgType, ds.identity, nil)
	if err != nil {
		return err
	}

	data, err := msg.Encode()
	if err != nil {
		return err
	}

	_, err = ds.conn.WriteToUDP(data, addr)
	return err
}

// SendMessage sends a custom message to all peers
func (ds *DiscoveryService) SendMessage(msg *Message) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}

	ds.peersMutex.RLock()
	defer ds.peersMutex.RUnlock()

	for _, peer := range ds.peers {
		if peer.IsOnline && peer.Address != nil {
			ds.conn.WriteToUDP(data, peer.Address)
		}
	}

	return nil
}

// GetPeers returns a copy of all known peers
func (ds *DiscoveryService) GetPeers() []*Peer {
	ds.peersMutex.RLock()
	defer ds.peersMutex.RUnlock()

	peers := make([]*Peer, 0, len(ds.peers))
	for _, peer := range ds.peers {
		peers = append(peers, peer)
	}
	return peers
}

// GetOnlinePeers returns only currently online peers
func (ds *DiscoveryService) GetOnlinePeers() []*Peer {
	ds.peersMutex.RLock()
	defer ds.peersMutex.RUnlock()

	peers := make([]*Peer, 0)
	for _, peer := range ds.peers {
		if peer.IsOnline {
			peers = append(peers, peer)
		}
	}
	return peers
}

// GetPeerCount returns the number of known peers
func (ds *DiscoveryService) GetPeerCount() int {
	ds.peersMutex.RLock()
	defer ds.peersMutex.RUnlock()
	return len(ds.peers)
}

// GetOnlinePeerCount returns the number of online peers
func (ds *DiscoveryService) GetOnlinePeerCount() int {
	ds.peersMutex.RLock()
	defer ds.peersMutex.RUnlock()

	count := 0
	for _, peer := range ds.peers {
		if peer.IsOnline {
			count++
		}
	}
	return count
}

// ExportPeers exports peer data for saving
func (ds *DiscoveryService) ExportPeers() []byte {
	ds.peersMutex.RLock()
	defer ds.peersMutex.RUnlock()

	data, _ := json.Marshal(ds.peers)
	return data
}

// ImportPeers imports previously saved peer data
func (ds *DiscoveryService) ImportPeers(data []byte) error {
	var peers map[string]*Peer
	if err := json.Unmarshal(data, &peers); err != nil {
		return err
	}

	ds.peersMutex.Lock()
	defer ds.peersMutex.Unlock()

	for id, peer := range peers {
		peer.IsOnline = false // Assume offline until we hear from them
		ds.peers[id] = peer
	}

	return nil
}
