package helpers

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/nblumenfeld/fetch-take-home/models"
	"go.uber.org/zap"
)

// validates the receipt object
func ValidateReceipt(receipt models.Receipt) (bool, string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	if _, err := time.Parse(time.DateOnly, receipt.PurchaseDate); err != nil {
		sugar.Error(
			"Invalid Purchase Date Format",
			zap.String("purchaseDate", receipt.PurchaseDate),
		)
		return false, "Invalid Purchase Date Format"
	}

	timeLayout := "[0-9]{2}:[0-9]{2}"

	if result, _ := regexp.MatchString(timeLayout, receipt.PurchaseTime); !result {
		sugar.Error(
			"Invalid Purchase Time Format",
			zap.String("purchaseTime", receipt.PurchaseTime),
		)
		return false, "Invalid Purchase Time Format"
	}

	costRegex := regexp.MustCompile("^[0-9]*.[0-9]{2}$")

	if result := costRegex.MatchString(receipt.Total); !result {
		sugar.Error(
			"Invalid Total Format",
			zap.String("total", receipt.Total),
		)
		return false, "Invalid Total Format"
	}

	for _, item := range receipt.Items {
		if result := costRegex.MatchString(item.Price); !result {
			sugar.Error(
				"Invalid Item Format for price",
				zap.String("shortDescription", item.ShortDescription),
				zap.String("price", item.Price),
			)
			return false, "Invalid Item Format for price with decription: " + item.ShortDescription
		}
	}

	return true, ""
}

// calulates the total points for a receipt, based on the rules of the assignment
func CalculateTotalPoints(receipt models.Receipt) int {
	var totalPoints = 0

	totalPoints += CalculateAlphaNumeric(receipt)

	totalPoints += CalculateTotalCostPoints(receipt)

	totalPoints += CalculateItemPoints(receipt)

	totalPoints += CalculateDatePoints(receipt)

	totalPoints += CalculateTimePoints(receipt)

	return totalPoints
}

// calculates the number of alphanumeric characters
func CalculateAlphaNumeric(receipt models.Receipt) int {
	var total = 0

	// One point for every alphanumeric character in the retailer name.
	for _, r := range receipt.Retailer {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			total++
		}
	}

	return total
}

func CalculateTotalCostPoints(receipt models.Receipt) int {
	total := 0
	// 50 points if the total is a round dollar amount with no cents.
	totalParts := strings.Split(receipt.Total, ".")
	if totalParts[1] == "00" {
		total += 50
	}

	// 25 points if the total is a multiple of 0.25.
	cents, _ := strconv.Atoi(totalParts[1])
	if cents%25 == 0 {
		total += 25
	}

	return total
}

func CalculateItemPoints(receipt models.Receipt) int {
	total := 0
	// 5 points for every two items on the receipt.
	total += 5 * (len(receipt.Items) / 2)

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		costParts := strings.Split(item.Price, ".")
		dollars, _ := strconv.Atoi(costParts[0])
		cents, _ := strconv.Atoi(costParts[1])
		cost := (dollars * 100) + cents
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			cost /= 50
			if cost%100 == 0 {
				total += cost / 10
			} else {
				total += (cost / 10) + 1
			}
		}
	}

	return total
}

func CalculateDatePoints(receipt models.Receipt) int {
	total := 0
	// 6 points if the day in the purchase date is odd.
	dateParts := strings.Split(receipt.PurchaseDate, "-")
	day, _ := strconv.Atoi(dateParts[2])
	if day%2 != 0 {
		total += 6
	}

	return total
}

func CalculateTimePoints(receipt models.Receipt) int {
	total := 0
	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	timeParts := strings.Split(receipt.PurchaseTime, ":")
	hours, _ := strconv.Atoi(timeParts[0])
	mins, _ := strconv.Atoi(timeParts[1])
	if (hours == 14 && mins >= 1) || (hours > 14 && hours < 16) {
		total += 10
	}

	return total
}
