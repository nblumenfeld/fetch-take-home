package helpers

import (
	"testing"

	"github.com/nblumenfeld/fetch-take-home/models"
)

var targetReceipt = models.Receipt{
	Retailer:     "Target",
	PurchaseDate: "2022-01-01",
	PurchaseTime: "13:01",
	Total:        "35.35",
	Items: []models.Item{
		models.Item{
			ShortDescription: "Mountain Dew 12PK",
			Price:            "6.49",
		},
		models.Item{
			ShortDescription: "Emils Cheese Pizza",
			Price:            "12.25",
		},
		models.Item{
			ShortDescription: "Knorr Creamy Chicken",
			Price:            "1.26",
		},
		models.Item{
			ShortDescription: "Doritos Nacho Cheese",
			Price:            "3.35",
		},
		models.Item{
			ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
			Price:            "12.00",
		},
	},
}

// I wanted to see how it was to set up unit tests for a golang project
// The following are all positive tests related to the target example that was provided
// This is not exhaustive of all testing angles that could be taken, but it
// was fun to learn how to do this in Go!

// Larger function to test everything in one go
func TestPointCalculations(t *testing.T) {

	points := CalculateTotalPoints(targetReceipt)

	if points != 28 {
		t.Fatalf("Points %v did not match %v", points, 28)
	}
}

func TestAlphNumericCalculations(t *testing.T) {
	points := CalculateAlphaNumeric(targetReceipt)
	if points != 6 {
		t.Fatalf("Points %v did not match %v for alpha numeric points", points, 6)
	}
}
func TestTotalCostCalculations(t *testing.T) {
	points := CalculateTotalCostPoints(targetReceipt)
	if points != 0 {
		t.Fatalf("Points %v did not match %v for total cost points", points, 0)
	}
}
func TestItemCalculations(t *testing.T) {
	points := CalculateItemPoints(targetReceipt)
	if points != 16 {
		t.Fatalf("Points %v did not match %v for item points", points, 16)
	}
}
func TestDatecCalculations(t *testing.T) {
	points := CalculateDatePoints(targetReceipt)
	if points != 6 {
		t.Fatalf("Points %v did not match %v for date points", points, 6)
	}
}
func TestTimeCalculations(t *testing.T) {
	points := CalculateTimePoints(targetReceipt)
	if points != 0 {
		t.Fatalf("Points %v did not match %v for time points", points, 0)
	}
}

func TestReceiptValidation(t *testing.T) {
	result, errMsg := ValidateReceipt(targetReceipt)

	if !result {
		t.Fatalf("Failed validation with message %s", errMsg)
	}
}
