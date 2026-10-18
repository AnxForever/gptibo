package procurement

import (
	"encoding/json"
	"testing"
)

const documentedOrder = `{"success":true,"order":{"order_id":"ORD-1","service_id":"S_03","service":"Duolingo Super 12M","quantity":1,"amount":0.5,"status":"success","created_at":"2026-10-18 16:08:18","delivered_products":["https://www.duolingo.com/redeem?code=X"]}}`

func decodeDocumented(t *testing.T, raw string) NYOrder {
	t.Helper()
	var response struct {
		Order NYOrder `json:"order"`
	}
	if e := json.Unmarshal([]byte(raw), &response); e != nil {
		t.Fatal(e)
	}
	normalize(&response.Order)
	return response.Order
}

func TestDocumentedOrderShapeSettlesOnStatusAndCharge(t *testing.T) {
	v := decodeDocumented(t, documentedOrder)
	if !v.Minimal || v.Status != "success" || v.Requested == nil || *v.Requested != 1 {
		t.Fatalf("short fields not adopted: %+v", v)
	}
	charged, e1 := usdMicros(v.Charged)
	refunded, e2 := usdMicros(v.Refunded)
	net, e3 := usdMicros(v.Net)
	if e1 != nil || e2 != nil || e3 != nil || charged != 500000 || refunded != 0 || net != charged {
		t.Fatal("amounts not derived exactly", charged, refunded, net, e1, e2, e3)
	}
	if len(v.Values) != *v.Requested {
		t.Fatal("delivered values do not match the requested quantity")
	}
}

func TestDocumentedOrderShapeKeepsExtendedFieldsAuthoritative(t *testing.T) {
	raw := `{"order":{"order_id":"ORD-2","service_id":"S_03","requested_quantity":1,"order_status":"success","charged_amount":0.5,"refunded_amount":0.5,"net_charged_amount":0,"delivered_quantity":0,"partial_delivery":false,"delivered_products":[]}}`
	v := decodeDocumented(t, raw)
	if v.Minimal {
		t.Fatal("extended counters must never be replaced")
	}
	if refunded, _ := usdMicros(v.Refunded); refunded != 500000 {
		t.Fatal("reported refund was overwritten")
	}
	if net, _ := usdMicros(v.Net); net != 0 {
		t.Fatal("reported net charge was overwritten")
	}
}

func TestDocumentedOrderShapeDoesNotAssumeSuccess(t *testing.T) {
	for _, status := range []string{"pending", "refunded", "failed", "cancelled"} {
		raw := `{"order":{"order_id":"ORD-3","service_id":"S_03","quantity":1,"amount":0.5,"status":"` + status + `","delivered_products":[]}}`
		v := decodeDocumented(t, raw)
		if v.Minimal || v.Net != nil {
			t.Fatalf("derived a charge for %s", status)
		}
	}
}

func TestDocumentedOrderShapeRejectsUnsafeAmounts(t *testing.T) {
	for _, amount := range []string{`"0.5"`, `null`, `-0.5`, `0.5000001`} {
		raw := `{"order":{"order_id":"ORD-4","service_id":"S_03","quantity":1,"amount":` + amount + `,"status":"success","delivered_products":["X"]}}`
		v := decodeDocumented(t, raw)
		if _, e := usdMicros(v.Charged); e == nil && amount != `"0.5"` {
			t.Fatal("accepted a malformed amount", amount)
		}
	}
}
