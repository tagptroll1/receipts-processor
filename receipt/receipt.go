package receipt

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// state is a persistent storage if the receipts sent for processing
// Technically not required as we're just returning points
var state = map[string]Receipt{}

// receipt -> uuid to points
// sync.Map matches the usecase of only ever writing once and reading multiple times
var points = sync.Map{}

var ErrFailedToCreateReceipt = errors.New("failed to create receipt")
var ErrReceiptNotFound = errors.New("receipt not found")
var ErrFailedToFetchPoints = errors.New("failed to cast receipt points to int")
var ErrPointNotReady = errors.New("points are not ready yet, try again later")

// NewReceipt generated an empty recept with a new random uuid id
func SaveNewReceipt(ctx context.Context, rec Receipt) (string, []error) {
	// Retry logic, which would include network, db hiccups etc
	// In this implementation we just use an exists, it shouldnt happen with uuid however..
	// If it was network errors we could check for that and optionally start a background job to save with increasing backup timeouts
	// and return the id early to prevent holding up the api response. This would require a memory state like this tho.
	for range 5 {
		id := uuid.New().String()

		if _, ok := state[id]; ok {
			continue
		}

		errs, ok := ValidateReceipt(ctx, rec)
		if !ok {
			return "", errs
		}

		state[id] = rec
		// lets us hold for a bit incase calcultions arent done when fetching
		points.Store(id, -1)

		// Make calculating points a one time job that is started now, and should finish
		// before the score is requested
		go func() {
			pnts, err := CalculatePoints(ctx, rec)
			if err != nil {
				// -2 is error
				points.Store(id, -2)
				return
			}

			points.Store(id, pnts)
		}()
		return id, nil
	}
	return "", []error{ErrFailedToCreateReceipt}
}

// GetReceipt will return a Receipt based on id
func GetReceipt(_ context.Context, id string) (Receipt, error) {
	// Retry logic here doesnt make sense, just fyi
	for range 2 {
		switch rec, ok := state[id]; {
		case ok:
			return rec, nil
		// case check some error from db...
		// !ok represents 0 results from db
		case !ok:
			return rec, ErrReceiptNotFound
		}
	}
	return Receipt{}, ErrReceiptNotFound
}

func GetReceiptPoints(_ context.Context, id string) (int, error) {
	for range 5 {
		switch pnts, ok := points.Load(id); {
		case ok:
			p, ok := pnts.(int)
			if !ok {
				return 0, ErrFailedToFetchPoints
			}
			if p == -1 {
				// Still waiting on calculate results
				// Should not hold the response for more than 250ms
				//nolint:staticcheck // SA1004 ignore this!
				time.Sleep(20)
				continue
			}

			if p == -2 {
				return 0, ErrFailedToFetchPoints
			}
			return p, nil
		case !ok:
			return 0, ErrReceiptNotFound
		}
	}

	return 0, ErrPointNotReady
}
