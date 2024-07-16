package app

import (
	"errors"
)

type Operation struct {
	accountID string
	amount    float64
	result    chan float64
	err       chan error
}

func (s *AccountsController) HandleDeposits() {
	for op := range s.depositCh {

		// s.mutex.Lock()
		account, exists := s.accounts[op.accountID]
		// s.mutex.Unlock()
		if !exists {
			op.err <- errors.New("account not found")
			continue
		}
		err := account.Deposit(op.amount)
		op.err <- err
	}
}

func (s *AccountsController) HandleWithdrawals() {
	for op := range s.withdrawCh {
		// s.mutex.Lock()
		account, exists := s.accounts[op.accountID]
		// s.mutex.Unlock()
		if !exists {
			op.err <- errors.New("account not found")
			continue
		}
		err := account.Withdraw(op.amount)
		op.err <- err
	}
}

func (c *AccountsController) HandleBalances() {
	for op := range c.balanceCh {
		// c.mutex.Lock()
		account, exists := c.accounts[op.accountID]
		// c.mutex.Unlock()
		if !exists {
			op.err <- errors.New("account not found")
			continue
		}
		balance := account.GetBalance()
		op.result <- balance
	}
}
