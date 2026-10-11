package domain

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

func ValidateAction(a Action) error {
	if a.Reason == "" || len(a.Reason) > 300 {
		return errors.New("reason must be a short decision summary")
	}
	switch a.Tool {
	case "click":
		if a.Selector == "" || a.Value != "" {
			return errors.New("click requires selector only")
		}
	case "fill":
		if a.Selector == "" || a.Value == "" || len(a.Value) > 100 {
			return errors.New("fill requires selector and a short value")
		}
	case "observe", "verify_total", "order_status", "finish":
		if a.Selector != "" || a.Value != "" {
			return errors.New("this tool takes no selector or value")
		}
	default:
		return errors.New("unknown tool")
	}
	return nil
}

var shopeeProduct = regexp.MustCompile(`(?:-i\.[0-9]+\.[0-9]+$|^/product/[0-9]+/[0-9]+/?$)`)

func MarketHome(environment string) string {
	if environment == "shopee" {
		return "https://shopee.co.id/"
	}
	if environment == "tokopedia" {
		return "https://www.tokopedia.com/"
	}
	return ""
}

func ValidateTarget(environment, raw, base string) error {
	if environment == "sandbox" {
		if raw != "" && strings.TrimRight(raw, "/") != base+"/sandbox" {
			return errors.New("sandbox mode accepts only the local synthetic sandbox URL")
		}
		return nil
	}
	if MarketHome(environment) == "" {
		return errors.New("choose sandbox, shopee, or tokopedia")
	}
	u, err := url.Parse(raw)
	if err != nil || !MarketPageAllowed(environment, raw) || u.RawPath != "" || u.Fragment != "" {
		return errors.New("enter a full HTTPS product URL for the selected marketplace")
	}
	if environment == "shopee" && !shopeeProduct.MatchString(u.Path) {
		return errors.New("use a Shopee product URL ending in -i.shopID.itemID or /product/shopID/itemID")
	}
	if environment == "tokopedia" {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		reserved := map[string]bool{"login": true, "register": true, "cart": true, "checkout": true, "order": true, "payment": true, "discovery": true, "search": true, "help": true, "user": true, "settings": true, "account": true, "api": true}
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" || reserved[strings.ToLower(parts[0])] || parts[0] == "." || parts[0] == ".." || parts[1] == "." || parts[1] == ".." {
			return errors.New("use a Tokopedia product URL: https://www.tokopedia.com/shop/product")
		}
	}
	return nil
}

func MarketPageAllowed(environment, raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if environment == "shopee" {
		return host == "shopee.co.id" || host == "www.shopee.co.id"
	}
	if environment == "tokopedia" {
		return host == "tokopedia.com" || host == "www.tokopedia.com"
	}
	return false
}
