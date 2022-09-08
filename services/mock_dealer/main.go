package main

import (
	"auction-back/adapters/bank"
	"auction-back/adapters/database"
	"auction-back/models"
	"auction-back/ports"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
)

func main() {
	db := database.Connect()
	bank := bank.New(&db)
	engine := gin.Default()

	dealer := NewDealer(&db, &bank)
	defer dealer.Stop()
	go dealer.Run()

	engine.POST("finish", func(ctx *gin.Context) {
		var input struct {
			AuctionId string `json:"auctionId"`
		}

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		auctions, err := dealer.Process(ports.FindAndSetFinishConfig{
			Filter: &models.AuctionsFilter{
				IDs: []string{input.AuctionId},
			},
			DefaultDuration:      "0 mins",
			TimeGapFromLastOffer: "0 mins",
		})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, auctions)
	})

	engine.POST("owner-accept", func(ctx *gin.Context) {
		var input struct {
			OfferId string  `json:"offerId"`
			UserId  *string `json:"userId"`
			Comment *string `json:"comment"`
		}

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := dealer.TransferringProduct(input.OfferId, input.UserId, input.Comment)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	engine.POST("consumer-accept", func(ctx *gin.Context) {
		var input struct {
			OfferId string  `json:"offerId"`
			UserId  *string `json:"userId"`
			Comment *string `json:"comment"`
		}

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := dealer.Succeeded(input.OfferId, input.UserId, input.Comment)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	engine.POST("fail", func(ctx *gin.Context) {
		var input struct {
			OfferId string  `json:"offerId"`
			UserId  *string `json:"userId"`
			Comment *string `json:"comment"`
		}

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var errs error

		err := dealer.TransferMoneyFailed(input.OfferId, input.Comment)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
			return
		} else if !errors.Is(err, ErrWrongState) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		errs = multierror.Append(errs, err)

		err = dealer.TransferProductFailed(input.OfferId, input.UserId, input.Comment)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
			return
		} else if !errors.Is(err, ErrWrongState) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
	})

	engine.Run()
}
