package config

import (
	"errors"
	"gptibo/storefront/internal/procurement"
	"os"
	"strconv"
	"strings"
)

func LoadProcurement() (*procurement.Service, error) {
	s := &procurement.Service{}
	var e error
	s.EmailEnabled, e = flag("EMAIL_PRODUCTS_ENABLED")
	if e != nil {
		return nil, e
	}
	s.GmailEnabled, e = flag("ACCOUNT_SUPPLIER_CHECKOUT_ENABLED")
	if e != nil {
		return nil, e
	}
	s.NYEnabled, e = flag("DELIVERY_SUPPLIER_CHECKOUT_ENABLED")
	if e != nil {
		return nil, e
	}
	if os.Getenv("ORDER_TESTING_ENABLED") == "true" {
		s.EmailEnabled = false
	}
	enabled, e := flag("ACCOUNT_SUPPLIER_ENABLED")
	if e != nil {
		return nil, e
	}
	if enabled {
		secret, e := Secret("ACCOUNT_SUPPLIER_APP_SECRET")
		if e != nil {
			return nil, e
		}
		s.AccountSupplier, e = procurement.NewAccountSupplier(os.Getenv("ACCOUNT_SUPPLIER_API_BASE"), os.Getenv("ACCOUNT_SUPPLIER_APP_ID"), secret)
		if e != nil {
			return nil, e
		}
		s.GoodsID, _ = strconv.ParseInt(os.Getenv("ACCOUNT_SUPPLIER_GOODS_ID"), 10, 64)
		s.MaxCents, _ = strconv.ParseInt(os.Getenv("ACCOUNT_SUPPLIER_MAX_COST_CENTS"), 10, 64)
		if s.GmailEnabled && (s.GoodsID < 1 || s.MaxCents < 1 || s.MaxCents > 10000000) {
			return nil, errors.New("invalid Gmail goods or cost ceiling")
		}
	}
	enabled, e = flag("DELIVERY_SUPPLIER_ENABLED")
	if e != nil {
		return nil, e
	}
	if enabled {
		key, e := Secret("DELIVERY_SUPPLIER_API_KEY")
		if e != nil {
			return nil, e
		}
		base := os.Getenv("DELIVERY_SUPPLIER_API_BASE")
		if base == "" {
			return nil, errors.New("DELIVERY_SUPPLIER_API_BASE is required when DELIVERY_SUPPLIER_ENABLED is true")
		}
		s.DeliverySupplier, e = procurement.NewDeliverySupplier(base, key)
		if e != nil {
			return nil, e
		}
		s.FX, e = strconv.ParseFloat(os.Getenv("DELIVERY_SUPPLIER_CNY_PER_USD"), 64)
		if e != nil || !(s.FX > 0 && s.FX <= 100) {
			return nil, errors.New("invalid DELIVERY_SUPPLIER_CNY_PER_USD")
		}
	}
	s.CanbosoEnabled, e = flag("CANBOSO_CHECKOUT_ENABLED")
	if e != nil {
		return nil, e
	}
	enabled, e = flag("CANBOSO_ENABLED")
	if e != nil {
		return nil, e
	}
	if enabled {
		key, e := Secret("CANBOSO_KEY")
		if e != nil {
			return nil, e
		}
		s.Canboso, e = procurement.NewCanbosoSupplier(os.Getenv("CANBOSO_API_BASE"), key)
		if e != nil {
			return nil, e
		}
		s.CanbosoKeys, e = canbosoKeys()
		if e != nil {
			return nil, e
		}
		s.CanbosoFX, e = strconv.ParseFloat(os.Getenv("CANBOSO_CNY_PER_USD"), 64)
		if e != nil || !(s.CanbosoFX > 0 && s.CanbosoFX <= 100) {
			return nil, errors.New("invalid CANBOSO_CNY_PER_USD")
		}
	}
	s.Prepare()
	return s, nil
}

// canbosoKeys loads the named buyer keys from CANBOSO_KEY_FILES
// ("key1=/run/integrations/canboso-key-1,key4=/run/integrations/canboso-key-4").
// Each entry is a file next to the default CANBOSO_KEY_FILE, so secrets never
// travel in the environment.
func canbosoKeys() (map[string]*procurement.CanbosoSupplier, error) {
	spec := strings.TrimSpace(os.Getenv("CANBOSO_KEY_FILES"))
	if spec == "" {
		return nil, nil
	}
	accounts, e := procurement.ParseCanbosoKeyFiles(spec)
	if e != nil {
		return nil, e
	}
	out := make(map[string]*procurement.CanbosoSupplier, len(accounts))
	for _, a := range accounts {
		b, e := os.ReadFile(a.Path)
		if e != nil {
			return nil, errors.New("cannot read canboso key file for " + a.Name)
		}
		client, e := procurement.NewCanbosoSupplier(os.Getenv("CANBOSO_API_BASE"), strings.TrimRight(string(b), "\r\n"))
		if e != nil {
			return nil, e
		}
		out[a.Name] = client
	}
	return out, nil
}
