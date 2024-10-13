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
	Genesis		 block
	Head         block
	Previous     block
}

func calculateHash(b block) string {
	var BlockData = []string{
		strconv.Itoa(b.Index),
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

func generateBlock(previousBlock block, data []interface{}) block {
	newBlock := block{
		Index:        data[0].(int),
		Age:          data[1].(int),
		Height:       data[2].(float32),
		Weight:       data[3].(float32),
		Time:         time.Now().String(),
		PreviousHash: previousBlock.CurrentHash,
		Medications:  data[4].([]string),
		Conditions:   data[5].([]string),
	}

	newBlock.CurrentHash = calculateHash(newBlock)
	return newBlock
}

func main() {
	genesisBlock := block{
		Index:        0,
		Age:          62,
		Height:       173.4,
		Weight:       78.2,
		Time:         time.Now().String(),
		PreviousHash: "",
		Medications:  []string{"medication1", "medication2"},
		Conditions:   []string{"destructive disease"},
	}

	genesisBlock.CurrentHash = calculateHash(genesisBlock)

	info := []interface{}{
		1, 63,
		float32(173.4), float32(78.0),
		[]string{"medication1", "medication2", "new medication"},
		[]string{"destructive disease"},
	}

	newBlock := generateBlock(genesisBlock, info)

	fmt.Printf("Genesis Block: %+v\n", genesisBlock)
	fmt.Printf("New Block: %+v\n", newBlock)
}