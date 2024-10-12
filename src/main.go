// PROOF OF CONCEPT

package main

import (
	"fmt"
	"crypto/sha256"
	"strconv"
	"time"
)

type block struct {
	Index int
	Age int
	Height float32
	Weight float32
	Time string
	PreviousHash string
	CurrentHash string
	Medications []string
	Conditions []string
}

func calculateHash(b block) string {
	var blockData = []string{
		strconv.Itoa(b.Index),
		strconv.Itoa(b.Age),
		b.Time,
		strconv.Format(b.Height, 'f', -1, 32),
		strconv.Format(b.Weight, 'f', -1, 32),
		b.PreviousHash,
		b.CurrentHash,
	}
	var record string
	for i := 0; i < len(blockData); i++ {
		record += blockData[i]
	}
	for j := 0; j < len(b.Medications); j++ {
		record += blockData[j]
	}
	for k := 0; k < len(b.Conditions); k++ {
		record += blockData[k]
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(i []interface) block {
	var newBlock block
	newBlock.Index = i[0] // int
	newBlock.Age = i[1] // int
	newBlock.Height = i[2] // float32, cm
	newBlock.Weight = i[3] // float32, kg0
	newBlock.Time = time.Now().string()// string
	newBlock.PreviousHash = i[4] //string
	newBlock.CurrentHash = i[5] // string
	newBlock.Medications = i[6] // []string
	newBlock.Conditions = i[7] // []string
	return newBlock
}

func main() {
	info := []interface{
		0,
		62,
		173.4,
		78.2,
		"", // Need to untangle messy entanglement for hash management
		"",
		[]string{"chug jug", "med mist"},
		[]string{"skill issue"},
	}
}
