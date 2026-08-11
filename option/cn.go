package main

import (
	"gonum.org/v1/gonum/mat"
)

func CrankNicolson(dl, d, du []float64, n int) (*mat.VecDense, error) {
	A := mat.NewTridiag(n, dl, d, du)
	b := mat.NewVecDense(n, []float64{1, 0, 1})

	var x mat.VecDense
	err := A.SolveVecTo(&x, false, b) // solves tridiagonal system
	if err != nil {
		return &x, err
	}
	return &x, nil
}
