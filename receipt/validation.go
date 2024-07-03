package receipt

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	ErrInvalidReceiptRetailer = errors.New("invalid receipt retailer")
	ErrInvalidReceiptDate     = errors.New("invalid receipt purchaseDate")
	ErrInvalidReceiptTime     = errors.New("invalid receipt purchaseTime")
	ErrInvalidReceiptItems    = errors.New("invalid receipt items")
	ErrInvalidReceiptTotal    = errors.New("invalid receipt total")

	ErrInvalidItemDesc  = errors.New("invalid receipt item shortDescription")
	ErrInvalidItemPrice = errors.New("invalid receipt item price")
)

var (
	regText, _  = regexp.Compile(`^[\w\s\-&]+$`)
	regFloat, _ = regexp.Compile(`^\d+\.\d{2}$`)
)

func ValidateReceipt(ctx context.Context, rec Receipt) ([]error, bool) {
	var errs []error

	if !regText.MatchString(rec.Retailer) {
		errs = append(errs, fmt.Errorf("%w: '%s' is not a valid retailer", ErrInvalidReceiptRetailer, rec.Retailer))
	}

	if _, err := time.Parse(time.DateOnly, rec.PurchaseDate); err != nil {
		errs = append(errs, fmt.Errorf("%w, '%s' is not a valid purchaseDate", ErrInvalidReceiptDate, rec.PurchaseDate))
	}

	if _, err := time.Parse(time.TimeOnly, rec.PurchaseTime); err != nil {
		errs = append(errs, fmt.Errorf("%w, '%s' is not a valid purchaseTime", ErrInvalidReceiptTime, rec.PurchaseTime))
	}

	if len(rec.Items) == 0 {
		errs = append(errs, fmt.Errorf("%w, receipt doesnt contain any items", ErrInvalidReceiptItems))
	}

	if !regFloat.MatchString(rec.Total) {
		errs = append(errs, fmt.Errorf("%w, '%s' is not a valid total", ErrInvalidReceiptTotal, rec.Total))
	}

	for i, item := range rec.Items {
		if itemErrs := ValidateReceiptItem(ctx, i, item); len(itemErrs) > 0 {
			errs = append(errs, itemErrs...)
		}
	}

	return errs, true
}

func ValidateReceiptItem(_ context.Context, i int, item Item) []error {
	var errs []error
	if !regText.MatchString(item.ShortDescription) {
		errs = append(errs, fmt.Errorf("%w, '%s' is not a valid shortDescription for item (%d)", ErrInvalidItemDesc, item.ShortDescription, i))
	}

	if !regFloat.MatchString(item.Price) {
		errs = append(errs, fmt.Errorf("%w, '%s' is not a valid price for item (%d)", ErrInvalidItemPrice, item.Price, i))
	}

	return errs
}
