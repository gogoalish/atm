package app

import (
	"errors"
	"time"
)

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

type Account struct {
	ID        int32
	CreatedAt time.Time

	balance float64
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	a.balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {

	if amount > a.balance {
		return errors.New("insufficient funds")
	}
	a.balance -= amount
	return nil
}

func (a *Account) GetBalance() float64 {

	return a.balance
}
