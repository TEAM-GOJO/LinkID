package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type block struct {
	Index        	int
	Initials     	string
	Sex			 	string
	Gender		 	string
	Age          	int
	Height       	float32
	Weight       	float32
	BMI				float32
	Blood		 	string
	Time         	string
	Location		string
	PreviousHash 	string
	CurrentHash  	string
	Prescriptions	[]string
	Conditions   	[]string
	VisitLogs		[]string
	History			[]string
}

type chain struct {
	ChainID    int
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
		b.Sex,
		b.Gender,
		strconv.Itoa(b.Age),
		b.Time,
		strconv.FormatFloat(float64(b.Height), 'f', -1, 32),
		strconv.FormatFloat(float64(b.Weight), 'f', -1, 32),
		strconv.FormatFloat(float64(b.BMI), 'f', -1, 32),
		b.Blood,
		b.Location,
		b.PreviousHash,
	}

	var record string
	for _, data := range BlockData {
		record += data
	}

	for _, med := range b.Prescriptions {
		record += med
	}

	for _, cond := range b.Conditions {
		record += cond
	}

	for _, visit := range b.VisitLogs {
		record += visit
	}

	for _, hist := range b.History {
		record += hist
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
		Sex:          data[2].(string),
		Gender:       data[3].(string),
		Age:          data[4].(int),
		Height:       data[5].(float32),
		Weight:       data[6].(float32),
		BMI:          data[7].(float32),
		Blood:        data[8].(string),
		Time:         time.Now().String(),
		Location:     data[9].(string),
		PreviousHash: previousBlock.CurrentHash,
		Prescriptions:  data[10].([]string),
		Conditions:   data[11].([]string),
		VisitLogs:    data[12].([]string),
		History:      data[13].([]string),
	}

	NewBlock.CurrentHash = calculateHash(NewBlock)
	return NewBlock
}

func mineBlock(previousBlock block, data []interface{}, difficulty int) block {
	var nonce int
	var NewBlock block

	for {
		NewBlock = block{
			Index:        data[0].(int),
			Initials:     data[1].(string),
			Sex:          data[2].(string),
			Gender:       data[3].(string),
			Age:          data[4].(int),
			Height:       data[5].(float32),
			Weight:       data[6].(float32),
			BMI:          data[7].(float32),
			Blood:        data[8].(string),
			Time:         time.Now().String(),
			Location:     data[9].(string),
			PreviousHash: previousBlock.CurrentHash,
			Prescriptions:  data[10].([]string),
			Conditions:   data[11].([]string),
			VisitLogs:    data[12].([]string),
			History:      data[13].([]string),
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

	max := big.NewInt(90000000)
	min := 10000000

	for {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return 0, err
		}
		chainID := int(n.Int64()) + min
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

func loadBlockFromFile(filePath string) (block, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return block{}, err
	}

	var newBlock block
	err = json.Unmarshal(file, &newBlock)
	if err != nil {
		return block{}, err
	}

	return newBlock, nil
}

func main() {
	createCommand := flag.String("c", "", "Create a new chain with the path to the input JSON file.")
	accessCommand := flag.String("a", "", "Access an existing chain by ID.")
	addBlockCommand := flag.String("add", "", "Add a block to the chain with the path to the input JSON file.")
	blockChainID := flag.String("id", "", "Chain ID for adding the block or accessing the chain.")
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

		ChainID, err := generateChainID()
		if err != nil {
			fmt.Println("Error generating Chain ID:", err)
			return
		}

		TestChain := chain{
			ChainID:    ChainID,
			BlockCount: 1,
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

		keys := map[string]interface{}{
			"id":	ChainID,
			"public_key":  publicKey,
			"private_key": privateKey,
		}

		jsonData, err := json.MarshalIndent(keys, "", "  ")
		if err != nil {
			fmt.Println("Error converting keys to JSON:", err)
			return
		}

		fmt.Println(string(jsonData))

		err = exportEncryptedChain(TestChain, privateKey)
		if err != nil {
			fmt.Println("Error exporting encrypted chain:", err)
			return
		}
	} else if *accessCommand != "" && *key != "" {
		TargetChain, err := loadEncryptedChain(*accessCommand, *key)
		if err != nil {
			fmt.Println("Error loading encrypted chain:", err)
			return
		}

		jsonData, err := json.MarshalIndent(TargetChain, "", "  ")
		if err != nil {
			fmt.Println("Error converting chain to JSON:", err)
			return
		}

		fmt.Println(string(jsonData))
	} else if *addBlockCommand != "" && *blockChainID != "" {
		TargetChain, err := loadEncryptedChain(*blockChainID, *key)
		if err != nil {
			fmt.Println("Error loading encrypted chain:", err)
			return
		}

		newBlockData, err := loadBlockFromFile(*addBlockCommand)
		if err != nil {
			fmt.Println("Error loading block data:", err)
			return
		}

		addedBlock := generateBlock(TargetChain.Head, []interface{}{
			newBlockData.Index,
			newBlockData.Initials,
			newBlockData.Sex,
			newBlockData.Gender,
			newBlockData.Age,
			newBlockData.Height,
			newBlockData.Weight,
			newBlockData.BMI,
			newBlockData.Blood,
			newBlockData.Location,
			newBlockData.Prescriptions,
			newBlockData.Conditions,
			newBlockData.VisitLogs,
			newBlockData.History,
		})

		addBlockToChain(addedBlock, &TargetChain)

		err = exportEncryptedChain(TargetChain, *key)
		if err != nil {
			fmt.Println("Error exporting encrypted chain:", err)
			return
		}

		fmt.Println("Block added successfully.")
	} else {
		fmt.Println("Please provide a valid command.")
	}
}