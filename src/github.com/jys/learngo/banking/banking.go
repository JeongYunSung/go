package banking

type Account struct {
	owner   string
	balance int
}

func NewAccount(owner string, balance int) *Account {
	return &Account{owner, balance}
}

func (a Account) Balance() int {
	return a.balance
}

func (a *Account) Deposit(amount int) {
	a.balance += amount
}
