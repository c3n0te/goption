package main

import (
	"fmt"
	"log/slog"

	"gonum.org/v1/gonum/mat"
)

func main() {
	n := 3
	dl := []float64{-1, -1} // sub-diagonal
	d := []float64{2, 2, 2} // main diagonal
	du := []float64{-1, -1} // super-diagonal
	x, err := CrankNicolson(dl, d, du, n)
	if err != nil {
		slog.Error("Failed to solve tridiagonal system", "error", err)
	}

	slog.Info(fmt.Sprintf("x = \n%.4f", mat.Formatted(x, mat.Prefix(""))))
}
