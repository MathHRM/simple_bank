package util

const (
	USD = "USD"
	BRL = "BRL"
	KWZ = "KWZ"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, BRL, KWZ:
		return true
	}
	return false
}
