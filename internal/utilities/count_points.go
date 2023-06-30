package utilities

import (
	"math"
	"receipt-processor-challenge/internal/entity"
	"receipt-processor-challenge/pkg/errors"
	"regexp"
	"strings"
)

func CountPoints(receipt *entity.Receipt) (float64, error) {
	if receipt == nil {
		return 0, errors.ErrNoObj
	}

	points := 0
	// Rule 1: One point for every alphanumeric character in the retailer name

	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	points += len(reg.ReplaceAllString(receipt.Retailer, ""))

	// Rule 2: 50 points if the total is a round dollar amount with no cents

	if receipt.Total != 0 && math.Mod(receipt.Total, 1) == 0 {
		points += 50
	}
	// Rule 3: 25 points if the total is a multiple of 0.25
	if receipt.Total > 0 && math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
	}
	// Rule 4: 5 points for every two items on the receipt
	points += 5 * (len(receipt.Items) / 2)

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}
	// Rule 6: 6 points if the day in the purchase date is odd
	if receipt.PurchaseDate.Day()%2 != 0 && !receipt.PurchaseDate.IsZero() {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	if !receipt.PurchaseTime.IsZero() && receipt.PurchaseTime.Hour() >= 14 && receipt.PurchaseTime.Hour() < 16 {
		points += 10
	}

	return float64(points), nil
}
