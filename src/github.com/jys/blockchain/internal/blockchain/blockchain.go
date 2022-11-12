package blockchain

import (
	"block"
	"github.com/boltdb/bolt"
	"log"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type Iterator struct {
	currentHash []byte
	db          *bolt.DB
	size        int
	point       int
}

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	newBlock := block.NewBlock(data, lastHash)

	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())

		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash
		return nil
	})
}

func (bc *Blockchain) Iterator() *Iterator {
	lastIndex := 0

	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastIndex = b.Stats().KeyN
		return nil
	})

	bci := &Iterator{bc.tip, bc.db, lastIndex, 1}
	return bci
}

func (i *Iterator) Next() *block.Block {
	var b *block.Block

	i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(i.currentHash)
		b = block.DeserializeBlock(encodedBlock)
		return nil
	})

	i.currentHash = b.PrevBlockHash
	i.point++
	return b
}

func (i *Iterator) HasNext() bool {
	return i.point < i.size
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := block.NewGenesisBlock()
			b, _ := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	bc := Blockchain{tip, db}
	return &bc
}
