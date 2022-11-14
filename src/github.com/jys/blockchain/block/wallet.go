package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "wallet.dat"
const addressChecksumLen = 4

type PriKey struct {
	pkey []byte
	x    *big.Int
	y    *big.Int
}

type Wallet struct {
	PrivateKey []byte
	Data       *big.Int
	PublicKey  []byte
}

func NewWallet() (*Wallet, error) {
	private, public, err := newKeyPair()

	if err != nil {
		return nil, fmt.Errorf("지갑을 생성하는도중 에러가 발생했습니다. : %w\n", err)
	}

	wallet := Wallet{elliptic.MarshalCompressed(private.Curve, private.X, private.Y), private.D, public}

	return &wallet, nil
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash, err := HashPubKey(w.PublicKey)

	if err != nil {

	}

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

func (w Wallet) UnmarshalPrivateKey() ecdsa.PrivateKey {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), w.PrivateKey)
	priKey := ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}}
	priKey.D = w.Data
	return priKey
}

func newKeyPair() (ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, nil, fmt.Errorf("키를 생성하는도중 에러가 발생했습니다. : %w\n", err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey, nil
}

func HashPubKey(pubKey []byte) ([]byte, error) {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		return nil, fmt.Errorf("PHK를 생성하는도중 에러가 발생했습니다. : %w\n", err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160, nil
}

func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
