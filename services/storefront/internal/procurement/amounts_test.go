package procurement

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestDeliverySupplierExactAmountsAndCeilings(t *testing.T) {
	for raw, want := range map[string]int64{"0": 0, "0.445": 445000, "0.38": 380000, "0.000001": 1, "20.123456": 20123456} {
		v := json.Number(raw)
		got, e := usdMicros(&v)
		if e != nil || got != want {
			t.Fatal(raw, got, e)
		}
	}
	for _, raw := range []string{"-1", "NaN", "Infinity", "0.4450001", "100000001", "1/2", "1e999"} {
		v := json.Number(raw)
		if _, e := usdMicros(&v); e == nil {
			t.Fatal("invalid amount accepted", raw)
		}
	}
	if _, e := usdMicros(nil); e == nil {
		t.Fatal("missing amount accepted")
	}
	for raw, valid := range map[string]bool{"0.445": true, "0.6": true, "0.600001": false, "0": false} {
		v := json.Number(raw)
		e := validateNY([]NYProduct{{ID: "service", Name: "Gemini", Price: &v, Stock: 1}}, "service", 60)
		if (e == nil) != valid {
			t.Fatal("ceiling rounded", raw, e)
		}
	}
}

func TestDeliverySupplierStockFailureIsDistinctFromInvalidQuote(t *testing.T) {
	price := json.Number("0.41")
	for _, tc := range []struct {
		name, availability string
		stock              int64
		price              *json.Number
		outOfStock, valid  bool
	}{
		{"zero stock", "available", 0, &price, true, false},
		{"sold out", "sold_out", 10, &price, true, false},
		{"available", "available", 1, &price, false, true},
		{"legacy available", "", 1, &price, false, true},
		{"invalid stock", "available", -1, &price, false, false},
		{"invalid price", "available", 1, nil, false, false},
		{"disabled", "disabled", 10, &price, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := validateNY([]NYProduct{{ID: "service", Name: "Gemini", Price: tc.price, Stock: tc.stock, Availability: tc.availability}}, "service", 60)
			if errors.Is(err, ErrOutOfStock) != tc.outOfStock || (err == nil) != tc.valid {
				t.Fatalf("unexpected validation result: %v", err)
			}
		})
	}
	// A service removed from the upstream catalogue is stock, not an outage:
	// checkout must answer 409 instead of an upstream failure.
	if err := validateNY([]NYProduct{{ID: "another", Stock: 0}}, "service", 60); !errors.Is(err, ErrOutOfStock) {
		t.Fatal("delisted product not classified as out of stock", err)
	}
	orderable := false
	if err := validateNY([]NYProduct{{ID: "service", Name: "Gemini", Price: &price, Stock: 10, Orderable: &orderable}}, "service", 60); !errors.Is(err, ErrOutOfStock) {
		t.Fatal("unorderable product not classified as out of stock", err)
	}
	if err := validateNY([]NYProduct{{ID: "service", Name: "Gemini", Price: &price, Stock: 10}}, "service", 40); !errors.Is(err, ErrPriceChanged) || errors.Is(err, ErrOutOfStock) {
		t.Fatal("above-ceiling quote not classified as a price change", err)
	}
	available := 3.0
	inStock := CanbosoProduct{ID: "p1", Price: CanbosoPrice{Amount: 9.87, Currency: "USD"}}
	inStock.Availability.Available = &available
	if err := validateCanboso([]CanbosoProduct{inStock}, "p1", 900); !errors.Is(err, ErrPriceChanged) {
		t.Fatal("canboso above-ceiling quote not classified as a price change", err)
	}
}
