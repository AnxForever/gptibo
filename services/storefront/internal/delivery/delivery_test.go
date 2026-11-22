package delivery

import (
	"strings"
	"testing"
)

func TestDevinClaimLinkUsesItsOwnMail(t *testing.T) {
	m, e := Message("link", "devin-max-1m", "YMX-TEST", "https://claim.example.invalid/devin/abc", "用你的 Devin 账号打开链接。", "https://shop.example.invalid")
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(m.Subject, "Devin Max") {
		t.Fatalf("subject does not name the product: %q", m.Subject)
	}
	// The generic link template is written for the Gemini 18-month line; a
	// Devin order must never ship that copy.
	if strings.Contains(m.Subject, "Gemini") || strings.Contains(m.HTML, "Gemini") {
		t.Fatal("Gemini copy leaked into the Devin mail")
	}
	if !strings.Contains(m.HTML, "https://claim.example.invalid/devin/abc") {
		t.Fatal("claim link missing from the mail")
	}
}

func TestDevinClaimLinkMustBeHTTPS(t *testing.T) {
	if _, e := Message("link", "devin-max-1m", "YMX-TEST", "http://claim.example.invalid/devin/abc", "说明", "https://shop.example.invalid"); e == nil {
		t.Fatal("plain HTTP claim link accepted")
	}
}
