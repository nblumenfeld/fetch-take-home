package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// receipt object and the in memory cache
type Receipt struct {
	ID           uuid.UUID `json:"id"`
	Retailer     string    `json:"retailer""`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Total        string    `json:"total"`
	Items        []Item    `json:"items"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json"price"`
}

type Points struct {
	ID     uuid.UUID
	Points int
}

var receipts = []Receipt{}
var receiptPoints = []Points{}

func main() {
	router := gin.Default()
	router.POST("/receipts/process", postReceipts)
	router.GET("/receipts/:id/points", getReceiptPoints)

	router.Run("localhost:8080")
}

// add a receipt to the "database" and process how many points it is worth
func postReceipts(c *gin.Context) {
	var newReceipt Receipt

	// fmt.Println("Parsing the Request Body")

	// Bind the received JSON to newReceipt
	if err := c.BindJSON(&newReceipt); err != nil {
		fmt.Println("Error Parsing the Request Body")
		return
	}

	// fmt.Println("Parsed the Request Body")

	if result, errMsg := validateReceipt(newReceipt); !result {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": errMsg})
		return
	}

	newReceipt.ID = uuid.New()
	points := calculateTotalPoints(newReceipt)

	// add the receipt to the "db" and save it's points
	receipts = append(receipts, newReceipt)
	receiptPoints = append(receiptPoints, Points{ID: newReceipt.ID, Points: points})

	c.IndentedJSON(http.StatusCreated, gin.H{"id": newReceipt.ID})
}

// get a receipt by ID
func getReceiptPoints(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	fmt.Println("Parsed the UUID")

	for _, a := range receiptPoints {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, gin.H{"points": a.Points})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
}

// calulates the total points for a receipt, based on the rules of the assignment
func calculateTotalPoints(receipt Receipt) int {
	var totalPoints = 0

	// One point for every alphanumeric character in the retailer name.
	totalPoints += countAlphaNumeric(receipt.Retailer)
	// fmt.Println(totalPoints)

	// 50 points if the total is a round dollar amount with no cents.
	totalParts := strings.Split(receipt.Total, ".")
	if totalParts[1] == "00" {
		totalPoints += 50
	}
	// fmt.Println(totalPoints)

	// 25 points if the total is a multiple of 0.25.
	cents, _ := strconv.Atoi(totalParts[1])
	if cents%25 == 0 {
		totalPoints += 25
	}
	// fmt.Println(totalPoints)

	// 5 points for every two items on the receipt.
	totalPoints += 5 * (len(receipt.Items) / 2)
	// fmt.Println(totalPoints)

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		costParts := strings.Split(item.Price, ".")
		dollars, _ := strconv.Atoi(costParts[0])
		cents, _ := strconv.Atoi(costParts[1])
		cost := (dollars * 100) + cents
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			fmt.Println(cost)
			cost /= 50
			if cost%100 == 0 {
				totalPoints += cost / 10
			} else {
				totalPoints += (cost / 10) + 1
			}
		}
	}
	// fmt.Println(totalPoints)

	// 6 points if the day in the purchase date is odd.
	dateParts := strings.Split(receipt.PurchaseDate, "-")
	day, _ := strconv.Atoi(dateParts[2])
	if day%2 != 0 {
		totalPoints += 6
	}
	// fmt.Println(totalPoints)

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	timeParts := strings.Split(receipt.PurchaseTime, ":")
	hours, _ := strconv.Atoi(timeParts[0])
	mins, _ := strconv.Atoi(timeParts[1])
	if (hours == 14 && mins >= 1) || (hours > 14 && hours < 16) {
		totalPoints += 10
	}
	// fmt.Println(totalPoints)

	return totalPoints
}

// calculates the number of alphanumeric characters
func countAlphaNumeric(s string) int {
	var total = 0

	for _, r := range s {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			total++
		}
	}
	return total
}

// validates the receipt object
func validateReceipt(receipt Receipt) (bool, string) {
	if _, err := time.Parse(time.DateOnly, receipt.PurchaseDate); err != nil {
		fmt.Println(err)
		return false, "Invalid Purchase Date Format"
	}

	timeLayout := "[0-9]{2}:[0-9]{2}"

	if result, _ := regexp.MatchString(timeLayout, receipt.PurchaseTime); !result {
		return false, "Invalid Purchase Time Format"
	}

	costRegex := regexp.MustCompile("^[0-9]*.[0-9]{2}$")

	if result := costRegex.MatchString(receipt.Total); !result {
		return false, "Invalid Total Format"
	}

	for _, item := range receipt.Items {
		if result := costRegex.MatchString(item.Price); !result {
			return false, "Invalid Item Format for price with decription: " + item.ShortDescription
		}
	}

	return true, ""
}
