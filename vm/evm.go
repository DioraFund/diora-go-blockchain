package vm

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/diora-blockchain/diora/core"
	"github.com/diora-blockchain/diora/crypto"
)

type EVM struct {
	state       *core.State
	config      *core.Config
	interpreter *Interpreter
	gasTable    *GasTable
	mu          sync.RWMutex
}

type Interpreter struct {
	evm   *EVM
	table *JumpTable
}

type GasTable struct {
	Zero         uint64
	Base         uint64
	VeryLow      uint64
	Low          uint64
	Mid          uint64
	High         uint64
	ExtCode      uint64
	Balance      uint64
	SLoad        uint64
	SStore       uint64
	Create       uint64
	Call         uint64
	SelfDestruct uint64
}

type Context struct {
	Origin     core.Address
	GasPrice   *big.Int
	Coinbase   core.Address
	Number     *big.Int
	Timestamp  uint64
	Difficulty *big.Int
	GasLimit   uint64
}

type Contract struct {
	Caller   core.Address
	Address  core.Address
	Value    *big.Int
	Input    []byte
	Gas      uint64
	Code     []byte
	CodeHash []byte
}

type Stack struct {
	data []*big.Int
}

type Memory struct {
	store []byte
}

func NewEVM(state *core.State, config *core.Config) *EVM {
	return &EVM{
		state:       state,
		config:      config,
		interpreter: NewInterpreter(nil),
		gasTable:    defaultGasTable(),
	}
}

func NewInterpreter(evm *EVM) *Interpreter {
	return &Interpreter{
		evm:   evm,
		table: newJumpTable(),
	}
}

func (evm *EVM) ExecuteTransaction(tx *core.Transaction) (*core.Receipt, error) {
	evm.mu.Lock()
	defer evm.mu.Unlock()

	// Create execution context
	ctx := evm.createContext(tx)

	// Create contract for execution
	contract := &core.Contract{
		Caller:   tx.From,
		Address:  *tx.To,
		Value:    tx.Value,
		Input:    tx.Data,
		Gas:      tx.GasLimit,
		Code:     evm.state.GetCode(*tx.To),
		CodeHash: evm.state.GetCodeHash(*tx.To),
	}

	// Execute the contract
	result, err := evm.interpreter.Run(contract, ctx)
	if err != nil {
		return &core.Receipt{
			TransactionHash: tx.Hash,
			Status:          0, // Failure
			GasUsed:         tx.GasLimit - contract.Gas,
		}, err
	}

	// Create receipt
	receipt := &core.Receipt{
		TransactionHash: tx.Hash,
		Status:          1, // Success
		GasUsed:         tx.GasLimit - contract.Gas,
		Logs:            result.Logs,
	}

	// Update sender nonce
	evm.state.SetNonce(tx.From, evm.state.GetNonce(tx.From)+1)

	return receipt, nil
}

func (evm *EVM) CreateContract(caller core.Address, value *big.Int, code []byte, gas uint64) (core.Address, uint64, error) {
	evm.mu.Lock()
	defer evm.mu.Unlock()

	// Generate contract address
	nonce := evm.state.GetNonce(caller)
	contractAddr := crypto.CreateAddress(caller, nonce)

	// Create contract
	contract := &core.Contract{
		Caller:  caller,
		Address: contractAddr,
		Value:   value,
		Input:   code,
		Gas:     gas,
		Code:    code,
	}

	// Execute contract creation
	ctx := evm.createContextForContract(caller)
	result, err := evm.interpreter.Run(contract, ctx)
	if err != nil {
		return core.Address{}, 0, err
	}

	// Store contract code
	evm.state.SetCode(contractAddr, result.ReturnData)
	evm.state.SetCodeHash(contractAddr, crypto.Keccak256Hash(result.ReturnData).Bytes())

	// Update caller nonce
	evm.state.SetNonce(caller, nonce+1)

	return contractAddr, contract.Gas, nil
}

func (evm *EVM) Call(caller core.Address, addr core.Address, value *big.Int, input []byte, gas uint64) ([]byte, uint64, error) {
	evm.mu.Lock()
	defer evm.mu.Unlock()

	contract := &core.Contract{
		Caller:  caller,
		Address: addr,
		Value:   value,
		Input:   input,
		Gas:     gas,
		Code:    evm.state.GetCode(addr),
	}

	ctx := evm.createContextForContract(caller)
	result, err := evm.interpreter.Run(contract, ctx)
	if err != nil {
		return nil, 0, err
	}

	return result.ReturnData, contract.Gas, nil
}

func (evm *EVM) createContext(tx *core.Transaction) *Context {
	return &Context{
		Origin:     tx.From,
		GasPrice:   tx.GasPrice,
		Coinbase:   core.Address{}, // Will be set by block
		Number:     big.NewInt(0),  // Will be set by block
		Timestamp:  0,              // Will be set by block
		Difficulty: big.NewInt(1),
		GasLimit:   tx.GasLimit,
	}
}

func (evm *EVM) createContextForContract(caller core.Address) *Context {
	return &Context{
		Origin:     caller,
		GasPrice:   evm.config.MinGasPrice,
		Coinbase:   core.Address{},
		Number:     big.NewInt(0),
		Timestamp:  0,
		Difficulty: big.NewInt(1),
		GasLimit:   evm.config.GasLimit,
	}
}

func (in *Interpreter) Run(contract *core.Contract, ctx *Context) (*ExecutionResult, error) {
	// Initialize stack and memory
	stack := newStack()
	memory := newMemory()

	// Execute bytecode
	for {
		if len(contract.Code) == 0 {
			break
		}

		// Get next opcode
		opcode := OpCode(contract.Code[0])
		contract.Code = contract.Code[1:]

		// Check if we have enough gas
		gasCost := in.evm.gasTable.getGasCost(opcode)
		if contract.Gas < gasCost {
			return nil, fmt.Errorf("out of gas")
		}
		contract.Gas -= gasCost

		// Execute operation
		operation := in.table[opcode]
		if operation.execute == nil {
			return nil, fmt.Errorf("invalid opcode: %d", opcode)
		}

		result := operation.execute(contract, stack, memory, ctx)
		if result.err != nil {
			return nil, result.err
		}

		// Check for STOP or RETURN
		if opcode == STOP || opcode == RETURN {
			break
		}
	}

	return &ExecutionResult{
		ReturnData: memory.Data(),
		GasUsed:    contract.Gas,
		Logs:       []*core.Log{},
	}, nil
}

type OpCode byte

const (
	STOP       OpCode = 0x00
	ADD        OpCode = 0x01
	MUL        OpCode = 0x02
	SUB        OpCode = 0x03
	DIV        OpCode = 0x04
	SDIV       OpCode = 0x05
	MOD        OpCode = 0x06
	SMOD       OpCode = 0x07
	ADDMOD     OpCode = 0x08
	MULMOD     OpCode = 0x09
	EXP        OpCode = 0x0a
	SIGNEXTEND OpCode = 0x0b
	LT         OpCode = 0x10
	GT         OpCode = 0x11
	SLT        OpCode = 0x12
	SGT        OpCode = 0x13
	EQ         OpCode = 0x14
	ISZERO     OpCode = 0x15
	AND        OpCode = 0x16
	OR         OpCode = 0x17
	XOR        OpCode = 0x18
	NOT        OpCode = 0x19
	BYTE       OpCode = 0x1a
	SHL        OpCode = 0x1b
	SHR        OpCode = 0x1c
	SAR        OpCode = 0x1d
	POP        OpCode = 0x50
	MLOAD      OpCode = 0x51
	MSTORE     OpCode = 0x52
	MSTORE8    OpCode = 0x53
	SLOAD      OpCode = 0x54
	SSTORE     OpCode = 0x55
	JUMP       OpCode = 0x56
	JUMPI      OpCode = 0x57
	PC         OpCode = 0x58
	MSIZE      OpCode = 0x59
	GAS        OpCode = 0x5a
	JUMPDEST   OpCode = 0x5b
	PUSH1      OpCode = 0x60
	PUSH2      OpCode = 0x61
	PUSH32     OpCode = 0x7f
	DUP1       OpCode = 0x80
	DUP16      OpCode = 0x8f
	SWAP1      OpCode = 0x90
	SWAP16     OpCode = 0x9f
	RETURN     OpCode = 0xf3
	REVERT     OpCode = 0xfd
)

type operation struct {
	execute func(*core.Contract, *Stack, *Memory, *Context) *result
}

type result struct {
	err error
}

type ExecutionResult struct {
	ReturnData []byte
	GasUsed    uint64
	Logs       []*core.Log
}

type JumpTable [256]operation

func newJumpTable() *JumpTable {
	table := &JumpTable{}

	// Arithmetic operations
	table[ADD] = operation{execute: opAdd}
	table[MUL] = operation{execute: opMul}
	table[SUB] = operation{execute: opSub}
	table[DIV] = operation{execute: opDiv}

	// Comparison operations
	table[LT] = operation{execute: opLt}
	table[GT] = operation{execute: opGt}
	table[EQ] = operation{execute: opEq}
	table[ISZERO] = operation{execute: opIszero}

	// Bitwise operations
	table[AND] = operation{execute: opAnd}
	table[OR] = operation{execute: opOr}
	table[XOR] = operation{execute: opXor}
	table[NOT] = operation{execute: opNot}

	// Stack operations
	table[POP] = operation{execute: opPop}
	table[DUP1] = operation{execute: opDup1}
	table[SWAP1] = operation{execute: opSwap1}

	// Memory operations
	table[MLOAD] = operation{execute: opMload}
	table[MSTORE] = operation{execute: opMstore}
	table[MSTORE8] = operation{execute: opMstore8}

	// Storage operations
	table[SLOAD] = operation{execute: opSload}
	table[SSTORE] = operation{execute: opSstore}

	// Control flow
	table[JUMP] = operation{execute: opJump}
	table[JUMPI] = operation{execute: opJumpi}
	table[PC] = operation{execute: opPc}
	table[MSIZE] = operation{execute: opMsize}
	table[GAS] = operation{execute: opGas}

	// Push operations
	for i := 0; i < 32; i++ {
		table[PUSH1+OpCode(i)] = operation{execute: makePush(i + 1)}
	}

	return table
}

// Operation implementations
func opAdd(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	stack.push(new(big.Int).Add(x, y))
	return &result{}
}

func opMul(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	stack.push(new(big.Int).Mul(x, y))
	return &result{}
}

func opSub(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	stack.push(new(big.Int).Sub(x, y))
	return &result{}
}

func opDiv(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	if y.Sign() == 0 {
		stack.push(big.NewInt(0))
	} else {
		stack.push(new(big.Int).Div(x, y))
	}
	return &result{}
}

func opLt(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) < 0 {
		stack.push(big.NewInt(1))
	} else {
		stack.push(big.NewInt(0))
	}
	return &result{}
}

func opGt(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) > 0 {
		stack.push(big.NewInt(1))
	} else {
		stack.push(big.NewInt(0))
	}
	return &result{}
}

func opEq(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) == 0 {
		stack.push(big.NewInt(1))
	} else {
		stack.push(big.NewInt(0))
	}
	return &result{}
}

func opIszero(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x := stack.pop()
	if x.Sign() == 0 {
		stack.push(big.NewInt(1))
	} else {
		stack.push(big.NewInt(0))
	}
	return &result{}
}

func opAnd(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	stack.push(new(big.Int).And(x, y))
	return &result{}
}

func opOr(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	stack.push(new(big.Int).Or(x, y))
	return &result{}
}

func opXor(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	stack.push(new(big.Int).Xor(x, y))
	return &result{}
}

func opNot(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x := stack.pop()
	stack.push(new(big.Int).Not(x))
	return &result{}
}

func opPop(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	stack.pop()
	return &result{}
}

func opDup1(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x := stack.peek()
	stack.push(new(big.Int).Set(x))
	return &result{}
}

func opSwap1(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	x, y := stack.pop(), stack.pop()
	stack.push(x)
	stack.push(y)
	return &result{}
}

func opMload(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	offset := stack.pop().Uint64()
	data := memory.Get(offset, 32)
	stack.push(new(big.Int).SetBytes(data))
	return &result{}
}

func opMstore(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	offset, value := stack.pop(), stack.pop()
	memory.Set(offset.Uint64(), 32, value.Bytes())
	return &result{}
}

func opMstore8(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	offset, value := stack.pop(), stack.pop()
	memory.Set(offset.Uint64(), 1, value.Bytes()[:1])
	return &result{}
}

func opSload(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	hash := core.BytesToHash(stack.pop().Bytes())
	value := contract.EVM.State.GetState(contract.Address, hash)
	stack.push(HashToBig(value))
	return &result{}
}

func opSstore(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	hash := core.BytesToHash(stack.pop().Bytes())
	value := stack.pop()
	contract.EVM.State.SetState(contract.Address, hash, core.BytesToHash(value.Bytes()))
	return &result{}
}

func opJump(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	// Simplified jump implementation
	return &result{}
}

func opJumpi(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	// Simplified conditional jump implementation
	return &result{}
}

func opPc(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	stack.push(big.NewInt(int64(len(contract.Code))))
	return &result{}
}

func opMsize(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	stack.push(big.NewInt(int64(memory.Len())))
	return &result{}
}

func opGas(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
	stack.push(big.NewInt(int64(contract.Gas)))
	return &result{}
}

func makePush(size int) func(*core.Contract, *Stack, *Memory, *Context) *result {
	return func(contract *core.Contract, stack *Stack, memory *Memory, ctx *Context) *result {
		if len(contract.Code) < size {
			return &result{err: fmt.Errorf("insufficient data for PUSH%d", size)}
		}

		data := contract.Code[:size]
		contract.Code = contract.Code[size:]

		value := new(big.Int).SetBytes(data)
		stack.push(value)

		return &result{}
	}
}

func newStack() *Stack {
	return &Stack{data: make([]*big.Int, 0)}
}

func (s *Stack) push(d *big.Int) {
	s.data = append(s.data, d)
}

func (s *Stack) pop() *big.Int {
	if len(s.data) == 0 {
		return big.NewInt(0)
	}

	d := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return d
}

func (s *Stack) peek() *big.Int {
	if len(s.data) == 0 {
		return big.NewInt(0)
	}
	return s.data[len(s.data)-1]
}

func newMemory() *Memory {
	return &Memory{store: make([]byte, 0)}
}

func (m *Memory) Set(offset, size uint64, value []byte) {
	end := offset + size
	if end > uint64(len(m.store)) {
		m.store = append(m.store, make([]byte, end-uint64(len(m.store)))...)
	}
	copy(m.store[offset:], value)
}

func (m *Memory) Get(offset, size uint64) []byte {
	end := offset + size
	if end > uint64(len(m.store)) {
		end = uint64(len(m.store))
	}
	if offset >= end {
		return make([]byte, size)
	}

	result := make([]byte, size)
	copy(result, m.store[offset:end])
	return result
}

func (m *Memory) Len() uint64 {
	return uint64(len(m.store))
}

func (m *Memory) Data() []byte {
	return append([]byte(nil), m.store...)
}

func defaultGasTable() *GasTable {
	return &GasTable{
		Zero:         0,
		Base:         2,
		VeryLow:      3,
		Low:          5,
		Mid:          8,
		High:         10,
		ExtCode:      700,
		Balance:      700,
		SLoad:        800,
		SStore:       20000,
		Create:       32000,
		Call:         700,
		SelfDestruct: 5000,
	}
}

func (gt *GasTable) getGasCost(opcode OpCode) uint64 {
	switch opcode {
	case STOP:
		return gt.Zero
	case ADD, SUB, MUL, DIV, SDIV, MOD, SMOD, ADDMOD, MULMOD:
		return gt.VeryLow
	case LT, GT, SLT, SGT, EQ, ISZERO:
		return gt.VeryLow
	case AND, OR, XOR, NOT, BYTE:
		return gt.VeryLow
	case SHL, SHR, SAR:
		return gt.VeryLow
	case POP:
		return gt.Base
	case MLOAD, MSTORE, MSTORE8:
		return gt.VeryLow
	case SLOAD:
		return gt.SLoad
	case SSTORE:
		return gt.SStore
	case JUMP, JUMPI, PC, MSIZE, GAS:
		return gt.Base
	case PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7, PUSH8:
		return gt.VeryLow
	case PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16:
		return gt.Low
	case PUSH17, PUSH18, PUSH19, PUSH20, PUSH21, PUSH22, PUSH23, PUSH24:
		return gt.Mid
	case PUSH25, PUSH26, PUSH27, PUSH28, PUSH29, PUSH30, PUSH31, PUSH32:
		return gt.High
	case DUP1, DUP2, DUP3, DUP4, DUP5, DUP6, DUP7, DUP8:
		return gt.VeryLow
	case DUP9, DUP10, DUP11, DUP12, DUP13, DUP14, DUP15, DUP16:
		return gt.Low
	case SWAP1, SWAP2, SWAP3, SWAP4, SWAP5, SWAP6, SWAP7, SWAP8:
		return gt.VeryLow
	case SWAP9, SWAP10, SWAP11, SWAP12, SWAP13, SWAP14, SWAP15, SWAP16:
		return gt.Low
	case RETURN:
		return gt.Zero
	case REVERT:
		return gt.Zero
	default:
		return gt.Base
	}
}
