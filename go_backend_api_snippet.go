// go_backend_api_snippet.go

package main

import (
	"context"
	"crypto/sha3"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// "your_project/proto" // In a real project, this would be your generated gRPC proto package
)

// --- Mock gRPC Protobuf Definitions (replace with actual generated code) ---
// These structs mimic the generated gRPC types for demonstration.
type Transaction struct {
	Hash      []byte
	Sender    []byte
	Recipient []byte
	Amount    uint64
	Signature []byte
}

type BlockHeader struct {
	Version       uint32
	PrevBlockHash []byte
	MerkleRoot    []byte
	Timestamp     uint64
	Height        uint64
}

type Block struct {
	Header *BlockHeader
	Transactions []*Transaction
}

type GetKnownPeersRequest struct{}
type GetKnownPeersResponse struct {
	PeerAddresses []string
}

type SendTransactionRequest struct {
	Transaction *Transaction
}
type SendTransactionResponse struct {
	Success bool
}

type SendBlockRequest struct {
	Block *Block
}
type SendBlockResponse struct {
	Success bool
}

// NodeServiceServer interface (mimics generated gRPC server interface)
type NodeServiceServer interface {
	GetKnownPeers(context.Context, *GetKnownPeersRequest) (*GetKnownPeersResponse, error)
	SendTransaction(context.Context, *SendTransactionRequest) (*SendTransactionResponse, error)
	SendBlock(context.Context, *SendBlockRequest) (*SendBlockResponse, error)
}

// NodeServiceClient interface (mimics generated gRPC client interface)
type NodeServiceClient interface {
	GetKnownPeers(ctx context.Context, in *GetKnownPeersRequest, opts ...grpc.CallOption) (*GetKnownPeersResponse, error)
	SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*SendTransactionResponse, error)
	SendBlock(ctx context.Context, in *SendBlockRequest, opts ...grpc.CallOption) (*SendBlockResponse, error)
}

// --- End Mock gRPC Protobuf Definitions ---

// P2PNode represents a lightweight network node for NaijaVote
type P2PNode struct {
	Addr        string
	Peers       map[string]NodeServiceClient
	KnownNodes  map[string]bool
	TxPool      chan *Transaction // For incoming vote transactions
	BlockChan   chan *Block       // For incoming blocks
	mu          sync.RWMutex
	grpcServer  *grpc.Server
	// Mock Rust Consensus Engine interaction
	// rustEngine *consensus.NaijaConsensusEngine // Conceptual link
}

// NewP2PNode creates a new P2P network node
func NewP2PNode(addr string) *P2PNode {
	return &P2PNode{
		Addr:       addr,
		Peers:      make(map[string]NodeServiceClient),
		KnownNodes: make(map[string]bool),
		TxPool:     make(chan *Transaction, 1000),
		BlockChan:  make(chan *Block, 100),
	}
}

// StartGRPCServer starts the gRPC server for the node.
// This method should be run in a goroutine.
func (n *P2PNode) StartGRPCServer() {
	lis, err := net.Listen("tcp", n.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	n.grpcServer = grpc.NewServer()
	// In a real project: pb.RegisterNodeServiceServer(n.grpcServer, n)
	log.Printf("gRPC server listening on %s", n.Addr)
	if err := n.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// ConnectToPeer establishes a gRPC connection to another peer.
func (n *P2PNode) ConnectToPeer(peerAddr string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, ok := n.Peers[peerAddr]; ok {
		return nil
	}

	conn, err := grpc.Dial(peerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s: %v", peerAddr, err)
	}
	client := &mockNodeServiceClient{} // Replace with actual gRPC client
	n.Peers[peerAddr] = client
	n.KnownNodes[peerAddr] = true
	log.Printf("Connected to peer: %s", peerAddr)
	return nil
}

// DiscoverPeers periodically discovers and connects to new peers.
// This method should be run in a goroutine.
func (n *P2PNode) DiscoverPeers(initialPeers []string) {
	for _, peer := range initialPeers {
		n.KnownNodes[peer] = true
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		n.mu.RLock()
		peersToQuery := make([]string, 0, len(n.Peers))
		for addr := range n.Peers {
			peersToQuery = append(peersToQuery, addr)
		}
		n.mu.RUnlock()

		for _, peerAddr := range peersToQuery {
			client, ok := n.Peers[peerAddr]
			if !ok {
				continue
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			resp, err := client.GetKnownPeers(ctx, &GetKnownPeersRequest{})
			cancel()
			if err != nil {
				log.Printf("Failed to get peers from %s: %v", peerAddr, err)
				n.mu.Lock()
				delete(n.Peers, peerAddr)
				n.mu.Unlock()
				continue
			}
			for _, newPeerAddr := range resp.GetPeerAddresses() {
				if newPeerAddr != n.Addr {
					n.mu.Lock()
					if _, known := n.KnownNodes[newPeerAddr]; !known {
						n.KnownNodes[newPeerAddr] = true
						go n.ConnectToPeer(newPeerAddr)
					}
					n.mu.Unlock()
				}
			}
		}
	}
}

// BroadcastTransaction broadcasts a transaction to all connected peers.
func (n *P2PNode) BroadcastTransaction(tx *Transaction) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for addr, client := range n.Peers {
		go func(addr string, client NodeServiceClient) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err := client.SendTransaction(ctx, &SendTransactionRequest{Transaction: tx})
			cancel()
			if err != nil {
				log.Printf("Failed to send transaction to %s: %v", addr, err)
			}
		}(addr, client)
	}
}

// --- gRPC Service Method Implementations (for P2PNode to act as a server) ---

func (n *P2PNode) GetKnownPeers(ctx context.Context, req *GetKnownPeersRequest) (*GetKnownPeersResponse, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	peers := make([]string, 0, len(n.KnownNodes))
	for addr := range n.KnownNodes {
		peers = append(peers, addr)
	}
	return &GetKnownPeersResponse{PeerAddresses: peers}, nil
}

func (n *P2PNode) SendTransaction(ctx context.Context, req *SendTransactionRequest) (*SendTransactionResponse, error) {
	log.Printf("Node %s received transaction: %x", n.Addr, req.GetTransaction().GetHash())
	// In a real system: Validate, add to mempool, re-broadcast if new.
	// Then, the full node would pass this to the Rust consensus engine.
	select {
	case n.TxPool <- req.GetTransaction():
	default:
		log.Printf("TxPool full, dropping transaction from %x", req.GetTransaction().GetHash())
	}
	return &SendTransactionResponse{Success: true}, nil
}

func (n *P2PNode) SendBlock(ctx context.Context, req *SendBlockRequest) (*SendBlockResponse, error) {
	log.Printf("Node %s received block: %x at height %d", n.Addr, req.GetBlock().GetHeader().GetHash(), req.GetBlock().GetHeader().GetHeight())
	// In a real system: Validate block using Rust consensus engine, add to chain, re-broadcast.
	select {
	case n.BlockChan <- req.GetBlock():
	default:
		log.Printf("BlockChan full, dropping block from %x", req.GetBlock().GetHeader().GetHash())
	}
	return &SendBlockResponse{Success: true}, nil
}

// --- Mock gRPC Client Implementation ---
type mockNodeServiceClient struct{}

func (m *mockNodeServiceClient) GetKnownPeers(ctx context.Context, in *GetKnownPeersRequest, opts ...grpc.CallOption) (*GetKnownPeersResponse, error) {
	return &GetKnownPeersResponse{PeerAddresses: []string{"localhost:50052", "localhost:50053"}}, nil
}

func (m *mockNodeServiceClient) SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*SendTransactionResponse, error) {
	return &SendTransactionResponse{Success: true}, nil
}

func (m *mockNodeServiceClient) SendBlock(ctx context.Context, in *SendBlockRequest, opts ...grpc.CallOption) (*SendBlockResponse, error) {
	return &SendBlockResponse{Success: true}, nil
}

// --- HTTP API Handlers for Frontend Interaction ---

// Mock voter registry (off-chain, for demonstration)
var voterRegistry = make(map[string]string) // NIN/BVN -> HashedPassword

func hashNINBVN(ninBvn string) string {
	hasher := sha3.New256()
	hasher.Write([]byte(ninBvn))
	return hex.EncodeToString(hasher.Sum(nil))
}

// RegisterVoter handles voter registration requests
func RegisterVoter(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		NIN_BVN  string `json:"nin_bvn"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedNINBVN := hashNINBVN(req.NIN_BVN)
	// In a real system: Verify NIN/BVN with NIMC, then generate and sign voting token.
	// For now, just store a mock.
	voterRegistry[hashedNINBVN] = "mock_hashed_password" // Store hashed password
	log.Printf("Voter registered: %s (hashed)", hashedNINBVN)

	// Simulate generating a cryptographically signed voting token
	votingToken := fmt.Sprintf("VOTETOKEN_%s_%d", hashedNINBVN, time.Now().Unix())

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Voter registered successfully",
		"voting_token": votingToken,
	})
}

// SubmitVote handles vote submission requests
func SubmitVote(node *P2PNode, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		VoterID    string `json:"voter_id"` // This would be the voting token or public key
		ElectionID string `json:"election_id"`
		Candidate  string `json:"candidate"`
		Signature  string `json:"signature"` // Signed transaction by the client
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// In a real system:
	// 1. Verify signature using Ed25519.
	// 2. Check if voter_id (voting token) is valid and hasn't voted.
	// 3. Create a VoteTransaction struct.
	// 4. Pass to P2PNode to broadcast.
	log.Printf("Received vote from %s for %s in election %s", req.VoterID, req.Candidate, req.ElectionID)

	// Simulate creating a blockchain transaction
	txHash := sha3.New256()
	txHash.Write([]byte(fmt.Sprintf("%s%s%s", req.VoterID, req.ElectionID, req.Candidate)))
	mockTx := &Transaction{
		Hash:      txHash.Sum(nil),
		Sender:    []byte(req.VoterID),
		Recipient: []byte(req.Candidate),
		Amount:    1, // Represents one vote
		Signature: []byte(req.Signature),
	}

	node.BroadcastTransaction(mockTx) // Broadcast to other P2P nodes

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vote submitted and broadcasted successfully. Awaiting blockchain finality.",
		"tx_hash": hex.EncodeToString(mockTx.Hash),
	})
}

// GetElectionStatus provides real-time election data
func GetElectionStatus(w http.ResponseWriter, r *http.Request) {
	// In a real system: Query the local blockchain state (managed by Rust engine)
	// For now, return mock data.
	status := map[string]interface{}{
		"total_votes":   12345,
		"candidates": map[string]int{
			"Candidate A": 5000,
			"Candidate B": 4000,
			"Candidate C": 3345,
		},
		"latest_block_hash": "0x" + hex.EncodeToString([]byte{0xab, 0xcd, 0xef, 0x12, 0x34, 0x56, 0x78, 0x90}) + "...",
		"block_height":      123456,
		"finality_time_seconds": 3,
		"validators_active": 21,
	}
	json.NewEncoder(w).Encode(status)
}

func main() {
	// Initialize P2P Node (conceptual)
	p2pNode := NewP2PNode("localhost:50051")
	go p2pNode.StartGRPCServer()
	go p2pNode.DiscoverPeers([]string{"localhost:50052"}) // Seed with a dummy peer

	// Start HTTP API Server
	http.HandleFunc("/register", RegisterVoter)
	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		SubmitVote(p2pNode, w, r)
	})
	http.HandleFunc("/status", GetElectionStatus)

	log.Println("HTTP API server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
</pre>

  <bindAction type="shell">npm install