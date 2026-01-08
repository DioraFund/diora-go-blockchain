package tests

import (
	"math/big"
	"testing"
	"time"

	"github.com/diora-blockchain/diora/core"
	"github.com/diora-blockchain/diora/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockchainCreation(t *testing.T) {
	config := &core.Config{
		ChainID:        big.NewInt(1337),
		NetworkID:      1,
		BlockTime:      6 * time.Second,
		GasLimit:       30000000,
		MinGasPrice:    big.NewInt(1000000000),
		MaxBlockSize:   1048576,
		StakeAmount:    big.NewInt(1000000),
		ValidatorCount: 100,
	}

	// Create temporary database
	db, err := core.NewLevelDBDatabase("/tmp/testdb")
	require.NoError(t, err)
	defer db.Close()

	blockchain, err := core.NewBlockchain(config, "/tmp/testdb")
	require.NoError(t, err)
	require.NotNil(t, blockchain)

	// Check genesis block
	genesis := blockchain.GetCurrentBlock()
	assert.NotNil(t, genesis)
	assert.Equal(t, uint64(0), genesis.Header.Number.Uint64())
}

func TestTransactionValidation(t *testing.T) {
	config := &core.Config{
		ChainID:        big.NewInt(1337),
		NetworkID:      1,
		BlockTime:      6 * time.Second,
		GasLimit:       30000000,
		MinGasPrice:    big.NewInt(1000000000),
		MaxBlockSize:   1048576,
		StakeAmount:    big.NewInt(1000000),
		ValidatorCount: 100,
	}

	db, err := core.NewLevelDBDatabase("/tmp/testdb")
	require.NoError(t, err)
	defer db.Close()

	blockchain, err := core.NewBlockchain(config, "/tmp/testdb")
	require.NoError(t, err)

	// Create test transaction
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	from := crypto.PubkeyToAddress(&privateKey.PublicKey)
	to := core.HexToAddress("0x1234567890123456789012345678901234567890")
	value := big.NewInt(1000000000000000000000) // 1 DIO

	tx := &core.Transaction{
		From:     from,
		To:       &to,
		Value:    value,
		GasPrice: big.NewInt(1000000000), // 1 Gwei
		GasLimit: 21000,
		Nonce:    0,
		Data:     []byte{},
	}

	// Sign transaction
	tx.Hash = tx.ComputeHash()
	tx.V, tx.R, tx.S = crypto.Sign(tx.Hash.Bytes(), privateKey)

	// Validate transaction
	err = blockchain.ValidateTransaction(tx)
	assert.NoError(t, err)
}

func TestBlockValidation(t *testing.T) {
	config := &core.Config{
		ChainID:        big.NewInt(1337),
		NetworkID:      1,
		BlockTime:      6 * time.Second,
		GasLimit:       30000000,
		MinGasPrice:    big.NewInt(1000000000),
		MaxBlockSize:   1048576,
		StakeAmount:    big.NewInt(1000000),
		ValidatorCount: 100,
	}

	db, err := core.NewLevelDBDatabase("/tmp/testdb")
	require.NoError(t, err)
	defer db.Close()

	blockchain, err := core.NewBlockchain(config, "/tmp/testdb")
	require.NoError(t, err)

	// Create test block
	block := &core.Block{
		Header: &core.BlockHeader{
			ParentHash:  core.Hash{},
			Coinbase:    core.HexToAddress("0x1234567890123456789012345678901234567890"),
			StateRoot:   core.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
			TxRoot:      core.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
			ReceiptRoot: core.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
			Difficulty:  big.NewInt(1),
			Number:      big.NewInt(1),
			GasLimit:    30000000,
			GasUsed:     21000,
			Timestamp:   uint64(time.Now().Unix()),
			ExtraData:   []byte("test block"),
			MixHash:     core.Hash{},
			Nonce:       [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		Transactions: []*core.Transaction{},
	}

	block.Hash = block.ComputeHash()
	block.Size = block.ComputeSize()

	// Validate block
	err = blockchain.ValidateBlock(block)
	assert.NoError(t, err)
}

func TestAccountBalance(t *testing.T) {
	config := &core.Config{
		ChainID:        big.NewInt(1337),
		NetworkID:      1,
		BlockTime:      6 * time.Second,
		GasLimit:       30000000,
		MinGasPrice:    big.NewInt(1000000000),
		MaxBlockSize:   1048576,
		StakeAmount:    big.NewInt(1000000),
		ValidatorCount: 100,
	}

	db, err := core.NewLevelDBDatabase("/tmp/testdb")
	require.NoError(t, err)
	defer db.Close()

	blockchain, err := core.NewBlockchain(config, "/tmp/testdb")
	require.NoError(t, err)

	// Test address
	address := core.HexToAddress("0x1234567890123456789012345678901234567890")

	// Check initial balance
	balance := blockchain.GetBalance(address)
	assert.Equal(t, big.NewInt(0), balance)

	// Check initial nonce
	nonce := blockchain.GetNonce(address)
	assert.Equal(t, uint64(0), nonce)
}

func TestStateManagement(t *testing.T) {
	config := &core.Config{
		ChainID:        big.NewInt(1337),
		NetworkID:      1,
		BlockTime:      6 * time.Second,
		GasLimit:       30000000,
		MinGasPrice:    big.NewInt(1000000000),
		MaxBlockSize:   1048576,
		StakeAmount:    big.NewInt(1000000),
		ValidatorCount: 100,
	}

	db, err := core.NewLevelDBDatabase("/tmp/testdb")
	require.NoError(t, err)
	defer db.Close()

	blockchain, err := core.NewBlockchain(config, "/tmp/testdb")
	require.NoError(t, err)

	// Get state
	state := blockchain.GetState()
	assert.NotNil(t, state)

	// Test account creation
	address := core.HexToAddress("0x1234567890123456789012345678901234567890")
	balance := big.NewInt(1000000000000000000000) // 1000 DIO

	blockchain.SetBalance(address, balance)

	// Verify balance
	retrievedBalance := blockchain.GetBalance(address)
	assert.Equal(t, balance, retrievedBalance)
}

func BenchmarkTransactionValidation(b *testing.B) {
	config := &core.Config{
		ChainID:        big.NewInt(1337),
		NetworkID:      1,
		BlockTime:      6 * time.Second,
		GasLimit:       30000000,
		MinGasPrice:    big.NewInt(1000000000),
		MaxBlockSize:   1048576,
		StakeAmount:    big.NewInt(1000000),
		ValidatorCount: 100,
	}

	db, err := core.NewLevelDBDatabase("/tmp/testdb")
	require.NoError(b, err)
	defer db.Close()

	blockchain, err := core.NewBlockchain(config, "/tmp/testdb")
	require.NoError(b, err)

	privateKey, err := crypto.GenerateKey()
	require.NoError(b, err)

	from := crypto.PubkeyToAddress(&privateKey.PublicKey)
	to := core.HexToAddress("0x1234567890123456789012345678901234567890")
	value := big.NewInt(1000000000000000000000)

	tx := &core.Transaction{
		From:     from,
		To:       &to,
		Value:    value,
		GasPrice: big.NewInt(1000000000),
		GasLimit: 21000,
		Nonce:    0,
		Data:     []byte{},
	}

	tx.Hash = tx.ComputeHash()
	tx.V, tx.R, tx.S = crypto.Sign(tx.Hash.Bytes(), privateKey)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := blockchain.ValidateTransaction(tx)
		if err != nil {
			b.Fatalf("Transaction validation failed: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkBlockValidation(b *testing.B) {
	config := &core.Config{
		ChainID:        big.NewInt(1337),
		NetworkID:      1,
		BlockTime:      6 * time.Second,
		GasLimit:       30000000,
		MinGasPrice:    big.NewInt(1000000000),
		MaxBlockSize:   1048576,
		StakeAmount:    big.NewInt(1000000),
		ValidatorCount: 100,
	}

	db, err := core.NewLevelDBDatabase("/tmp/testdb")
	require.NoError(b, err)
	defer db.Close()

	blockchain, err := core.NewBlockchain(config, "/tmp/testdb")
	require.NoError(b, err)

	block := &core.Block{
		Header: &core.BlockHeader{
			ParentHash:  core.Hash{},
			Coinbase:    core.HexToAddress("0x1234567890123456789012345678901234567890"),
			StateRoot:   core.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
			TxRoot:      core.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
			ReceiptRoot: core.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
			Difficulty:  big.NewInt(1),
			Number:      big.NewInt(1),
			GasLimit:    30000000,
			GasUsed:     21000,
			Timestamp:   uint64(time.Now().Unix()),
			ExtraData:   []byte("test block"),
			MixHash:     core.Hash{},
			Nonce:       [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		Transactions: []*core.Transaction{},
	}

	block.Hash = block.ComputeHash()
	block.Size = block.ComputeSize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := blockchain.ValidateBlock(block)
		if err != nil {
			b.Fatalf("Block validation failed: %v", err)
		}
	}
	b.StopTimer()
}
