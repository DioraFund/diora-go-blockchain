package core

import (
	"math/big"
	"time"
)

// Common types used across the blockchain

type Address [20]byte

func (a Address) Hex() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a Address) Bytes() []byte {
	return a[:]
}

func HexToAddress(s string) Address {
	if len(s) >= 2 && s[0:2] == "0x" {
		s = s[2:]
	}
	if len(s) != 40 {
		return Address{}
	}
	
	var a Address
	hex.Decode(a[:], []byte(s))
	return a
}

type Hash [32]byte

func (h Hash) Hex() string {
	return "0x" + hex.EncodeToString(h[:])
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func HexToHash(s string) Hash {
	if len(s) >= 2 && s[0:2] == "0x" {
		s = s[2:]
	}
	if len(s) != 64 {
		return Hash{}
	}
	
	var h Hash
	hex.Decode(h[:], []byte(s))
	return h
}

type Receipt struct {
	TransactionHash Hash
	TransactionIndex uint64
	BlockHash       Hash
	BlockNumber     *big.Int
	CumulativeGasUsed uint64
	GasUsed         uint64
	ContractAddress Address
	Logs            []*Log
	Status          uint64 // 1 = success, 0 = failure
}

type Log struct {
	Address     Address
	Topics      []Hash
	Data        []byte
	BlockNumber uint64
	TxHash      Hash
	TxIndex     uint32
	Index       uint32
	Removed     bool
}

type Bloom [256]byte

func (b Bloom) Set(topic Hash) {
	// Simplified bloom filter implementation
	// In production, use proper bloom filter with multiple hash functions
}

func (b Bloom) Test(topic Hash) bool {
	// Simplified bloom filter test
	return false
}

type Trie struct {
	root   Hash
	db     Database
	cache  map[Hash][]byte
}

func NewTrie(db Database) *Trie {
	return &Trie{
		root:  Hash{},
		db:    db,
		cache: make(map[Hash][]byte),
	}
}

func (t *Trie) Root() Hash {
	return t.root
}

func (t *Trie) Get(key []byte) ([]byte, error) {
	// Simplified trie get implementation
	// In production, implement proper Patricia trie
	return nil, nil
}

func (t *Trie) Put(key, value []byte) error {
	// Simplified trie put implementation
	// In production, implement proper Patricia trie
	return nil
}

func (t *Trie) Delete(key []byte) error {
	// Simplified trie delete implementation
	// In production, implement proper Patricia trie
	return nil
}

type Database interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	Close() error
}

type LevelDBDatabase struct {
	db *leveldb.DB
}

func NewLevelDBDatabase(path string) (*LevelDBDatabase, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBDatabase{db: db}, nil
}

func (db *LevelDBDatabase) Get(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

func (db *LevelDBDatabase) Put(key, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *LevelDBDatabase) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *LevelDBDatabase) Has(key []byte) (bool, error) {
	_, err := db.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return false, nil
	}
	return err == nil, err
}

func (db *LevelDBDatabase) Close() error {
	return db.db.Close()
}

// Network types
type PeerInfo struct {
	ID          string
	Address     string
	LastSeen    time.Time
	ChainHeight *big.Int
	Version     string
}

type NetworkStats struct {
	PeerCount      int
	TotalPeers     int
	ChainHeight    *big.Int
	ChainTip       Hash
	GasPrice       *big.Int
	BlockTime      time.Duration
	Tps            float64
	ActiveValidators int
}

// Validator types
type Validator struct {
	Address     Address
	Stake       *big.Int
	Commission  uint64
	Status      ValidatorStatus
	LastActive  time.Time
	TotalBlocks uint64
	Reward      *big.Int
}

type ValidatorStatus int

const (
	ValidatorStatusInactive ValidatorStatus = iota
	ValidatorStatusActive
	ValidatorStatusSlashed
)

// Staking types
type StakingInfo struct {
	Validator   Address
	Delegator   Address
	Amount      *big.Int
	Reward      *big.Int
	LockPeriod  uint64
	StartTime   time.Time
	EndTime     time.Time
	Status      StakingStatus
}

type StakingStatus int

const (
	StakingStatusActive StakingStatus = iota
	StakingStatusUnbonding
	StakingStatusCompleted
	StakingStatusSlashed
)

// Governance types
type Proposal struct {
	ID          uint64
	Proposer    Address
	Title       string
	Description string
	Type        ProposalType
	Value       *big.Int
	StartTime   time.Time
	EndTime     time.Time
	Status      ProposalStatus
	Votes       map[Address]bool
	TotalVotes  uint64
	YesVotes    uint64
	NoVotes     uint64
}

type ProposalType int

const (
	ProposalTypeText ProposalType = iota
	ProposalTypeParameterChange
	ProposalTypeUpgrade
	ProposalTypeSpending
)

type ProposalStatus int

const (
	ProposalStatusPending ProposalStatus = iota
	ProposalStatusActive
	ProposalStatusPassed
	ProposalStatusRejected
	ProposalStatusExecuted
)

// Token types
type TokenInfo struct {
	Address     Address
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply *big.Int
	Holders     uint64
	Transfers   uint64
}

type NFTInfo struct {
	Contract    Address
	TokenID     *big.Int
	Owner       Address
	URI         string
	Metadata    map[string]interface{}
	Creator     Address
	CreatedAt   time.Time
}

// Utility functions
func BigToHash(n *big.Int) Hash {
	if n == nil {
		return Hash{}
	}
	return Hash(n.Bytes())
}

func HashToBig(h Hash) *big.Int {
	return new(big.Int).SetBytes(h[:])
}

func BytesToHash(b []byte) Hash {
	var h Hash
	if len(b) > len(h) {
		b = b[len(b)-len(h):]
	}
	copy(h[len(h)-len(b):], b)
	return h
}

func BytesToAddress(b []byte) Address {
	var a Address
	if len(b) > len(a) {
		b = b[len(b)-len(a):]
	}
	copy(a[len(a)-len(b):], b)
	return a
}

// RLP encoding/decoding helpers
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		b.Header.ParentHash,
		b.Header.Coinbase,
		b.Header.StateRoot,
		b.Header.TxRoot,
		b.Header.ReceiptRoot,
		b.Header.Difficulty,
		b.Header.Number,
		b.Header.GasLimit,
		b.Header.GasUsed,
		b.Header.Timestamp,
		b.Header.ExtraData,
		b.Header.MixHash,
		b.Header.Nonce,
		b.Header.Validator,
		b.Header.Signature,
		b.Transactions,
	})
}

func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var header struct {
		ParentHash  Hash
		Coinbase    Address
		StateRoot   Hash
		TxRoot      Hash
		ReceiptRoot Hash
		Difficulty  *big.Int
		Number      *big.Int
		GasLimit    uint64
		GasUsed     uint64
		Timestamp   uint64
		ExtraData   []byte
		MixHash     Hash
		Nonce       [8]byte
		Validator   Address
		Signature   []byte
	}
	
	if err := s.Decode(&header); err != nil {
		return err
	}
	
	b.Header = &BlockHeader{
		ParentHash:  header.ParentHash,
		Coinbase:    header.Coinbase,
		StateRoot:   header.StateRoot,
		TxRoot:      header.TxRoot,
		ReceiptRoot: header.ReceiptRoot,
		Difficulty:  header.Difficulty,
		Number:      header.Number,
		GasLimit:    header.GasLimit,
		GasUsed:     header.GasUsed,
		Timestamp:   header.Timestamp,
		ExtraData:   header.ExtraData,
		MixHash:     header.MixHash,
		Nonce:       header.Nonce,
		Validator:   header.Validator,
		Signature:   header.Signature,
	}
	
	return s.Decode(&b.Transactions)
}
