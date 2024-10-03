package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type Transaction struct {
	OwnerID    string    `json:"owner_id"`
	AccessorID string    `json:"downloader_id"`
	DocID      string    `json:"doc_id"`
	AccessTime time.Time `json:"access_time"`
}

type Block struct {
	Index        int         `json:"index"`
	Timestamp    time.Time   `json:"timestamp"`
	Data         Transaction `json:"data"`
	PreviousHash string      `json:"previous_hash"`
	Hash         string      `json:"hash"`
}

type Blockchain struct {
	Blocks []Block    `json:"blocks"`
	mutex  sync.Mutex `json:"-"`
}

func NewBlockchain() *Blockchain {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now(),
		Data:         Transaction{},
		PreviousHash: "0",
		Hash:         "",
	}
	genesisBlock.Hash = calculateHash(genesisBlock)
	return &Blockchain{
		Blocks: []Block{genesisBlock},
	}
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp.String() + block.Data.OwnerID + block.Data.AccessorID + block.Data.DocID + block.PreviousHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (bc *Blockchain) AddBlock(data Transaction) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    time.Now(),
		Data:         data,
		PreviousHash: lastBlock.Hash,
		Hash:         "",
	}
	newBlock.Hash = calculateHash(newBlock)

	if newBlock.PreviousHash != lastBlock.Hash {
		return nil
	}

	bc.Blocks = append(bc.Blocks, newBlock)

	if err := bc.SaveToFile("blockchain.json"); err != nil {
		log.Printf("Gagal menyimpan blockchain: %v", err)
	}

	return nil
}

func (bc *Blockchain) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func LoadFromFile(filename string) (*Blockchain, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return NewBlockchain(), nil
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var bc Blockchain
	err = json.Unmarshal(data, &bc)
	if err != nil {
		return nil, err
	}
	return &bc, nil
}

func (bc *Blockchain) GetHistoryByUserID(userID string) []Transaction {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	var historyList []Transaction
	for _, block := range bc.Blocks {
		if block.Data.OwnerID == userID {
			historyList = append(historyList, block.Data)
		}
	}
	return historyList
}
