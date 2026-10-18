package procurement

import (
	"gptibo/storefront/internal/delivery"
	"testing"
)

func TestDeliveryValueUpgradesOnlyWellFormedHTTPClaimLinks(t *testing.T) {
	input := "http://claim.example.invalid/claim/a%2Fb?token=x%2By&next=%2F#claim"
	want := "https://claim.example.invalid/claim/a%2Fb?token=x%2By&next=%2F#claim"
	if got := delivery_supplierDeliveryValue(input); got != want || !delivery.ValidValue("link", got) {
		t.Fatal("claim bytes changed", got)
	}
	for _, value := range []string{
		"http://claim.example.invalid:8080/claim",
		"http://user@claim.example.invalid/claim",
		"http://claim.example.invalid/claim\nother",
		"http://claim.example.invalid/claim\"onclick=evil",
		"http://localhost/claim",
		"http:///claim",
	} {
		if delivery_supplierDeliveryValue(value) != value {
			t.Fatal("rewrote an unsafe value", value)
		}
		if delivery.ValidValue("link", delivery_supplierDeliveryValue(value)) {
			t.Fatal("accepted unsafe claim", value)
		}
	}
	if delivery_supplierDeliveryValue(want) != want {
		t.Fatal("changed HTTPS claim")
	}
	// Redeem codes and account credentials must pass through byte for byte.
	for _, value := range []string{"PRIVATE-CDK", "user@example.test\npassword", "https://www.duolingo.com/redeem?code=X"} {
		if delivery_supplierDeliveryValue(value) != value {
			t.Fatal("changed a non-HTTP value", value)
		}
	}
}
