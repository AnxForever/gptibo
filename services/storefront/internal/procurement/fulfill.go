package procurement

import (
	"context"
	"encoding/json"
	"errors"
	"gptibo/storefront/internal/delivery"
	"gptibo/storefront/internal/orders"
	"math"
	"net/url"
	"strings"
	"time"
)

// Upstreams sometimes hand out claim links over HTTP even when the host only
// serves HTTPS. Upgrade only a well-formed link whose host carries no port,
// userinfo or unsafe bytes; every other value is left untouched so delivery
// validation still decides.
func delivery_supplierDeliveryValue(value string) string {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "http://") || strings.ContainsAny(value, "<>\"'\t\r\n") {
		return value
	}
	u, e := url.Parse(value)
	if e != nil || u.User != nil || u.Host != u.Hostname() || !strings.Contains(u.Hostname(), ".") {
		return value
	}
	return "https://" + strings.TrimPrefix(value, "http://")
}

func (s *Service) gmail(ctx context.Context, p Purchase) error {
	var e error
	if p.Attempted == nil {
		g, e := s.AccountSupplier.Goods(ctx, p.Goods)
		if e != nil {
			return e
		}
		if e = validateGoods(g, p.Ceiling); e != nil {
			return e
		}
		p, e = s.start(ctx, "gmail", p)
		if e != nil {
			return e
		}
		var created GmailOrder
		e = s.AccountSupplier.Call(ctx, "/order/create", map[string]any{"goodsId": p.Goods, "count": 1, "notifyUrl": "", "outerNumber": p.ID, "safePrice": float64(p.Ceiling) / 100, "sku": "", "attach": []any{}}, &created)
		if errors.Is(e, ErrRejected) {
			if err := s.reject(ctx, "gmail", p, "gmail-purchase-rejected", "供应商明确拒绝采购，请核对原订单售后"); err != nil {
				return err
			}
		}
		if e != nil {
			return e
		}
		if e = s.remember(ctx, "gmail", p, created.ID); e != nil {
			return e
		}
		p.Supplier = &created.ID
	}
	// AccountSupplier does not guarantee idempotent create: after any attempted write only query.
	id := ""
	if p.Supplier != nil {
		id = *p.Supplier
	}
	var result GmailOrder
	if e = s.AccountSupplier.Call(ctx, "/order/getDetail", map[string]any{"orderNumber": id, "outerNumber": p.ID}, &result); e != nil {
		return e
	}
	if !identifier.MatchString(result.ID) || result.Merchant != p.ID || (p.Supplier != nil && result.ID != *p.Supplier) {
		return errors.New("supplier order identity mismatch")
	}
	if result.Status == 4 || result.Status == 5 {
		return s.reject(ctx, "gmail", p, "gmail-purchase-rejected", "供应商已确认采购失败或退款，请通过原订单处理售后")
	}
	cost, e1 := cents(result.Money)
	refund, e2 := cents(result.Refund)
	if e1 != nil || e2 != nil || result.Count != 1 || result.Status != 3 || refund != 0 || cost > p.Ceiling || len(result.Cards) != 1 {
		return errors.New("supplier purchase amount or quantity mismatch")
	}
	card := result.Cards[0]
	if card.ID == "" || strings.TrimSpace(card.Number) == "" || card.Password == nil || card.Status == nil || *card.Status != 0 || len(card.Number)+len(*card.Password) > 8000 {
		return errors.New("supplier credentials not verified")
	}
	tx, e := s.Store.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	_, p, e = s.fence(ctx, tx, "gmail", p)
	if e != nil {
		return e
	}
	_, e = tx.Exec(ctx, `UPDATE account_supplier_orders SET state='purchased',supplier_order=$1,actual_cost_cents=$2,cards=$3,updated_at=$4 WHERE order_id=$5`, result.ID, cost, raw(result.Cards), time.Now().UnixMilli(), p.ID)
	if e != nil {
		return e
	}
	return tx.Commit(ctx)
}
func (s *Service) delivery_supplier(ctx context.Context, p Purchase) error {
	var e error
	if p.Attempted == nil {
		products, e := s.DeliverySupplier.Products(ctx)
		if e != nil {
			return e
		}
		if e = validateNY(products, p.Service, p.MaxUSD); e != nil {
			return e
		}
		var profile struct {
			Balance *json.Number `json:"wallet_balance"`
		}
		if e = s.DeliverySupplier.Call(ctx, "/me", "", nil, &profile); e != nil {
			return e
		}
		amount, e := usdMicros(profile.Balance)
		if e != nil || amount < p.MaxUSD*10000 {
			return errors.New("insufficient supplier balance")
		}
		p, e = s.start(ctx, "delivery_supplier", p)
		if e != nil {
			return e
		}
	}
	if p.Supplier == nil {
		var created NYOrder
		e = s.DeliverySupplier.Call(ctx, "/order", p.Key, &p.Body, &created)
		if errors.Is(e, ErrReceiptRejected) {
			if err := s.reject(ctx, "delivery_supplier", p, "delivery_supplier-receipt-rejected", "供应商明确拒绝采购编号，采购未发生，停止自动重试"); err != nil {
				return err
			}
		}
		if e != nil {
			return e
		}
		if e = s.remember(ctx, "delivery_supplier", p, created.ID); e != nil {
			return e
		}
		p.Supplier = &created.ID
	}
	var response struct {
		Order NYOrder `json:"order"`
	}
	if e = s.DeliverySupplier.Call(ctx, "/order/"+*p.Supplier, "", nil, &response); e != nil {
		return e
	}
	v := response.Order
	normalize(&v)
	if v.ID != *p.Supplier || (v.Service != "" && v.Service != p.Service) || v.Requested == nil || *v.Requested != 1 {
		return errors.New("supplier order identity mismatch")
	}
	charged, e1 := usdMicros(v.Charged)
	refunded, e2 := usdMicros(v.Refunded)
	net, e3 := usdMicros(v.Net)
	// A rejection needs the supplier's own refund counters to agree. Partial or
	// missing evidence stays pending: the charge may still be reversed later.
	if v.Status == "refunded" && e1 == nil && e2 == nil && e3 == nil && net == 0 && refunded == charged && v.Delivered != nil && *v.Delivered == 0 {
		return s.reject(ctx, "delivery_supplier", p, "delivery_supplier-refunded", "供应商已确认采购退款，请通过原订单处理客户售后")
	}
	if v.Status != "success" {
		delay := min(float64(3600), max(float64(5), v.Retry))
		_, e = s.Store.DB.Exec(ctx, `UPDATE delivery_supplier_orders SET state='pending',next_check_at=$1,updated_at=$2 WHERE order_id=$3 AND lease_token=$4 AND lease_until>$2`, time.Now().UnixMilli()+int64(delay*1000), time.Now().UnixMilli(), p.ID, p.Lease)
		return e
	}
	// An upstream that reports only the status, the charge and the delivered
	// values proves completion with those alone: refunds and partial deliveries
	// would have to appear as a different status.
	if e1 != nil || e2 != nil || e3 != nil || charged <= 0 || charged > p.MaxUSD*10000 || refunded != 0 || net != charged || len(v.Values) != *v.Requested {
		return errors.New("supplier purchase amount or quantity mismatch")
	}
	if !v.Minimal && (v.Delivered == nil || *v.Delivered != 1 || v.Partial == nil || *v.Partial) {
		return errors.New("supplier purchase amount or quantity mismatch")
	}
	tx, e := s.Store.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	o, p, e := s.fence(ctx, tx, "delivery_supplier", p)
	if e != nil {
		return e
	}
	value := delivery_supplierDeliveryValue(v.Values[0])
	if !delivery.ValidValue(o.Delivery(), value) {
		return errors.New("invalid supplier delivery value")
	}
	_, e = tx.Exec(ctx, `UPDATE delivery_supplier_orders SET state='purchased',actual_usd_cents=$1,actual_cost_cents=$2,value=$3,updated_at=$4,actual_usd_micros=$6 WHERE order_id=$5`, (charged+5000)/10000, int64(math.Round(float64(charged)*p.FX/10000)), value, time.Now().UnixMilli(), p.ID, charged)
	if e != nil {
		return e
	}
	return tx.Commit(ctx)
}

// canboso buys one account from the buyer-key shop. The purchase response is
// the only place the credentials appear: the history endpoint reports status
// only, so an order that is not completed at purchase time can never be
// delivered automatically and is parked for an operator.
func (s *Service) canboso(ctx context.Context, p Purchase) error {
	if p.Product == nil {
		return errors.New("missing canboso product")
	}
	if p.Attempted == nil {
		products, e := s.Canboso.Products(ctx)
		if e != nil {
			return e
		}
		if e = validateCanboso(products, *p.Product, p.MaxUSD); e != nil {
			return e
		}
		if p, e = s.start(ctx, "canboso", p); e != nil {
			return e
		}
		order, e := s.Canboso.Purchase(ctx, *p.Product, p.Key)
		if errors.Is(e, ErrRejected) || errors.Is(e, ErrCanbosoBalance) {
			return s.reject(ctx, "canboso", p, "canboso-purchase-rejected", "供应商拒绝采购或钱包余额不足，请核对原订单售后")
		}
		if errors.Is(e, ErrCanbosoInventory) {
			return s.deferCanboso(ctx, p, 60000)
		}
		if e != nil {
			return e
		}
		if e = s.remember(ctx, "canboso", p, order.Code); e != nil {
			return e
		}
		return s.settleCanboso(ctx, p, order, true)
	}
	// A write may already have happened: only read from here on. Replaying the
	// same idempotency key is safe but the history read costs nothing.
	if p.Supplier == nil {
		return orders.Conflict()
	}
	order, e := s.Canboso.Order(ctx, *p.Supplier)
	if errors.Is(e, ErrRejected) {
		return s.deferCanboso(ctx, p, 60000)
	}
	if e != nil {
		return e
	}
	return s.settleCanboso(ctx, p, order, false)
}

// deferCanboso parks the purchase until the next reconciliation pass.
func (s *Service) deferCanboso(ctx context.Context, p Purchase, delay int64) error {
	_, e := s.Store.DB.Exec(ctx, `UPDATE canboso_orders SET state='pending',next_check_at=$1,updated_at=$2 WHERE order_id=$3 AND lease_token=$4 AND lease_until>$2`, time.Now().UnixMilli()+delay, time.Now().UnixMilli(), p.ID, p.Lease)
	return e
}

func (s *Service) settleCanboso(ctx context.Context, p Purchase, order CanbosoOrder, direct bool) error {
	switch order.Status {
	case "completed":
		if len(order.Delivery.Accounts) != 1 {
			if !direct {
				return s.deferCanboso(ctx, p, 3600000)
			}
			return errors.New("canboso delivery not ready")
		}
	case "paid", "pending", "processing":
		return s.deferCanboso(ctx, p, 600000)
	case "cancelled", "refunded", "failed":
		return s.reject(ctx, "canboso", p, "canboso-purchase-failed", "供应商确认采购失败或已退款，请通过原订单处理售后")
	default:
		return errors.New("unexpected canboso order status")
	}
	if order.Product != "" && *p.Product != order.Product {
		return errors.New("canboso product mismatch")
	}
	value, e := canbosoValue(order.Delivery.Accounts[0])
	if e != nil {
		return e
	}
	usd := p.MaxUSD
	if direct {
		if usd, e = canbosoUSDCents(order.Payment.Amount); e != nil || order.Payment.Currency != "USD" {
			return errors.New("invalid canboso charge")
		}
		if usd > p.MaxUSD {
			return errors.New("canboso charge above ceiling")
		}
	}
	tx, e := s.Store.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	o, p, e := s.fence(ctx, tx, "canboso", p)
	if e != nil {
		return e
	}
	if !delivery.ValidValue(o.Delivery(), value) {
		return errors.New("invalid supplier delivery value")
	}
	_, e = tx.Exec(ctx, `UPDATE canboso_orders SET state='purchased',actual_usd_cents=$1,actual_cost_cents=$2,delivery=$3,value=$4,updated_at=$5,actual_usd_micros=$6 WHERE order_id=$7`, usd, int64(math.Round(float64(usd)*p.FX)), raw(order.Delivery.Accounts), value, time.Now().UnixMilli(), usd*10000, p.ID)
	if e != nil {
		return e
	}
	return tx.Commit(ctx)
}
