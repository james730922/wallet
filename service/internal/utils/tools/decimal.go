package tools

import "github.com/shopspring/decimal"

// DecimalToFloat64 converts an exact internal amount only at an API boundary
// whose protobuf contract still uses double.
func DecimalToFloat64(value decimal.Decimal) float64 {
	result, _ := value.Float64()
	return result
}

func DecimalPtr(value decimal.Decimal) *decimal.Decimal {
	return &value
}
