package util

import "github.com/shopspring/decimal"

// Decimal 封装 decimal 操作
type Decimal struct {
	decimal.Decimal
}

// NewDecimalFromString 从字符串创建 Decimal
func NewDecimalFromString(value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	return Decimal{d}, err
}

// NewDecimalFromFloat 从浮点数创建 Decimal
func NewDecimalFromFloat(value float64) Decimal {
	return Decimal{decimal.NewFromFloat(value)}
}

// Add 加法
func (d Decimal) Add(value Decimal) Decimal {
	return Decimal{d.Decimal.Add(value.Decimal)}
}

// Sub 减法
func (d Decimal) Sub(value Decimal) Decimal {
	return Decimal{d.Decimal.Sub(value.Decimal)}
}

// String 返回字符串表示
func (d Decimal) String() string {
	return d.Decimal.String()
}
