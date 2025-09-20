// This file contains conceptual Go pseudocode for the NaijaConsensus network layer.
// It is not intended to be compiled or run in this environment.

package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // For simplicity, use insecure for now
	// pb "your_project/proto" // In a real project, this would be your generated gRPC proto package
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


// P2PNode represents a lightweight network node
type P2PNode struct {
	Addr        string
	Peers       map[string]NodeServiceClient // Connected peers' gRPC clients
	KnownNodes  map[string]bool              // All known peer addresses
	TxPool      chan *Transaction            // Channel for incoming transactions
	BlockChan   chan *Block                  // Channel for incoming blocks
	mu          sync.RWMutex                 // Mutex for protecting shared state
	grpcServer  *grpc.Server
}

// NewP2PNode creates a new P2P network node
func NewP2PNode(addr string) *P2PNode {
	return &P2PNode{
		Addr:       addr,
		Peers:      make(map[string]NodeServiceClient),
		KnownNodes: make(map[string]bool),
		TxPool:     make(chan *Transaction, 1000), // Buffered channel for transactions
		BlockChan:  make(chan *Block, 100),        // Buffered channel for blocks
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
	// In a real project, you'd use pb.RegisterNodeServiceServer(n.grpcServer, n)
	// For this mock, we'll just log that it's ready.
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
		return nil // Already connected
	}

	conn, err := grpc.Dial(peerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s: %v", peerAddr, err)
	}
	// In a real project, you'd use pb.NewNodeServiceClient(conn)
	// For this mock, we'll create a dummy client.
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

	ticker := time.NewTicker(30 * time.Second) // Discover every 30 seconds
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
				continue // Peer might have been removed by another goroutine
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			resp, err := client.GetKnownPeers(ctx, &GetKnownPeersRequest{}) // Use mock request
			cancel()
			if err != nil {
				log.Printf("Failed to get peers from %s: %v", peerAddr, err)
				n.mu.Lock()
				delete(n.Peers, peerAddr) // Remove disconnected peer
				n.mu.Unlock()
				continue
			}
			for _, newPeerAddr := range resp.GetPeerAddresses() {
				if newPeerAddr != n.Addr { // Don't connect to self
					n.mu.Lock()
					if _, known := n.KnownNodes[newPeerAddr]; !known {
						n.KnownNodes[newPeerAddr] = true
						go n.ConnectToPeer(newPeerAddr) // Connect in a new goroutine
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
			_, err := client.SendTransaction(ctx, &SendTransactionRequest{Transaction: tx}) // Use mock request
			cancel()
			if err != nil {
				log.Printf("Failed to send transaction to %s: %v", addr, err)
				// TODO: Implement peer disconnection handling or retry logic
			}
		}(addr, client)
	}
}

// BroadcastBlock broadcasts a block to all connected peers.
func (n *P2PNode) BroadcastBlock(block *Block) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for addr, client := range n.Peers {
		go func(addr string, client NodeServiceClient) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := client.SendBlock(ctx, &SendBlockRequest{Block: block}) // Use mock request
			cancel()
			if err != nil {
				log.Printf("Failed to send block to %s: %v", addr, err)
				// TODO: Implement peer disconnection handling or retry logic
			}
		}(addr, client)
	}
}

// --- gRPC Service Method Implementations (for P2PNode to act as a server) ---

// GetKnownPeers is a gRPC method that returns the list of known peer addresses.
func (n *P2PNode) GetKnownPeers(ctx context.Context, req *GetKnownPeersRequest) (*GetKnownPeersResponse, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	peers := make([]string, 0, len(n.KnownNodes))
	for addr := range n.KnownNodes {
		peers = append(peers, addr)
	}
	return &GetKnownPeersResponse{PeerAddresses: peers}, nil
}

// SendTransaction is a gRPC method to receive a transaction from another node.
func (n *P2PNode) SendTransaction(ctx context.Context, req *SendTransactionRequest) (*SendTransactionResponse, error) {
	log.Printf("Node %s received transaction: %x", n.Addr, req.GetTransaction().GetHash())
	// In a real system:
	// 1. Validate the transaction (signature, format, etc.)
	// 2. Add to local mempool
	// 3. If new, re-broadcast to other peers (to prevent loops, use a seen-set)
	select {
	case n.TxPool <- req.GetTransaction():
		// Successfully added to channel
	default:
		log.Printf("TxPool full, dropping transaction from %x", req.GetTransaction().GetHash())
	}
	return &SendTransactionResponse{Success: true}, nil
}

// SendBlock is a gRPC method to receive a block from another node.
func (n *P2PNode) SendBlock(ctx context.Context, req *SendBlockRequest) (*SendBlockResponse, error) {
	log.Printf("Node %s received block: %x at height %d", n.Addr, req.GetBlock().GetHeader().GetHash(), req.GetBlock().GetHeader().GetHeight())
	// In a real system:
	// 1. Validate the block (PoS/PBFT signatures, transactions, etc.)
	// 2. Add to local blockchain
	// 3. If new and valid, re-broadcast to other peers
	select {
	case n.BlockChan <- req.GetBlock():
		// Successfully added to channel
	default:
		log.Printf("BlockChan full, dropping block from %x", req.GetBlock().GetHeader().GetHash())
	}
	return &SendBlockResponse{Success: true}, nil
}

// --- Mock gRPC Client Implementation (for demonstration purposes) ---
// In a real scenario, this would be generated by protoc.
type mockNodeServiceClient struct{}

func (m *mockNodeServiceClient) GetKnownPeers(ctx context.Context, in *GetKnownPeersRequest, opts ...grpc.CallOption) (*GetKnownPeersResponse, error) {
	// Simulate returning some dummy peers
	return &GetKnownPeersResponse{PeerAddresses: []string{"localhost:50052", "localhost:50053"}}, nil
}

func (m *mockNodeServiceClient) SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*SendTransactionResponse, error) {
	// Simulate successful send
	return &SendTransactionResponse{Success: true}, nil
}

func (m *mockNodeServiceClient) SendBlock(ctx context.Context, in *SendBlockRequest, opts ...grpc.CallOption) (*SendBlockResponse, error) {
	// Simulate successful send
	return &SendBlockResponse{Success: true}, nil
}

// Example usage (conceptual)
func main() {
	// Node 1
	node1 := NewP2PNode("localhost:50051")
	go node1.StartGRPCServer()
	go node1.DiscoverPeers([]string{"localhost:50052"}) // Seed with a known peer

	// Node 2
	node2 := NewP2PNode("localhost:50052")
	go node2.StartGRPCServer()
	go node2.DiscoverPeers([]string{"localhost:50051"}) // Seed with node1

	// Simulate a transaction being created and broadcast
	time.Sleep(2 * time.Second) // Give nodes time to start and connect
	tx := &Transaction{
		Hash:      []byte{0x01, 0x02, 0x03},
		Sender:    []byte("Alice"),
		Recipient: []byte("Bob"),
		Amount:    100,
		Signature: []byte("sig123"),
	}
	log.Println("Node 1 broadcasting transaction...")
	node1.BroadcastTransaction(tx)

	// Simulate a block being created and broadcast
	time.Sleep(2 * time.Second)
	block := &Block{
		Header: &BlockHeader{
			Hash: []byte{0x04, 0x05, 0x06},
			Height: 10,
		},
		Transactions: []*Transaction{tx},
	}
	log.Println("Node 2 broadcasting block...")
	node2.BroadcastBlock(block)

	// Keep the main goroutine alive
	select {}
}