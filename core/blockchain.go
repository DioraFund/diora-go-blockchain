package core

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/diora-blockchain/diora/consensus"
	"github.com/diora-blockchain/diora/crypto"
	"github.com/diora-blockchain/diora/vm"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	ErrBlockNotFound      = errors.New("block not found")
	ErrInvalidBlock       = errors.New("invalid block")
	ErrInvalidChain       = errors.New("invalid chain")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Blockchain struct {
	config      *Config
	db          *leveldb.DB
	genesis     *Block
	currentBlock *Block
	state       *State
	consensus   consensus.Consensus
	vm          *vm.EVM
	
	// Caches
	blockCache  map[common.Hash]*Block
	txPool      *TxPool
	
	// Channels
	newBlockCh  chan *Block
	newTxCh     chan *Transaction
	
	// Synchronization
	mu          sync.RWMutex
	wg          sync.WaitGroup
	stopCh      chan struct{}
}

type Config struct {
	ChainID     *big.Int
	NetworkID   uint64
	BlockTime   time.Duration
	GasLimit    uint64
	MinGasPrice *big.Int
	MaxBlockSize uint64
	StakeAmount *big.Int
	ValidatorCount int
}

type Block struct {
	Header       *BlockHeader
	Transactions []*Transaction
	Hash         common.Hash
	Size         uint64
}

type BlockHeader struct {
	ParentHash  common.Hash
	Coinbase    common.Address
	StateRoot   common.Hash
	TxRoot      common.Hash
	ReceiptRoot common.Hash
	Difficulty  *big.Int
	Number      *big.Int
	GasLimit    uint64
	GasUsed     uint64
	Timestamp   uint64
	ExtraData   []byte
	MixHash     common.Hash
	Nonce       [8]byte
	Validator   common.Address
	Signature   []byte
}

type Transaction struct {
	Nonce      uint64
	GasPrice   *big.Int
	GasLimit   uint64
	To         *common.Address
	Value      *big.Int
	Data       []byte
	V, R, S    *big.Int
	Hash       common.Hash
	From       common.Address
}

type State struct {
	db     *leveldb.DB
	trie   *Trie
	cache  map[common.Address]*Account
	mu     sync.RWMutex
}

type Account struct {
	Nonce    uint64
	Balance  *big.Int
	CodeHash []byte
	Code     []byte
	Storage  map[common.Hash]common.Hash
}

type TxPool struct {
	pending   map[common.Hash]*Transaction
	queued    map[common.Hash]*Transaction
	all       map[common.Hash]*Transaction
	mu        sync.RWMWMutex
	maxSize   int
}

func NewBlockchain(config *Config, dbPath string) (*Blockchain, error) {
	// Open database
	db, err := leveldb.OpenFile(dbPath, &opt.Options{
		WriteBuffer: 64 * 1024 * 1024,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Initialize state
	state := NewState(db)
	
	// Initialize consensus engine
	consensus := consensus.NewPoS(config.StakeAmount, config.ValidatorCount)
	
	// Initialize EVM
	evm := vm.NewEVM(state, config)
	
	bc := &Blockchain{
		config:     config,
		db:         db,
		state:      state,
		consensus:  consensus,
		vm:         evm,
		blockCache: make(map[common.Hash]*Block),
		txPool:     NewTxPool(10000),
		newBlockCh: make(chan *Block, 100),
		newTxCh:    make(chan *Transaction, 1000),
		stopCh:     make(chan struct{}),
	}

	// Load or create genesis block
	if err := bc.loadGenesis(); err != nil {
		return nil, fmt.Errorf("failed to load genesis: %w", err)
	}

	// Start background workers
	bc.wg.Add(2)
	go bc.blockProcessor()
	go bc.txProcessor()

	return bc, nil
}

func (bc *Blockchain) loadGenesis() error {
	// Try to load current state
	currentHash, err := bc.db.Get([]byte("currentBlock"), nil)
	if err == nil {
		// Load existing blockchain
		hash := common.BytesToHash(currentHash)
		block, err := bc.GetBlockByHash(hash)
		if err != nil {
			return fmt.Errorf("failed to load current block: %w", err)
		}
		bc.currentBlock = block
		return nil
	}

	if err != leveldb.ErrNotFound {
		return fmt.Errorf("database error: %w", err)
	}

	// Create genesis block
	genesis := bc.createGenesisBlock()
	if err := bc.writeBlock(genesis); err != nil {
		return fmt.Errorf("failed to write genesis block: %w", err)
	}
	
	bc.genesis = genesis
	bc.currentBlock = genesis
	
	return nil
}

func (bc *Blockchain) createGenesisBlock() *Block {
	timestamp := uint64(time.Now().Unix())
	
	header := &BlockHeader{
		ParentHash:  common.Hash{},
		Coinbase:    common.HexToAddress("0x0000000000000000000000000000000000000000"),
		StateRoot:   common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
		TxRoot:      common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
		ReceiptRoot: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
		Difficulty:  big.NewInt(1),
		Number:      big.NewInt(0),
		GasLimit:    bc.config.GasLimit,
		GasUsed:     0,
		Timestamp:   timestamp,
		ExtraData:   []byte("Diora Genesis Block"),
		MixHash:     common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		Nonce:       [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}

	block := &Block{
		Header:       header,
		Transactions: []*Transaction{},
	}
	
	block.Hash = block.ComputeHash()
	block.Size = block.ComputeSize()
	
	return block
}

func (bc *Blockchain) AddTransaction(tx *Transaction) error {
	// Validate transaction
	if err := bc.ValidateTransaction(tx); err != nil {
		return err
	}

	// Add to transaction pool
	bc.txPool.Add(tx)
	
	// Notify transaction processor
	select {
	case bc.newTxCh <- tx:
	default:
		// Channel full, transaction will be processed eventually
	}
	
	return nil
}

func (bc *Blockchain) ValidateTransaction(tx *Transaction) error {
	// Check signature
	if err := crypto.VerifySignature(tx.From, tx.Hash.Bytes(), tx.V, tx.R, tx.S); err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	// Check nonce
	account := bc.state.GetAccount(tx.From)
	if tx.Nonce != account.Nonce {
		return fmt.Errorf("invalid nonce: expected %d, got %d", account.Nonce, tx.Nonce)
	}

	// Check balance
	cost := new(big.Int).Mul(tx.GasPrice, big.NewInt(int64(tx.GasLimit)))
	cost.Add(cost, tx.Value)
	
	if account.Balance.Cmp(cost) < 0 {
		return ErrInsufficientBalance
	}

	// Check gas price
	if tx.GasPrice.Cmp(bc.config.MinGasPrice) < 0 {
		return fmt.Errorf("gas price too low: minimum %s", bc.config.MinGasPrice.String())
	}

	return nil
}

func (bc *Blockchain) ProcessBlock(block *Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Validate block
	if err := bc.ValidateBlock(block); err != nil {
		return err
	}

	// Execute transactions
	receipts, err := bc.ExecuteTransactions(block.Transactions)
	if err != nil {
		return fmt.Errorf("failed to execute transactions: %w", err)
	}

	// Update state
	if err := bc.state.Commit(); err != nil {
		return fmt.Errorf("failed to commit state: %w", err)
	}

	// Write block to database
	if err := bc.writeBlock(block); err != nil {
		return fmt.Errorf("failed to write block: %w", err)
	}

	// Update current block
	bc.currentBlock = block

	// Update consensus
	bc.consensus.UpdateBlock(block)

	return nil
}

func (bc *Blockchain) ValidateBlock(block *Block) error {
	// Check block hash
	computedHash := block.ComputeHash()
	if !bytes.Equal(block.Hash.Bytes(), computedHash.Bytes()) {
		return ErrInvalidBlock
	}

	// Check parent block exists
	if block.Header.Number.Cmp(big.NewInt(0)) > 0 {
		_, err := bc.GetBlockByHash(block.Header.ParentHash)
		if err != nil {
			return fmt.Errorf("parent block not found: %w", err)
		}
	}

	// Validate consensus
	if err := bc.consensus.ValidateBlock(block); err != nil {
		return fmt.Errorf("consensus validation failed: %w", err)
	}

	// Validate transactions
	for _, tx := range block.Transactions {
		if err := bc.ValidateTransaction(tx); err != nil {
			return fmt.Errorf("invalid transaction %s: %w", tx.Hash.Hex(), err)
		}
	}

	return nil
}

func (bc *Blockchain) ExecuteTransactions(txs []*Transaction) ([]*Receipt, error) {
	var receipts []*Receipt
	
	for _, tx := range txs {
		receipt, err := bc.vm.ExecuteTransaction(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to execute transaction %s: %w", tx.Hash.Hex(), err)
		}
		receipts = append(receipts, receipt)
	}
	
	return receipts, nil
}

func (bc *Blockchain) writeBlock(block *Block) error {
	// Serialize block
	data, err := rlp.EncodeToBytes(block)
	if err != nil {
		return fmt.Errorf("failed to serialize block: %w", err)
	}

	// Write to database
	if err := bc.db.Put(block.Hash.Bytes(), data, nil); err != nil {
		return fmt.Errorf("failed to write block: %w", err)
	}

	// Update current block pointer
	if err := bc.db.Put([]byte("currentBlock"), block.Hash.Bytes(), nil); err != nil {
		return fmt.Errorf("failed to update current block: %w", err)
	}

	// Cache block
	bc.blockCache[block.Hash] = block

	return nil
}

func (bc *Blockchain) GetBlockByHash(hash common.Hash) (*Block, error) {
	// Check cache first
	if block, exists := bc.blockCache[hash]; exists {
		return block, nil
	}

	// Read from database
	data, err := bc.db.Get(hash.Bytes(), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrBlockNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Deserialize block
	var block Block
	if err := rlp.DecodeBytes(data, &block); err != nil {
		return nil, fmt.Errorf("failed to deserialize block: %w", err)
	}

	// Cache block
	bc.blockCache[hash] = &block

	return &block, nil
}

func (bc *Blockchain) GetBlockByNumber(number *big.Int) (*Block, error) {
	// For now, implement simple linear search
	// In production, use block number index
	currentBlock := bc.currentBlock
	for currentBlock != nil {
		if currentBlock.Header.Number.Cmp(number) == 0 {
			return currentBlock, nil
		}
		if currentBlock.Header.Number.Cmp(number) < 0 {
			break
		}
		
		parent, err := bc.GetBlockByHash(currentBlock.Header.ParentHash)
		if err != nil {
			return nil, err
		}
		currentBlock = parent
	}
	
	return nil, ErrBlockNotFound
}

func (bc *Blockchain) GetCurrentBlock() *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.currentBlock
}

func (bc *Blockchain) GetBalance(address common.Address) *big.Int {
	account := bc.state.GetAccount(address)
	return account.Balance
}

func (bc *Blockchain) GetNonce(address common.Address) uint64 {
	account := bc.state.GetAccount(address)
	return account.Nonce
}

func (bc *Blockchain) blockProcessor() {
	defer bc.wg.Done()
	
	ticker := time.NewTicker(bc.config.BlockTime)
	defer ticker.Stop()
	
	for {
		select {
		case <-bc.stopCh:
			return
		case <-ticker.C:
			bc.tryCreateBlock()
		case block := <-bc.newBlockCh:
			if err := bc.ProcessBlock(block); err != nil {
				fmt.Printf("Failed to process block: %v\n", err)
			}
		}
	}
}

func (bc *Blockchain) txProcessor() {
	defer bc.wg.Done()
	
	for {
		select {
		case <-bc.stopCh:
			return
		case tx := <-bc.newTxCh:
			// Transaction already validated when added to pool
			// Just ensure it's in the pool
			bc.txPool.Add(tx)
		}
	}
}

func (bc *Blockchain) tryCreateBlock() {
	if !bc.consensus.IsValidator() {
		return
	}

	// Get pending transactions
	txs := bc.txPool.Pending()
	if len(txs) == 0 {
		return
	}

	// Create new block
	block, err := bc.CreateBlock(txs)
	if err != nil {
		fmt.Printf("Failed to create block: %v\n", err)
		return
	}

	// Process block
	if err := bc.ProcessBlock(block); err != nil {
		fmt.Printf("Failed to process created block: %v\n", err)
		return
	}

	// Broadcast block
	bc.BroadcastBlock(block)
}

func (bc *Blockchain) CreateBlock(txs []*Transaction) (*Block, error) {
	parent := bc.currentBlock
	
	header := &BlockHeader{
		ParentHash:  parent.Hash,
		Coinbase:    bc.consensus.GetValidatorAddress(),
		StateRoot:   bc.state.Root(),
		TxRoot:      computeTxRoot(txs),
		ReceiptRoot: common.Hash{}, // Will be set after execution
		Difficulty:  big.NewInt(1),
		Number:      new(big.Int).Add(parent.Header.Number, big.NewInt(1)),
		GasLimit:    bc.config.GasLimit,
		GasUsed:     computeGasUsed(txs),
		Timestamp:   uint64(time.Now().Unix()),
		ExtraData:   []byte{},
		MixHash:     common.Hash{},
		Nonce:       [8]byte{},
		Validator:   bc.consensus.GetValidatorAddress(),
	}

	block := &Block{
		Header:       header,
		Transactions: txs,
	}
	
	// Sign block
	signature, err := bc.consensus.SignBlock(block)
	if err != nil {
		return nil, fmt.Errorf("failed to sign block: %w", err)
	}
	block.Header.Signature = signature
	
	block.Hash = block.ComputeHash()
	block.Size = block.ComputeSize()
	
	return block, nil
}

func (bc *Blockchain) BroadcastBlock(block *Block) {
	// Implementation for P2P block broadcasting
	// This will be implemented in the P2P module
}

func (bc *Blockchain) Stop() {
	close(bc.stopCh)
	bc.wg.Wait()
	
	if bc.db != nil {
		bc.db.Close()
	}
}

func (b *Block) ComputeHash() common.Hash {
	// Create hash of header only (transactions are included via TxRoot)
	data, _ := rlp.EncodeToBytes(b.Header)
	return crypto.Keccak256Hash(data)
}

func (b *Block) ComputeSize() uint64 {
	data, _ := rlp.EncodeToBytes(b)
	return uint64(len(data))
}

func computeTxRoot(txs []*Transaction) common.Hash {
	// Simplified implementation
	// In production, use proper Merkle tree
	if len(txs) == 0 {
		return common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	}
	
	data, _ := rlp.EncodeToBytes(txs)
	return crypto.Keccak256Hash(data)
}

func computeGasUsed(txs []*Transaction) uint64 {
	var total uint64
	for _, tx := range txs {
		total += tx.GasLimit
	}
	return total
}
