package main

import (
	"auction-back/adapters/database"
	"auction-back/models"
	"net/http"

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

		transaction.Date = input.Date.Into()
		transaction.State = input.State
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

	engine.Run()
}
