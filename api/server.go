package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diora-blockchain/diora/core"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	blockchain *core.Blockchain
	router     *gin.Engine
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type BlockResponse struct {
	Number          uint64            `json:"number"`
	Hash            string            `json:"hash"`
	ParentHash      string            `json:"parentHash"`
	Timestamp       uint64            `json:"timestamp"`
	Transactions    []TransactionResp `json:"transactions"`
	GasUsed         uint64            `json:"gasUsed"`
	GasLimit        uint64            `json:"gasLimit"`
	Miner           string            `json:"miner"`
	Difficulty      string            `json:"difficulty"`
	TotalDifficulty string            `json:"totalDifficulty"`
	Size            uint64            `json:"size"`
}

type TransactionResp struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to,omitempty"`
	Value       string `json:"value"`
	GasPrice    string `json:"gasPrice"`
	GasLimit    uint64 `json:"gasLimit"`
	GasUsed     uint64 `json:"gasUsed"`
	Nonce       uint64 `json:"nonce"`
	Input       string `json:"input,omitempty"`
	Status      uint64 `json:"status"`
	BlockNumber uint64 `json:"blockNumber,omitempty"`
	BlockHash   string `json:"blockHash,omitempty"`
	Timestamp   uint64 `json:"timestamp,omitempty"`
}

type AccountResp struct {
	Address     string `json:"address"`
	Balance     string `json:"balance"`
	Nonce       uint64 `json:"nonce"`
	Code        string `json:"code,omitempty"`
	StorageRoot string `json:"storageRoot"`
}

type NetworkStats struct {
	ChainID           string  `json:"chainId"`
	NetworkID         uint64  `json:"networkId"`
	BlockNumber       uint64  `json:"blockNumber"`
	BlockTime         uint64  `json:"blockTime"`
	GasPrice          string  `json:"gasPrice"`
	TotalTransactions uint64  `json:"totalTransactions"`
	ActiveValidators  uint64  `json:"activeValidators"`
	TPS               float64 `json:"tps"`
	Difficulty        string  `json:"difficulty"`
	HashRate          string  `json:"hashRate"`
}

func NewServer(blockchain *core.Blockchain) *Server {
	return &Server{
		blockchain: blockchain,
		router:     gin.Default(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (s *Server) setupRoutes() {
	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	s.router.Use(cors.New(config))

	// API routes
	v1 := s.router.Group("/api/v1")
	{
		// Blockchain info
		v1.GET("/network/stats", s.getNetworkStats)
		v1.GET("/block/latest", s.getLatestBlock)
		v1.GET("/block/:number", s.getBlockByNumber)
		v1.GET("/block/hash/:hash", s.getBlockByHash)
		v1.GET("/blocks", s.getBlocks)

		// Transactions
		v1.GET("/transaction/:hash", s.getTransaction)
		v1.GET("/transactions", s.getTransactions)
		v1.GET("/transactions/pending", s.getPendingTransactions)
		v1.POST("/transaction", s.sendTransaction)

		// Accounts
		v1.GET("/account/:address", s.getAccount)
		v1.GET("/account/:address/balance", s.getBalance)
		v1.GET("/account/:address/nonce", s.getNonce)
		v1.GET("/account/:address/transactions", s.getAccountTransactions)

		// Contracts
		v1.POST("/contract/call", s.callContract)
		v1.POST("/contract/deploy", s.deployContract)
		v1.GET("/contract/:address/code", s.getContractCode)
		v1.GET("/contract/:address/storage", s.getContractStorage)

		// Tokens
		v1.GET("/tokens", s.getTokens)
		v1.GET("/token/:address", s.getTokenInfo)
		v1.GET("/token/:address/balance/:account", s.getTokenBalance)
		v1.GET("/token/:address/transfers", s.getTokenTransfers)

		// NFTs
		v1.GET("/nfts", s.getNFTs)
		v1.GET("/nft/:contract/:tokenId", s.getNFT)
		v1.GET("/nft/:contract/owner/:tokenId", s.getNFTOwner)
		v1.GET("/nft/:contract/transfers", s.getNFTTransfers)

		// Staking
		v1.GET("/validators", s.getValidators)
		v1.GET("/validator/:address", s.getValidator)
		v1.GET("/staking/info/:address", s.getStakingInfo)
		v1.POST("/staking/delegate", s.delegate)
		v1.POST("/staking/undelegate", s.undelegate)
		v1.POST("/staking/claim", s.claimRewards)

		// Governance
		v1.GET("/proposals", s.getProposals)
		v1.GET("/proposal/:id", s.getProposal)
		v1.POST("/proposal", s.createProposal)
		v1.POST("/proposal/:id/vote", s.voteProposal)

		// Events
		v1.GET("/events", s.getEvents)
		v1.GET("/events/logs", s.getLogs)
	}

	// WebSocket endpoint
	s.router.GET("/ws", s.handleWebSocket)

	// Health check
	s.router.GET("/health", s.healthCheck)
}

func (s *Server) getNetworkStats(c *gin.Context) {
	currentBlock := s.blockchain.GetCurrentBlock()
	validators := s.blockchain.GetValidators()

	stats := NetworkStats{
		ChainID:           "1337", // Diora testnet
		NetworkID:         1,
		BlockNumber:       currentBlock.Header.Number.Uint64(),
		BlockTime:         6,            // 6 seconds
		GasPrice:          "1000000000", // 1 Gwei
		TotalTransactions: 0,            // TODO: Implement transaction count
		ActiveValidators:  uint64(len(validators)),
		TPS:               100.5, // TODO: Calculate real TPS
		Difficulty:        currentBlock.Header.Difficulty.String(),
		HashRate:          "1.2 TH/s", // TODO: Calculate real hash rate
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

func (s *Server) getLatestBlock(c *gin.Context) {
	block := s.blockchain.GetCurrentBlock()
	response := s.formatBlock(block)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

func (s *Server) getBlockByNumber(c *gin.Context) {
	numberStr := c.Param("number")
	var number uint64
	_, err := fmt.Sscanf(numberStr, "%d", &number)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid block number",
		})
		return
	}

	block, err := s.blockchain.GetBlockByNumber(new(big.Int).SetUint64(number))
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Block not found",
		})
		return
	}

	response := s.formatBlock(block)
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

func (s *Server) getBlockByHash(c *gin.Context) {
	hashStr := c.Param("hash")
	hash := core.HexToHash(hashStr)

	block, err := s.blockchain.GetBlockByHash(hash)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Block not found",
		})
		return
	}

	response := s.formatBlock(block)
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

func (s *Server) getBlocks(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// TODO: Implement pagination
	// For now, return latest blocks
	var blocks []BlockResponse
	currentBlock := s.blockchain.GetCurrentBlock()

	for i := 0; i < 10; i++ {
		if currentBlock == nil {
			break
		}

		blocks = append(blocks, s.formatBlock(currentBlock))
		if currentBlock.Header.Number.Uint64() == 0 {
			break
		}

		parent, err := s.blockchain.GetBlockByHash(currentBlock.Header.ParentHash)
		if err != nil {
			break
		}
		currentBlock = parent
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"blocks": blocks,
			"page":   page,
			"limit":  limit,
		},
	})
}

func (s *Server) getTransaction(c *gin.Context) {
	hashStr := c.Param("hash")
	hash := core.HexToHash(hashStr)

	// TODO: Implement transaction lookup
	// For now, return not found
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Error:   "Transaction not found",
	})
}

func (s *Server) getTransactions(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// TODO: Implement transaction listing
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"transactions": []TransactionResp{},
			"page":         page,
			"limit":        limit,
		},
	})
}

func (s *Server) getPendingTransactions(c *gin.Context) {
	// TODO: Implement pending transactions
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    []TransactionResp{},
	})
}

func (s *Server) sendTransaction(c *gin.Context) {
	var txData struct {
		From      string `json:"from"`
		To        string `json:"to,omitempty"`
		Value     string `json:"value"`
		GasPrice  string `json:"gasPrice"`
		GasLimit  uint64 `json:"gasLimit"`
		Data      string `json:"data,omitempty"`
		Nonce     uint64 `json:"nonce"`
		Signature string `json:"signature"`
	}

	if err := c.ShouldBindJSON(&txData); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid transaction data",
		})
		return
	}

	// TODO: Implement transaction sending
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"txHash": "0x1234567890abcdef", // Mock hash
		},
	})
}

func (s *Server) getAccount(c *gin.Context) {
	addressStr := c.Param("address")
	address := core.HexToAddress(addressStr)

	balance := s.blockchain.GetBalance(address)
	nonce := s.blockchain.GetNonce(address)

	response := AccountResp{
		Address: addressStr,
		Balance: balance.String(),
		Nonce:   nonce,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

func (s *Server) getBalance(c *gin.Context) {
	addressStr := c.Param("address")
	address := core.HexToAddress(addressStr)

	balance := s.blockchain.GetBalance(address)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"balance": balance.String(),
		},
	})
}

func (s *Server) getNonce(c *gin.Context) {
	addressStr := c.Param("address")
	address := core.HexToAddress(addressStr)

	nonce := s.blockchain.GetNonce(address)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"nonce": nonce,
		},
	})
}

func (s *Server) getAccountTransactions(c *gin.Context) {
	// TODO: Implement account transaction history
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    []TransactionResp{},
	})
}

// WebSocket handler for real-time updates
func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	s.clients[conn] = true

	// Send initial data
	s.sendInitialData(conn)

	// Handle messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			delete(s.clients, conn)
			break
		}

		if messageType == websocket.TextMessage {
			log.Printf("Received message: %s", message)
		}
	}
}

func (s *Server) sendInitialData(conn *websocket.Conn) {
	// Send latest block
	currentBlock := s.blockchain.GetCurrentBlock()
	if currentBlock != nil {
		blockData := s.formatBlock(currentBlock)
		message, _ := json.Marshal(map[string]interface{}{
			"type":  "newBlock",
			"block": blockData,
		})
		conn.WriteMessage(websocket.TextMessage, message)
	}

	// Send network stats
	stats := NetworkStats{
		ChainID:          "1337",
		NetworkID:        1,
		BlockNumber:      currentBlock.Header.Number.Uint64(),
		ActiveValidators: 100, // TODO: Get real count
		TPS:              100.5,
	}

	statsMessage, _ := json.Marshal(map[string]interface{}{
		"type":  "networkStats",
		"stats": stats,
	})
	conn.WriteMessage(websocket.TextMessage, statsMessage)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		},
	})
}

func (s *Server) formatBlock(block *core.Block) BlockResponse {
	var txs []TransactionResp
	for _, tx := range block.Transactions {
		txs = append(txs, TransactionResp{
			Hash:        tx.Hash.Hex(),
			From:        tx.From.Hex(),
			Value:       tx.Value.String(),
			GasPrice:    tx.GasPrice.String(),
			GasLimit:    tx.GasLimit,
			GasUsed:     tx.GasLimit, // TODO: Get actual gas used
			Nonce:       tx.Nonce,
			Status:      1, // Success
			BlockNumber: block.Header.Number.Uint64(),
			BlockHash:   block.Hash.Hex(),
			Timestamp:   block.Header.Timestamp,
		})
	}

	return BlockResponse{
		Number:          block.Header.Number.Uint64(),
		Hash:            block.Hash.Hex(),
		ParentHash:      block.Header.ParentHash.Hex(),
		Timestamp:       block.Header.Timestamp,
		Transactions:    txs,
		GasUsed:         block.Header.GasUsed,
		GasLimit:        block.Header.GasLimit,
		Miner:           block.Header.Validator.Hex(),
		Difficulty:      block.Header.Difficulty.String(),
		TotalDifficulty: block.Header.Difficulty.String(), // TODO: Calculate total difficulty
		Size:            block.Size,
	}
}

// Placeholder implementations for remaining endpoints
func (s *Server) callContract(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) deployContract(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getContractCode(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getContractStorage(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getTokens(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getTokenInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getTokenBalance(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getTokenTransfers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getNFTs(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getNFT(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getNFTOwner(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getNFTTransfers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getValidators(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getValidator(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getStakingInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) delegate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) undelegate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) claimRewards(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getProposals(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getProposal(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) createProposal(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) voteProposal(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getEvents(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getLogs(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}
func (s *Server) getAccountTransactions(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, APIResponse{Success: false, Error: "Not implemented"})
}

func (s *Server) start() {
	s.setupRoutes()

	// Start broadcast goroutine
	go s.handleBroadcasts()

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		// TODO: Implement graceful shutdown
		os.Exit(0)
	}()

	log.Println("Diora API Server starting on :8080")
	if err := s.router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func (s *Server) handleBroadcasts() {
	for {
		message := <-s.broadcast
		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}

func main() {
	// Initialize blockchain
	config := &core.Config{
		ChainID:        big.NewInt(1337),
		NetworkID:      1,
		BlockTime:      6 * time.Second,
		GasLimit:       30000000,
		MinGasPrice:    big.NewInt(1000000000), // 1 Gwei
		MaxBlockSize:   1048576,                // 1MB
		StakeAmount:    big.NewInt(1000000),    // 1 DIO
		ValidatorCount: 100,
	}

	blockchain, err := core.NewBlockchain(config, "./data")
	if err != nil {
		log.Fatal("Failed to initialize blockchain:", err)
	}

	server := NewServer(blockchain)
	server.start()
}
