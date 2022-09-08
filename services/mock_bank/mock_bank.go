package main

import (
	"auction-back/adapters/database"
	"auction-back/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func main() {
	db := database.Connect()
	engine := gin.Default()

	engine.POST("transaction", func(ctx *gin.Context) {
		var input struct {
			Date          *Time                   `json:"date"`
			State         models.TransactionState `json:"state"`
			Type          models.TransactionType  `json:"type" binding:"required"`
			Currency      *models.CurrencyEnum    `json:"currency"`
			Amount        decimal.Decimal         `json:"amount" binding:"required"`
			Error         *string                 `json:"error"`
			AccountFromID *string                 `json:"accountFrom"`
			AccountToID   *string                 `json:"accountTo"`
			OfferID       *string                 `json:"order"`
		}

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		transaction := models.Transaction{}

		transaction.State = input.State

		if input.Date == nil && transaction.State == models.TransactionStateSucceeded {
			now := time.Now()
			transaction.Date = &now
		} else {
			transaction.Date = input.Date.Into()
		}
		transaction.Type = input.Type
		if input.Currency == nil {
			transaction.Currency = models.CurrencyEnumRub
		} else {
			transaction.Currency = *input.Currency
		}
		transaction.Amount = input.Amount
		transaction.Error = input.Error
		transaction.AccountFromID = input.AccountFromID
		transaction.AccountToID = input.AccountToID
		transaction.OfferID = input.OfferID

		if err := db.Transaction().Create(&transaction); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, (&Transaction{}).From(&transaction))
	})

	engine.PATCH("transaction", func(ctx *gin.Context) {
		var input struct {
			ID            int                      `json:"id"`
			Date          *Time                    `json:"date"`
			State         *models.TransactionState `json:"state"`
			Type          *models.TransactionType  `json:"type" binding:"required"`
			Currency      *models.CurrencyEnum     `json:"currency"`
			Amount        *decimal.Decimal         `json:"amount" binding:"required"`
			Error         *string                  `json:"error"`
			AccountFromID *string                  `json:"accountFrom"`
			AccountToID   *string                  `json:"accountTo"`
			OfferID       *string                  `json:"order"`
		}

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		transaction, err := db.Transaction().Get(input.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("db.Transaction().Get: %w").Error()})
			return
		}

		prevState := transaction.State
		if input.State != nil {
			transaction.State = *input.State
		}

		if input.Date == nil && transaction.State == models.TransactionStateSucceeded {
			now := time.Now()
			transaction.Date = &now
		} else if input.Date != nil {
			transaction.Date = input.Date.Into()
		}

		if input.Type != nil {
			transaction.Type = *input.Type
		}
		if input.Currency != nil {
			transaction.Currency = *input.Currency
		}
		if input.Amount != nil {
			transaction.Amount = *input.Amount
		}
		if input.Error != nil {
			transaction.Error = input.Error
		}
		if input.AccountFromID != nil {
			transaction.AccountFromID = input.AccountFromID
		}
		if input.AccountToID != nil {
			transaction.AccountToID = input.AccountToID
		}
		if input.OfferID != nil {
			transaction.OfferID = input.OfferID
		}

		if err := db.Transaction().Update(&transaction); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if prevState != transaction.State {
			go processChangeStateSideEffects(transaction, prevState)
		}

		ctx.JSON(http.StatusOK, (&Transaction{}).From(&transaction))
	})

	engine.Run()
}

func processChangeStateSideEffects(transaction models.Transaction, prevState models.TransactionState) {
	if transaction.Type == models.TransactionTypeBuy {
		if prevState == models.TransactionStateProcessing {
			if transaction.State == models.TransactionStateSucceeded {

			}
		}
	}
}
