package procurement

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testCanbosoKey = "tgb_0123456789abcdef0123456789abcdef01234567"

func canbosoFixture(t *testing.T, h http.HandlerFunc) *CanbosoSupplier {
	t.Helper()
	ts := httptest.NewTLSServer(h)
	t.Cleanup(ts.Close)
	return &CanbosoSupplier{Base: ts.URL, Key: testCanbosoKey, HTTP: ts.Client()}
}

func TestCanbosoPurchaseCarriesCredentials(t *testing.T) {
	var seenID, seenKey, seenMethod, seenPath string
	var body map[string]any
	c := canbosoFixture(t, func(w http.ResponseWriter, r *http.Request) {
		seenID, seenKey, seenMethod, seenPath = r.Header.Get("Idempotency-Key"), r.Header.Get("x-buyer-key"), r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = w.Write([]byte(`{"success":true,"order":{"orderCode":"ORDER1A2B","status":"completed","productId":"p1","quantity":1,"delivery":{"accounts":[{"user":"buyer@example.com","password":"pw123","verifyEmail":"rec@example.com","expiryText":"2026-11-22"}]},"payment":{"amount":9.87,"currency":"USD"}}}`))
	})
	order, e := c.Purchase(context.Background(), "p1", "cbx_ORDERID01")
	if e != nil {
		t.Fatal(e)
	}
	if seenMethod != http.MethodPost || seenPath != "/api/v2/telegram-buyer/purchase" {
		t.Fatal("wrong endpoint", seenMethod, seenPath)
	}
	if seenKey != testCanbosoKey || seenID != "cbx_ORDERID01" {
		t.Fatal("missing buyer or idempotency header", seenKey, seenID)
	}
	if body["product_id"] != "p1" || body["quantity"] != float64(1) {
		t.Fatal("wrong purchase body", body)
	}
	if order.Code != "ORDER1A2B" || order.Status != "completed" || order.Payment.Amount != 9.87 {
		t.Fatal("order not parsed", order)
	}
	if len(order.Delivery.Accounts) != 1 {
		t.Fatal("delivery lost", order.Delivery)
	}
	value, e := canbosoValue(order.Delivery.Accounts[0])
	if e != nil {
		t.Fatal(e)
	}
	if strings.ContainsAny(value, "\r\n\t") || !strings.Contains(value, "buyer@example.com") || !strings.Contains(value, "pw123") {
		t.Fatal("unexpected delivery value", value)
	}
}

func TestCanbosoBalanceErrorIsRecognised(t *testing.T) {
	c := canbosoFixture(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"success":false,"message":"钱包余额不足","payment":{"currency":"USD","balance":0}}`))
	})
	if _, e := c.Purchase(context.Background(), "p1", "cbx_ORDERID01"); !errors.Is(e, ErrCanbosoBalance) {
		t.Fatal("balance error not classified", e)
	}
}

func TestCanbosoRejectsMalformedConfiguration(t *testing.T) {
	for _, c := range [][2]string{
		{"http://canboso.com", testCanbosoKey},
		{"https://canboso.com/api", testCanbosoKey},
		{"https://canboso.com", "abc"},
		{"https://user:pw@canboso.com", testCanbosoKey},
	} {
		if _, e := NewCanbosoSupplier(c[0], c[1]); e == nil {
			t.Fatal("accepted", c[0])
		}
	}
	if _, e := NewCanbosoSupplier("https://canboso.com", testCanbosoKey); e != nil {
		t.Fatal(e)
	}
}

func TestValidateCanboso(t *testing.T) {
	product := func(id, currency string, price float64, available float64) CanbosoProduct {
		var p CanbosoProduct
		p.ID, p.Price.Currency, p.Price.Amount = id, currency, price
		p.Availability.Available = &available
		return p
	}
	if e := validateCanboso([]CanbosoProduct{product("p1", "USD", 9.87, 3)}, "p1", 987); e != nil {
		t.Fatal(e)
	}
	if e := validateCanboso([]CanbosoProduct{product("p1", "USD", 9.87, 0)}, "p1", 987); !errors.Is(e, ErrOutOfStock) {
		t.Fatal("sold-out product accepted", e)
	}
	if e := validateCanboso([]CanbosoProduct{product("p1", "USD", 9.88, 3)}, "p1", 987); e == nil {
		t.Fatal("price above ceiling accepted")
	}
	if e := validateCanboso([]CanbosoProduct{product("p1", "VND", 9.87, 3)}, "p1", 987); e == nil {
		t.Fatal("foreign currency accepted")
	}
}

func TestCanbosoValueRejectsInjectedCredentials(t *testing.T) {
	if _, e := canbosoValue(CanbosoAccount{User: "a@b.c", Password: "x\ny"}); e == nil {
		t.Fatal("newline accepted")
	}
	if _, e := canbosoValue(CanbosoAccount{User: "  "}); e == nil {
		t.Fatal("empty user accepted")
	}
}

func TestPickKind(t *testing.T) {
	cases := []struct {
		name     string
		sp       supply
		legacy   string
		canboso  bool
		delivery bool
		account  bool
		want     string
	}{
		{"unmapped plan stays unsupported", supply{}, "", true, true, true, ""},
		{"legacy mapping still works", supply{}, "gmail", false, true, true, "gmail"},
		{"catalogue row wins", supply{Supplier: "canboso", Product: "p1", Ceiling: 987}, "", true, true, true, "canboso"},
		{"unconfigured channel falls back to legacy", supply{Supplier: "canboso"}, "gmail", false, true, true, "gmail"},
		{"unconfigured channel cannot claim a new plan", supply{Supplier: "canboso"}, "", false, true, true, ""},
		{"delivery row", supply{Supplier: "delivery_supplier", Product: "S_03", Ceiling: 60}, "", true, true, false, "delivery_supplier"},
		{"delivery row ignored when client missing", supply{Supplier: "delivery_supplier"}, "", true, false, true, ""},
	}
	for _, c := range cases {
		if got := pickKind(c.sp, c.legacy, c.canboso, c.delivery, c.account); got != c.want {
			t.Fatalf("%s: pickKind=%q want %q", c.name, got, c.want)
		}
	}
}
