package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrossSiteIgnoresVisitorBorneRequests(t *testing.T) {
	cases := []struct {
		name    string
		headers map[string]string
		want    bool
	}{
		{"fetch metadata cross-site", map[string]string{"Sec-Fetch-Site": "cross-site"}, true},
		{"fetch metadata same-origin", map[string]string{"Sec-Fetch-Site": "same-origin"}, false},
		{"fetch metadata same-site", map[string]string{"Sec-Fetch-Site": "same-site"}, false},
		{"fetch metadata none", map[string]string{"Sec-Fetch-Site": "none"}, false},
		{"foreign origin", map[string]string{"Origin": "https://evil.example"}, true},
		{"matching origin", map[string]string{"Origin": "https://anxforever.shop"}, false},
		{"foreign referer", map[string]string{"Referer": "https://evil.example/page"}, true},
		{"matching referer", map[string]string{"Referer": "https://anxforever.shop/orders"}, false},
		{"unparseable referer", map[string]string{"Referer": "::"}, true},
		{"no origin evidence counts", nil, false},
	}
	for _, c := range cases {
		r := httptest.NewRequest(http.MethodGet, "https://anxforever.shop/api/orders/x", nil)
		for k, v := range c.headers {
			r.Header.Set(k, v)
		}
		if got := crossSite(r); got != c.want {
			t.Fatalf("%s: crossSite=%v want %v", c.name, got, c.want)
		}
	}
}
