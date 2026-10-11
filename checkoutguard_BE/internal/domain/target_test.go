package domain

import (
	"testing"
)

func TestEnvironmentURLBoundaries(t *testing.T) {
	for _, c := range []struct {
		environment, url string
		valid            bool
	}{
		{"sandbox", "http://127.0.0.1:8080/sandbox", true},
		{"shopee", "https://shopee.co.id/product/123/456", true},
		{"shopee", "https://shopee.co.id/Tote-i.123.456?sp_atk=tracking", true},
		{"tokopedia", "https://www.tokopedia.com/shop/tote", true},
		{"sandbox", "https://shopee.co.id/product/123/456", false},
		{"shopee", "https://www.tokopedia.com/shop/tote", false},
		{"tokopedia", "https://www.tokopedia.com/checkout/confirm", false},
		{"shopee", "https://shopee.co.id.evil.test/Tote-i.123.456", false},
		{"shopee", "https://shopee.co.id@evil.test/Tote-i.123.456", false},
		{"shopee", "https://user:password@shopee.co.id/Tote-i.123.456", false},
		{"shopee", "http://shopee.co.id/Tote-i.123.456", false},
		{"tokopedia", "https://www.tokopedia.com:444/shop/tote", false},
		{"tokopedia", "https://www.tokopedia.com/shop/%2e%2e", false},
		{"other", "https://example.com/shop/tote", false},
	} {
		if got := ValidateTarget(c.environment, c.url, "http://127.0.0.1:8080") == nil; got != c.valid {
			t.Errorf("%s %s accepted=%t", c.environment, c.url, got)
		}
	}
}
