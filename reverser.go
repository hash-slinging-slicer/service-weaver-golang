package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type Reverser interface {
	Reverse(context.Context, string) (string, error)
	Penjumlahan(context.Context, float32, float32) (float32, error)
}

type reverser struct {
	weaver.Implements[Reverser]
}

func (r *reverser) Reverse(_ context.Context, s string) (string, error) {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	return string(runes), nil
}

func (r *reverser) Penjumlahan(_ context.Context, angka1 float32, angka2 float32) (float32, error) {
	total := angka1 + angka2
	return total, nil
}
