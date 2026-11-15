package procurement

import (
	"errors"
	"fmt"
	"gptibo/storefront/internal/fault"
	"testing"
)

func TestQuoteErrorRendersClientConflicts(t *testing.T) {
	for _, tc := range []struct {
		name   string
		in     error
		status int
	}{
		{"stock", ErrOutOfStock, 409},
		{"wrapped stock", fmt.Errorf("quote: %w", ErrOutOfStock), 409},
		{"price", ErrPriceChanged, 409},
		{"drifted response stays an upstream failure", errors.New("supplier goods unavailable or quote changed"), 0},
		{"no error", nil, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := quoteError(tc.in)
			if tc.status == 0 {
				if got != tc.in {
					t.Fatalf("error rewritten: %v", got)
				}
				return
			}
			if fault.Status(got) != tc.status {
				t.Fatalf("status=%d want %d", fault.Status(got), tc.status)
			}
			if fault.Message(got) == "" || fault.Message(got) == fault.Message(errors.New("x")) {
				t.Fatal("conflict reused the generic upstream message")
			}
		})
	}
}
