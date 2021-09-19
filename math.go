package cpt

import (
	"math"
)

func Add(a, b int) int {
	return a + b
}

func Sub(a, b int) int {
	return a - b
}

func Mul(a, b int) int {
	return a * b
}

func Div(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

func Pow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}
