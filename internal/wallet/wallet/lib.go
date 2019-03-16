package wallet

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	clientService "github.com/harmony-one/harmony/api/client/service"
	"github.com/harmony-one/harmony/core/types"
	"github.com/harmony-one/harmony/node"
	"github.com/harmony-one/harmony/p2p"
)

// we will hard coded several rpc servers (archival nodes) for wallet
var (
	p1 = p2p.Peer{IP: "127.0.0.1", Port: "9001"}
	p2 = p2p.Peer{IP: "127.0.0.1", Port: "9002"}
)

// AccountState includes the balance and nonce of an account
type AccountState struct {
	Balance *big.Int
	Nonce   uint64
}

type Wallet struct {
	selfPeer p2p.Peer
	peers    map[uint32][]p2p.Peer // we don't use p2p, use grpc instead
	shardIds []uint32
}

// CreateWalletNode creates wallet server node.
// Peer information is hard coded
func New() *Wallet {
	// dummy host for wallet
	self := p2p.Peer{IP: "127.0.0.1", Port: "6999"}
	wallet := &Wallet{selfPeer: self, peers: createAllPeers()}
	return wallet
}

// createAllPeers get peers from given shard
// TODO (chao): only support beacon shard now; add support for normal shards
func createAllPeers() map[uint32][]p2p.Peer {
	peers := make(map[uint32][]p2p.Peer)
	peers[0] = []p2p.Peer{p1, p2}
	return peers
}

// GetPeersFromShardID gets all the peers for a given ShardID
func (wallet *Wallet) GetPeersFromShardID(shardID uint32) []p2p.Peer {
	return wallet.peers[shardID]
}

// pickRandomPeer will pick a random peer from a given shardID
func (wallet *Wallet) pickRandomPeer(shardID uint32) p2p.Peer {
	idx := rand.Intn(len(wallet.peers[shardID]))
	return wallet.peers[shardID][idx]
}

// FetchBalance fetches account balance of specified address from the Harmony network
func (wallet *Wallet) FetchBalance(address common.Address, shardIDs []uint32) map[uint32]AccountState {
	result := make(map[uint32]AccountState)
	for _, shardID := range shardIDs {
		peer := wallet.pickRandomPeer(shardID)
		port, _ := strconv.Atoi(peer.Port)
		client := clientService.NewClient(peer.IP, strconv.Itoa(port+node.ClientServicePortDiff))
		response := client.GetBalance(address)
		balance := big.NewInt(0)
		balance.SetBytes(response.Balance)
		result[shardID] = AccountState{balance, response.Nonce}
	}
	return result
}

// GetFreeToken requests for token test token on each shard
func (wallet *Wallet) GetFreeToken(address common.Address, shardID uint32) {
	peer := wallet.pickRandomPeer(shardID)
	port, _ := strconv.Atoi(peer.Port)
	client := clientService.NewClient(peer.IP, strconv.Itoa(port+node.ClientServicePortDiff))
	response := client.GetFreeToken(address)

	txID := common.Hash{}
	txID.SetBytes(response.TxId)
	fmt.Printf("Transaction Id requesting free token in shard %d: %s\n", int(0), txID.Hex())
}

func (wallet *Wallet) SubmitTransactions(txs []*types.Transaction, shardID uint32) {
	peer := wallet.pickRandomPeer(shardID)
	port, _ := strconv.Atoi(peer.Port)
	client := clientService.NewClient(peer.IP, strconv.Itoa(port+node.ClientServicePortDiff))
	client.SubmitTransactions(txs)
}
