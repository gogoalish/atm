package app

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gogoalish/atm/internal/logger"

	"go.uber.org/zap"
)

type AccountsController struct {
	accounts   map[string]BankAccount
	mutex      sync.Mutex
	depositCh  chan Operation
	withdrawCh chan Operation
	balanceCh  chan Operation
	nextID     int
}

func NewAccountsController() *AccountsController {
	return &AccountsController{
		accounts:   make(map[string]BankAccount),
		depositCh:  make(chan Operation),
		withdrawCh: make(chan Operation),
		balanceCh:  make(chan Operation),
		mutex:      sync.Mutex{},
		nextID:     1,
	}
}

var ErrNoLogger = errors.New("logger not found in context")

func (c *AccountsController) Create(ctx *gin.Context) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	l.Info("Received request to create account")
	id := strconv.Itoa(c.nextID)
	c.accounts[id] = &Account{ID: int32(c.nextID)}

	c.nextID++

	l.Info("Account created successfully", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

type depositReqBody struct {
	Amount float64 `json:"amount" binding:"required,min=1"`
}

func (c *AccountsController) Deposit(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	id := ctx.Param("id")
	var bodyReq depositReqBody

	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	l.Info("Received request to deposit", zap.String("account_id", id), zap.Float64("amount", bodyReq.Amount))

	op := Operation{accountID: id, amount: bodyReq.Amount, result: make(chan float64), err: make(chan error)}
	c.depositCh <- op

	err := <-op.err
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	l.Info("Deposit made", zap.String("account_id", id), zap.Float64("amount", bodyReq.Amount))
	ctx.JSON(http.StatusOK, gin.H{"status": "deposit successful"})
}

func (c *AccountsController) Balance(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	id := ctx.Param("id")

	resultCh := make(chan float64)
	errCh := make(chan error)
	op := Operation{accountID: id, result: resultCh, err: errCh}
	c.balanceCh <- op

	select {
	case err := <-errCh:
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case balance := <-resultCh:
		l.Info("Balance checked", zap.String("account_id", id), zap.Float64("balance", balance))
		ctx.JSON(http.StatusOK, gin.H{"balance": balance})
	}
}

type withdrawReqBody struct {
	Amount float64 `json:"amount" binding:"required,min=1"`
}

func (c *AccountsController) Withdraw(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	id := ctx.Param("id")
	var bodyReq withdrawReqBody

	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	l.Info("Received request to withdraw", zap.String("account_id", id), zap.Float64("amount", bodyReq.Amount))

	op := Operation{accountID: id, amount: bodyReq.Amount, result: make(chan float64), err: make(chan error)}
	c.withdrawCh <- op

	err := <-op.err
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	l.Info("Withdrawal made", zap.String("account_id", id), zap.Float64("amount", bodyReq.Amount))
	ctx.JSON(http.StatusOK, gin.H{"status": "deposit successful"})
}
