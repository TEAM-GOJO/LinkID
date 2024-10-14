package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type block struct {
	Index        int
	Initials	 string
	Age          int
	Height       float32
	Weight       float32
	Time         string
	PreviousHash string
	CurrentHash  string
	Medications  []string
	Conditions   []string
}

type chain struct {
	BlockCount   int
	Genesis      block
	Head         block
	Previous     block
	Chain        []block
}

func calculateHash(b block) string {
	var BlockData = []string{
		strconv.Itoa(b.Index),
		b.Initials,
		strconv.Itoa(b.Age),
		b.Time,
		strconv.FormatFloat(float64(b.Height), 'f', -1, 32),
		strconv.FormatFloat(float64(b.Weight), 'f', -1, 32),
		b.PreviousHash,
	}

	var record string
	for _, data := range BlockData {
		record += data
	}

	for _, med := range b.Medications {
		record += med
	}

	for _, cond := range b.Conditions {
		record += cond
	}

	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func addBlockToChain(AddedBlock block, TargetChain *chain) {
	TargetChain.Chain = append(TargetChain.Chain, AddedBlock)
	TargetChain.Previous = TargetChain.Head
	TargetChain.Head = AddedBlock
	TargetChain.BlockCount = TargetChain.BlockCount + 1
}

func generateBlock(previousBlock block, data []interface{}) block {
	NewBlock := block{
		Index:        data[0].(int),
		Initials:	  data[1].(string),
		Age:          data[2].(int),
		Height:       data[3].(float32),
		Weight:       data[4].(float32),
		Time:         time.Now().String(),
		PreviousHash: previousBlock.CurrentHash,
		Medications:  data[5].([]string),
		Conditions:   data[6].([]string),
	}

	NewBlock.CurrentHash = calculateHash(NewBlock)
	return NewBlock
}

func main() {
	GenesisBlock := block{
		Index:        0,
		Initials:	  "SP",
		Age:          62,
		Height:       173.4,
		Weight:       78.2,
		Time:         time.Now().String(),
		PreviousHash: "",
		Medications:  []string{"medication1", "medication2"},
		Conditions:   []string{"destructive disease"},
	}

	GenesisBlock.CurrentHash = calculateHash(GenesisBlock)

	TestChain := chain{
		BlockCount: 0,
		Genesis:    GenesisBlock,
		Head:       GenesisBlock,
		Previous:   GenesisBlock,
		Chain:      []block{GenesisBlock}, // Changed to []block
	}

	info := []interface{}{
		1, "SP", 63,
		float32(173.4), float32(78.0),
		[]string{"medication1", "medication2", "new medication"},
		[]string{"destructive disease"},
	}

	NewBlock := generateBlock(GenesisBlock, info)

	fmt.Printf("Genesis Block: %+v\n\n", GenesisBlock)
	fmt.Printf("Head Block (before): %+v\n\n", TestChain.Head)
	fmt.Printf("New Block: %+v\n\n", NewBlock)

	addBlockToChain(NewBlock, &TestChain) // Pass the chain by reference using &

	fmt.Printf("Head Block (after): %+v\n", TestChain.Head)
	fmt.Printf("Block Count: %d\n", TestChain.BlockCount)
}
