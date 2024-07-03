package receipt

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func CalculatePoints(ctx context.Context, rec Receipt) (int, error) {
	var points int

	// +1 for every alphanumeric char
	for _, c := range rec.Retailer {
		if unicode.IsLetter(c) {
			points += 1
		}
	}

	// +50 if round dollar
	if strings.HasSuffix(rec.Total, ".00") {
		points += 50
	}

	total, err := strconv.ParseFloat(rec.Total, 64)
	if err != nil {
		return 0, err
	}
	// +25 if multiple of .25
	if int(total*100)%25 == 0 {
		points += 25
	}

	// +5 for every 2 items (divide by 2 and floor)
	points += 5 * int(math.Floor(float64(len(rec.Items))/2))

	// +6 if the purchase day date is odd (Forced yyyy-mm-dd format)
	dateParts := strings.Split(rec.PurchaseDate, "-")
	if len(dateParts) != 3 {
		return 0, fmt.Errorf("wtf.. invalid date")
	}
	day, err := strconv.Atoi(dateParts[2])
	if err != nil {
		return 0, err
	}
	if isOdd(day) {
		points += 6
	}

	// +10 if hour between 14 and 16 (2pm - 4pm)
	hour, _, _ := strings.Cut(rec.PurchaseTime, ":")
	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		return 0, err
	}

	// Anything at 14:xx is between 14 and anything under 16: is before, ut to 15:59:59
	if hourInt >= 14 && hourInt < 16 {
		points += 10
	}

	pointsForItems, err := calculatePointsForItems(ctx, rec.Items)
	if err != nil {
		return 0, err
	}

	points += pointsForItems

	return points, nil
}

func calculatePointsForItems(_ context.Context, items []Item) (int, error) {
	var points int
	for _, item := range items {
		// +0.2*price for every item that has a trimmed desc length divisible by 3
		desc := strings.TrimSpace(item.ShortDescription)
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return 0, err
		}

		if len(desc)%3 == 0 {
			p := math.Ceil(price * 0.2)
			points += int(p)
		}
	}
	return points, nil
}

func isOdd(n int) bool {
	return n%2 != 0
}
