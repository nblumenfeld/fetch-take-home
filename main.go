package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nblumenfeld/fetch-take-home/helpers"
	"github.com/nblumenfeld/fetch-take-home/models"
	"go.uber.org/zap"
)

var receipts = []models.Receipt{}
var receiptPoints = []models.Points{}

func main() {
	router := gin.Default()
	router.POST("/receipts/process", postReceipts)
	router.GET("/receipts/:id/points", getReceiptPoints)

	router.Run("localhost:8080")
}

// add a receipt to the "database" and process how many points it is worth
func postReceipts(c *gin.Context) {
	var newReceipt models.Receipt

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Bind the received JSON to newReceipt
	if err := c.BindJSON(&newReceipt); err != nil {
		sugar.Error("Error processing receipt from JSON")
		return
	}

	if result, errMsg := helpers.ValidateReceipt(newReceipt); !result {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": errMsg})
		return
	}

	newReceipt.ID = uuid.New()
	points := helpers.CalculateTotalPoints(newReceipt)

	// add the receipt to the "db" and save it's points
	receipts = append(receipts, newReceipt)
	receiptPoints = append(receiptPoints, models.Points{ID: newReceipt.ID, Points: points})

	c.IndentedJSON(http.StatusCreated, gin.H{"id": newReceipt.ID})
}

// get a receipt by ID
func getReceiptPoints(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))

	for _, a := range receiptPoints {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, gin.H{"points": a.Points})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
}
