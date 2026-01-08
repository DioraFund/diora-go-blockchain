package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

type Block struct {
	Index        int       `json:"index"`
	Timestamp    time.Time `json:"timestamp"`
	Data         string    `json:"data"`
	PreviousHash string    `json:"previous_hash"`
	Hash         string    `json:"hash"`
}

type Blockchain struct {
	Blocks []Block `json:"blocks"`
}

var blockchain *Blockchain
var transactions []Transaction

func (b *Block) CalculateHash() string {
	return fmt.Sprintf("%d%s%s%s", b.Index, b.Timestamp.String(), b.Data, b.PreviousHash)
}

func NewBlock(index int, previousHash string, data string) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now(),
		Data:         data,
		PreviousHash: previousHash,
	}
	block.Hash = block.CalculateHash()
	return block
}

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock(0, "", "Genesis Block")
	return &Blockchain{
		Blocks: []Block{*genesisBlock},
	}
}

func (bc *Blockchain) AddBlock(data string) {
	previousBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(previousBlock.Index+1, previousBlock.Hash, data)
	bc.Blocks = append(bc.Blocks, *newBlock)
}

func main() {
	blockchain = NewBlockchain()

	// CLI
	var rootCmd = &cobra.Command{
		Use:   "diora",
		Short: "Diora Blockchain CLI",
		Long:  `Diora is a modern EVM-compatible blockchain with PoS consensus`,
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the blockchain node",
		Run:   startNode,
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show blockchain status",
		Run:   showStatus,
	}

	rootCmd.AddCommand(startCmd, statusCmd)
	rootCmd.Execute()
}

func startNode(cmd *cobra.Command, args []string) {
	// Start HTTP server
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	r.GET("/api/blocks", getBlocks)
	r.GET("/api/blocks/:index", getBlock)
	r.POST("/api/blocks", addBlock)
	r.GET("/api/status", getStatus)
	r.GET("/api/transactions", getTransactions)
	r.POST("/api/transactions", createTransaction)

	// Wallet endpoints
	r.GET("/api/wallets", getWallets)
	r.POST("/api/wallets", createWallet)
	r.GET("/api/wallets/:address", getWallet)
	r.POST("/api/wallets/:address/send", sendFromWallet)
	r.GET("/api/wallets/:address/balance", getWalletBalance)

	// Web interface
	r.StaticFile("/", "./web/dist/index.html")
	r.Static("/static", "./web/dist")

	fmt.Println("🚀 Diora Blockchain starting on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}

func getBlocks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"blocks": blockchain.Blocks,
		"count":  len(blockchain.Blocks),
	})
}

func getBlock(c *gin.Context) {
	indexStr := c.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(blockchain.Blocks) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
		return
	}

	c.JSON(http.StatusOK, blockchain.Blocks[index])
}

func addBlock(c *gin.Context) {
	var request struct {
		Data string `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blockchain.AddBlock(request.Data)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Block added successfully",
		"block":   blockchain.Blocks[len(blockchain.Blocks)-1],
	})
}

func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":     "running",
		"blocks":     len(blockchain.Blocks),
		"last_block": blockchain.Blocks[len(blockchain.Blocks)-1],
	})
}

func showStatus(cmd *cobra.Command, args []string) {
	fmt.Printf("📊 Diora Blockchain Status\n")
	fmt.Printf("========================\n")
	fmt.Printf("Status: Running\n")
	fmt.Printf("Blocks: %d\n", len(blockchain.Blocks))
	fmt.Printf("Transactions: %d\n", len(transactions))
	fmt.Printf("Last Block: #%d\n", blockchain.Blocks[len(blockchain.Blocks)-1].Index)
	fmt.Printf("Network: Diora Testnet\n")
	fmt.Printf("Node ID: %s\n", "diora-node-1")
}

func getTransactions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
	})
}

func createTransaction(c *gin.Context) {
	var request TransactionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new transaction
	tx, err := NewTransaction(request.From, request.To, request.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate transaction
	if !tx.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction"})
		return
	}

	// Add to transactions pool
	transactions = append(transactions, *tx)

	// Create a block with transaction
	blockData := fmt.Sprintf("Transaction: %s -> %s (%s)", tx.From, tx.To, tx.Amount.String())
	blockchain.AddBlock(blockData)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transaction created successfully",
		"transaction": tx.ToJSON(),
		"block_index": len(blockchain.Blocks) - 1,
	})
}

// Wallet functions
func getWallets(c *gin.Context) {
	wallets := GetAllWallets()
	walletData := make([]map[string]interface{}, len(wallets))

	for i, wallet := range wallets {
		walletData[i] = wallet.ToJSON()
		// Remove private key for security
		delete(walletData[i], "private_key")
	}

	c.JSON(http.StatusOK, gin.H{
		"wallets": walletData,
		"count":   len(wallets),
	})
}

func createWallet(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := NewWallet(request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Wallet created successfully",
		"wallet":  wallet.ToJSON(),
	})
}

func getWallet(c *gin.Context) {
	address := c.Param("address")

	wallet, exists := GetWallet(address)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	walletData := wallet.ToJSON()
	// Remove private key for security
	delete(walletData, "private_key")

	c.JSON(http.StatusOK, gin.H{
		"wallet": walletData,
	})
}

func sendFromWallet(c *gin.Context) {
	address := c.Param("address")

	var request struct {
		To     string `json:"to" binding:"required"`
		Amount string `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, exists := GetWallet(address)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Convert amount to big.Int
	amount, ok := new(big.Int).SetString(request.Amount, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	err := wallet.Send(request.To, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction sent successfully",
		"from":    wallet.Address,
		"to":      request.To,
		"amount":  amount.String(),
	})
}

func getWalletBalance(c *gin.Context) {
	address := c.Param("address")

	wallet, exists := GetWallet(address)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": wallet.Address,
		"balance": wallet.Balance.String(),
	})
}
