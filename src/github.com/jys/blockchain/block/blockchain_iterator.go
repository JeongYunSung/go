package block

import "github.com/boltdb/bolt"

type Iterator struct {
	currentHash []byte
	db          *bolt.DB
	size        int
	point       int
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

func (i *Iterator) Next() *Block {
	var b *Block

	i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(i.currentHash)
		b = DeserializeBlock(encodedBlock)
		return nil
	})

	i.currentHash = b.PrevBlockHash
	i.point++
	return b
}

func (i *Iterator) HasNext() bool {
	return i.point < i.size
}
