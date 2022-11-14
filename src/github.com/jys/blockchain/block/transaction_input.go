package block

import "bytes"

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash, err := HashPubKey(in.PubKey)

	if err != nil {
		panic(err)
	}

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
