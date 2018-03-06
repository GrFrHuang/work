package models

import (
	"testing"
	"math/big"
	"fmt"
)

func TestFloat(t *testing.T) {
	fsfs := float64(1234567.8)
	f_rat := big.NewRat(1, 1)
	f_rat.SetFloat64(fsfs)
	fmt.Println("test : ", f_rat.FloatString(1))
}