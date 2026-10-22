package procurement

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// The product catalogue runs to several hundred kilobytes (511 entries), far
// past the 256 KB bound that suits the other suppliers, so it gets its own
// limit while the UTF-8 and size checks stay.
const canbosoBodyLimit = 4 << 20

func readCanboso(res *http.Response) ([]byte, error) {
	b, e := io.ReadAll(io.LimitReader(res.Body, canbosoBodyLimit+1))
	if e != nil || len(b) > canbosoBodyLimit || !utf8.Valid(b) {
		return nil, errors.New("invalid bounded supplier response")
	}
	return b, nil
}

// Canboso is a Telegram "buyer key" shop: the wallet pays for the goods and the
// purchase response is the only place the delivered credentials appear. The
// order-history endpoint never returns delivery content, so a purchase that was
// accepted without an immediate completed delivery can only be closed by hand.
var (
	ErrCanbosoBalance   = errors.New("canboso wallet balance is not enough")
	ErrCanbosoInventory = errors.New("canboso inventory is not enough")
)

var (
	canbosoKeyPattern   = regexp.MustCompile(`^tgb_[A-Za-z0-9]{16,96}$`)
	canbosoCodePattern  = regexp.MustCompile(`^[A-Za-z0-9_-]{4,64}$`)
	canbosoIDPattern    = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)
	canbosoAmountLimits = 100000000.0
)

type CanbosoProduct struct {
	ID           string `json:"productId"`
	Name         string `json:"name"`
	Type         string `json:"productType"`
	Price        CanbosoPrice
	Availability struct {
		Available *float64 `json:"available"`
	} `json:"availability"`
}

type CanbosoPrice struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type CanbosoAccount struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	VerifyEmail string `json:"verifyEmail"`
	ExpiryText  string `json:"expiryText"`
	OtherInfo   string `json:"otherInfo"`
}

type CanbosoOrder struct {
	Code     string `json:"orderCode"`
	Status   string `json:"status"`
	Product  string `json:"productId"`
	Fulfill  string `json:"fulfillmentStatus"`
	Quantity int64  `json:"quantity"`
	Delivery struct {
		Accounts []CanbosoAccount `json:"accounts"`
	} `json:"delivery"`
	Payment struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	} `json:"payment"`
}

type CanbosoSupplier struct {
	Base, Key string
	HTTP      *http.Client
}

// NewCanbosoSupplier accepts only an https origin and a buyer key shaped like
// the one the platform issues; the wallet is charged by the purchase call, so
// there is no configurable gateway URL.
func NewCanbosoSupplier(base, key string) (*CanbosoSupplier, error) {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "https://canboso.com"
	}
	u, e := url.Parse(base)
	if e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") || !canbosoKeyPattern.MatchString(key) {
		return nil, errors.New("invalid Canboso configuration")
	}
	return &CanbosoSupplier{Base: "https://" + u.Host, Key: key, HTTP: client(25)}, nil
}

func (c *CanbosoSupplier) call(ctx context.Context, method, path, idempotency string, body, out any) error {
	var rdr *bytes.Reader
	if body != nil {
		b, e := json.Marshal(body)
		if e != nil {
			return e
		}
		rdr = bytes.NewReader(b)
	} else {
		rdr = bytes.NewReader(nil)
	}
	req, e := http.NewRequestWithContext(ctx, method, c.Base+path, rdr)
	if e != nil {
		return e
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-buyer-key", c.Key)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if idempotency != "" {
		req.Header.Set("Idempotency-Key", idempotency)
	}
	res, e := c.HTTP.Do(req)
	if e != nil {
		return e
	}
	defer res.Body.Close()
	raw, e := readCanboso(res)
	if e != nil {
		return e
	}
	if res.StatusCode != http.StatusOK {
		return canbosoError(res.StatusCode, raw)
	}
	if out == nil {
		return nil
	}
	if e = json.Unmarshal(raw, out); e != nil {
		return errors.New("invalid canboso response")
	}
	return nil
}

func canbosoError(code int, raw []byte) error {
	var body struct {
		Message string `json:"message"`
		Code    string `json:"code"`
		Payment *struct {
			Balance float64 `json:"balance"`
		} `json:"payment"`
	}
	_ = json.Unmarshal(raw, &body)
	msg := strings.TrimSpace(body.Message)
	switch code {
	case http.StatusBadRequest:
		if body.Payment != nil || strings.Contains(msg, "余额") || strings.Contains(strings.ToLower(msg), "balance") {
			return ErrCanbosoBalance
		}
		return ErrRejected
	case http.StatusNotFound:
		return ErrRejected
	case http.StatusConflict:
		return ErrCanbosoInventory
	case http.StatusTooManyRequests:
		return errors.New("canboso rate limited")
	default:
		if msg == "" {
			msg = "canboso upstream error"
		}
		return errors.New(msg)
	}
}

func (c *CanbosoSupplier) Products(ctx context.Context) ([]CanbosoProduct, error) {
	var out struct {
		Success  bool             `json:"success"`
		Products []CanbosoProduct `json:"products"`
	}
	if e := c.call(ctx, http.MethodGet, "/api/v2/telegram-buyer/products", "", nil, &out); e != nil {
		return nil, e
	}
	if !out.Success {
		return nil, errors.New("canboso product list rejected")
	}
	return out.Products, nil
}

// Purchase charges the wallet. The idempotency key must stay the same for a
// local order across every retry, otherwise a retried write buys twice.
func (c *CanbosoSupplier) Purchase(ctx context.Context, product, idempotency string) (CanbosoOrder, error) {
	var order CanbosoOrder
	if !canbosoIDPattern.MatchString(product) || len(idempotency) < 8 || len(idempotency) > 128 {
		return order, errors.New("invalid canboso purchase request")
	}
	var out struct {
		Success bool         `json:"success"`
		Order   CanbosoOrder `json:"order"`
	}
	body := map[string]any{"key": c.Key, "product_id": product, "quantity": 1}
	if e := c.call(ctx, http.MethodPost, "/api/v2/telegram-buyer/purchase", idempotency, body, &out); e != nil {
		return order, e
	}
	if !out.Success || !canbosoCodePattern.MatchString(out.Order.Code) {
		return order, errors.New("invalid canboso purchase response")
	}
	return out.Order, nil
}

func (c *CanbosoSupplier) Order(ctx context.Context, code string) (CanbosoOrder, error) {
	var order CanbosoOrder
	if !canbosoCodePattern.MatchString(code) {
		return order, errors.New("invalid canboso order code")
	}
	var out struct {
		Success bool           `json:"success"`
		Orders  []CanbosoOrder `json:"orders"`
	}
	if e := c.call(ctx, http.MethodGet, "/api/v2/telegram-buyer/orders?orderCode="+url.QueryEscape(code), "", nil, &out); e != nil {
		return order, e
	}
	for _, o := range out.Orders {
		if o.Code == code {
			return o, nil
		}
	}
	return order, ErrRejected
}

func canbosoUSDCents(amount float64) (int64, error) {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount < 0 || amount > canbosoAmountLimits || math.Abs(amount*100-math.Round(amount*100)) > 0.000001 {
		return 0, errors.New("invalid canboso amount")
	}
	return int64(math.Round(amount * 100)), nil
}

// validateCanboso rejects a product that is gone, priced over the ceiling or
// quoted in a currency the CNY conversion does not cover.
func validateCanboso(products []CanbosoProduct, id string, ceiling int64) error {
	for _, p := range products {
		if p.ID != id {
			continue
		}
		if p.Availability.Available != nil && *p.Availability.Available <= 0 {
			return ErrOutOfStock
		}
		if p.Price.Currency != "USD" {
			return errors.New("unexpected canboso currency")
		}
		cents, e := canbosoUSDCents(p.Price.Amount)
		if e != nil {
			return e
		}
		if cents > ceiling {
			return errors.New("canboso price above ceiling")
		}
		return nil
	}
	return ErrOutOfStock
}

// canbosoValue renders one delivered account as a single line: the delivery
// validator rejects control characters, so no newlines are allowed.
func canbosoValue(a CanbosoAccount) (string, error) {
	user := strings.TrimSpace(a.User)
	if user == "" || strings.ContainsAny(a.User+a.Password+a.VerifyEmail+a.ExpiryText, "\r\n\t") {
		return "", errors.New("invalid canboso credentials")
	}
	parts := []string{user}
	if p := strings.TrimSpace(a.Password); p != "" {
		parts = append(parts, "密码："+p)
	}
	if v := strings.TrimSpace(a.VerifyEmail); v != "" {
		parts = append(parts, "验证邮箱："+v)
	}
	if v := strings.TrimSpace(a.ExpiryText); v != "" {
		parts = append(parts, "到期："+v)
	}
	value := strings.Join(parts, "  ")
	if len(value) > 500 {
		return "", errors.New("canboso credentials too long")
	}
	return value, nil
}
