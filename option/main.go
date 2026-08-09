package main

import (
	"fmt"
	"log/slog"
	"time"

	"gonum.org/v1/gonum/mat"
)

func main() {
	n := 3
	a := []float64{-1, -1}  // sub-diagonal
	b := []float64{2, 2, 2} // main diagonal
	c := []float64{-1, -1}  // super-diagonal

	A := mat.NewTridiag(n, a, b, c)
	d := mat.NewVecDense(n, []float64{1, 0, 1})

	var x mat.VecDense
	start := time.Now()
	A.SolveVecTo(&x, false, d) // solves tridiagonal system
	elapsed := time.Since(start)

	var X mat.VecDense
	start1 := time.Now()
	X.SolveVec(A, d) // solves linear least squares
	elapsed1 := time.Since(start1)

	slog.Info(fmt.Sprintf("x = \n%.4f", mat.Formatted(&x, mat.Prefix(""))))
	slog.Info(fmt.Sprintf("X = \n%.4f", mat.Formatted(&X, mat.Prefix(""))))
	slog.Info(fmt.Sprintf("tridiag time: %v; least squares time: %v", elapsed.Nanoseconds(), elapsed1.Nanoseconds()))
}
