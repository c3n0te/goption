package main

import "testing"

func TestCrankNicolson(t *testing.T) {
	tol := 5e-4
	n := 3
	dl := []float64{-1, -1} // sub-diagonal
	d := []float64{2, 2, 2} // main diagonal
	du := []float64{-1, -1} // super-diagonal
	x, err := CrankNicolson(dl, d, du, n)
	if err != nil {
		t.Errorf("Failed to solve tridiagonal system: %v", err)
	}

	if x.At(0, 0)-1.0 >= tol || x.At(1, 0)-1.0 >= tol || x.At(2, 0)-1.0 >= tol {
		t.Errorf("Incorrect solution to tridiagonal system")
	}

}
