package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(cr string) bool {
	switch cr {
	case USD, EUR, CAD:
		return true
	}

	return false
}
