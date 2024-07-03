package receipt

import (
	"context"
	"testing"
)

func TestCalculatePrice(t *testing.T) {
	tests := []struct {
		name  string
		input Receipt
		want  int
	}{
		{
			name: "receipt to points #1",
			input: Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items: []Item{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            "6.49",
					},
					{
						ShortDescription: "Emils Cheese Pizza",
						Price:            "12.25",
					},
					{
						ShortDescription: "Knorr Creamy Chicken",
						Price:            "1.26",
					},
					{
						ShortDescription: "Doritos Nacho Cheese",
						Price:            "3.35",
					},
					{
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            "12.00",
					},
				},
			},
			want: 28,
		},

		{
			name: "receipt to points #2",
			input: Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Total:        "9.00",
				Items: []Item{
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
				},
			},
			want: 109,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculatePoints(context.Background(), tt.input)
			if err != nil {
				t.Fail()
			}

			if got != tt.want {
				t.Logf("got %d, want %d", got, tt.want)
				t.Fail()
			}
		})
	}
}
