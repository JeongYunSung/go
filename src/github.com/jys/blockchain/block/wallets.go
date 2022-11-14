package block

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile)

	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets

	json.Unmarshal(fileContent, &wallets)

	ws.Wallets = wallets.Wallets

	return nil
}

func (ws Wallets) SaveToFile() {
	content, _ := json.Marshal(ws)

	f, _ := os.Create(walletFile)

	_, err := f.Write(content)
	if err != nil {
		log.Panic(err)
	}
}
