package models

import "fmt"

type CurrencyType int32

func ToCurrencyType(f float32) CurrencyType {
	return CurrencyType((f * 100) + 0.5)
}

func (c CurrencyType) Float32() float32 {
	x := float32(c)
	return x / 100
}

func (c CurrencyType) String() string {
	x := float32(c)
	x = x / 100
	return fmt.Sprintf("%.2f", x)
}
