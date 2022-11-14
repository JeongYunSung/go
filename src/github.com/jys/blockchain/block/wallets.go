package block

import (
	"encoding/json"
	"fmt"
	"os"
)

type Wallets struct {
	Wallets map[string]*Wallet
	err     error
}

func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	wallets.loadFromFile()

	if wallets.err != nil {
		wallets.err = fmt.Errorf("지갑목록을 생성하는 도중 에러가 발생했습니다. : %w\n", wallets.err)
	}

	return &wallets, wallets.err
}

func (ws *Wallets) CreateWallet() (string, error) {
	wallet, err := NewWallet()

	if err != nil {
		return "", fmt.Errorf("지갑을 생성하는 도중 에러가 발생했습니다. : %w\n", err)
	}

	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address, nil
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

func (ws *Wallets) loadFromFile() {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		ws.err = fmt.Errorf("지갑 파일이 존재하지 않습니다. : %w\n", err)
		return
	}

	fileContent, err := os.ReadFile(walletFile)

	if err != nil {
		ws.err = fmt.Errorf("지갑 파일을 읽는 도중 에러가 발생했습니다. : %w\n", err)
		return
	}

	var wallets Wallets

	if err := json.Unmarshal(fileContent, &wallets); err != nil {
		ws.err = fmt.Errorf("지갑 파일을 언마셜 하는 도중 에러가 발생했습니다. : %w\n", err)
		return
	}

	ws.Wallets = wallets.Wallets
}

func (ws Wallets) SaveToFile() error {
	content, err := json.Marshal(ws)

	if err != nil {
		return fmt.Errorf("지갑 파일을 마셜 하는 도중 에러가 발생했습니다. : %w\n", err)
	}

	f, err := os.Create(walletFile)

	if err != nil {
		return fmt.Errorf("지갑 파일을 생성하는 도중 에러가 발생했습니다. : %w\n", err)
	}

	if _, err := f.Write(content); err != nil {
		return fmt.Errorf("지갑 파일에 쓰는 도중 에러가 발생했습니다. : %w\n", err)
	}

	return nil
}
