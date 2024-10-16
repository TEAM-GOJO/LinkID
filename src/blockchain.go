package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"path/filepath"
	"flag"
)

type block struct {
	Index        int
	Initials     string
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
	ChainID	   int
	BlockCount int
	Genesis    block
	Head       block
	Previous   block
	Chain      []block
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
	TargetChain.BlockCount++
}

func generateBlock(previousBlock block, data []interface{}) block {
	NewBlock := block{
		Index:        data[0].(int),
		Initials:     data[1].(string),
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

// Probably not even needed actually. Just here just in case we need to mine
func mineBlock(previousBlock block, data []interface{}, difficulty int) block {
	var nonce int
	var NewBlock block

	for {
		NewBlock = block{
			Index:        data[0].(int),
			Initials:     data[1].(string),
			Age:          data[2].(int),
			Height:       data[3].(float32),
			Weight:       data[4].(float32),
			Time:         time.Now().String(),
			PreviousHash: previousBlock.CurrentHash,
			Medications:  data[5].([]string),
			Conditions:   data[6].([]string),
		}

		NewBlock.CurrentHash = calculateHash(NewBlock)

		if NewBlock.CurrentHash[:difficulty] == string(make([]byte, difficulty)) {
			break
		}
		nonce++
	}

	return NewBlock
}

func getBlockByHash(TargetChain chain, hash string) (block, bool) {
	for _, b := range TargetChain.Chain {
		if b.CurrentHash == hash {
			return b, true
		}
	}
	return block{}, false
}

func loadGenesisFromFile(filePath string) (block, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return block{}, err
	}

	var genesisData block
	err = json.Unmarshal(file, &genesisData)
	if err != nil {
		return block{}, err
	}

	return genesisData, nil
}

func generateChainID() (int, error) {
	dir := "./records"
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return 0, err
	}

	for {
		chainID := rand.Intn(90000000) + 10000000
		encFilePath := filepath.Join(dir, strconv.Itoa(chainID) + ".enc")
		if _, err := os.Stat(encFilePath); os.IsNotExist(err) {
			return chainID, nil
		}
	}
}

func generateKeyPair() (string, string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(key), hex.EncodeToString(key), nil
}

func encrypt(data []byte, keyString string) ([]byte, error) {
	key, _ := hex.DecodeString(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt(ciphertext []byte, keyString string) ([]byte, error) {
	key, _ := hex.DecodeString(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func exportEncryptedChain(TargetChain chain, key string) error {
	dir := "./records"
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dir, strconv.Itoa(TargetChain.ChainID) + ".enc")

	chainJSON, err := json.MarshalIndent(TargetChain, "", "  ")
	if err != nil {
		return err
	}

	encryptedData, err := encrypt(chainJSON, key)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, encryptedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadEncryptedChain(id string, key string) (chain, error) {
	filePath := filepath.Join("./records", id+".enc")

	encryptedData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return chain{}, err
	}

	decryptedData, err := decrypt(encryptedData, key)
	if err != nil {
		return chain{}, err
	}

	var TargetChain chain
	err = json.Unmarshal(decryptedData, &TargetChain)
	if err != nil {
		return chain{}, err
	}

	return TargetChain, nil
}

func main() {
	createCommand := flag.String("c", "", "Create a new genesis block with the path to the input JSON file.")
	accessCommand := flag.String("a", "", "Access an existing chain by ID.")
	key := flag.String("k", "", "Private key for decrypting the chain.")
	flag.Parse()

	if *createCommand != "" {
		GenesisBlock, err := loadGenesisFromFile(*createCommand)
		if err != nil {
			fmt.Println("Error loading genesis data:", err)
			return
		}

		GenesisBlock.Index = 0
		GenesisBlock.Time = time.Now().String()
		GenesisBlock.PreviousHash = ""
		GenesisBlock.CurrentHash = calculateHash(GenesisBlock)

		ChainID := generateChainID()
		TestChain := chain{
			ChainID:    ChainID,
			BlockCount: 0,
			Genesis:    GenesisBlock,
			Head:       GenesisBlock,
			Previous:   GenesisBlock,
			Chain:      []block{GenesisBlock},
		}

		publicKey, privateKey, err := generateKeyPair()
		if err != nil {
			fmt.Println("Error generating key pair:", err)
			return
		}

		fmt.Println("Public Key:", publicKey)
		fmt.Println("Private Key:", privateKey)

		err = exportEncryptedChain(TestChain, privateKey)
		if err != nil {
			fmt.Println("Error exporting encrypted chain:", err)
			return
		}

		fmt.Println("Genesis block created and saved with Chain ID:", ChainID)
	}

	if *accessCommand != "" && *key != "" {
		TargetChain, err := loadEncryptedChain(*accessCommand, *key)
		if err != nil {
			fmt.Println("Error accessing chain:", err)
			return
		}

		chainJSON, _ := json.MarshalIndent(TargetChain, "", "  ")
		fmt.Println("Decrypted chain content:", string(chainJSON))
	}
}