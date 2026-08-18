package main

import (
	"log/slog"
	"math"

	"gonum.org/v1/gonum/mat"
)

// OptionParams stores the financial variables for the Black-Scholes PDE
type OptionParams struct {
	S0    float64 // Current stock price
	K     float64 // Strike price
	T     float64 // Time to maturity (years)
	R     float64 // Risk-free interest rate
	Sigma float64 // Volatility
}

// GridParams defines the resolution of the finite difference mesh
type GridParams struct {
	NS int // Number of stock price steps
	NT int // Number of time steps
}

// CrankNicolson prices a European Call Option using the Crank-Nicolson method
func CrankNicolson(op OptionParams, gp GridParams) float64 {
	// Upper boundary for stock price (typically 3 to 4 times the strike)
	sMax := 4.0 * op.K
	dS := sMax / float64(gp.NS)
	dT := op.T / float64(gp.NT)

	// Initialize vectors to hold option values
	// V represents the grid at the current time step (moving backwards)
	V := make([]float64, gp.NS+1)

	// 1. Terminal Condition: Populate payoff at maturity (T)
	for i := range gp.NS {
		S := float64(i) * dS
		V[i] = math.Max(S-op.K, 0.0) // Call option payoff
	}

	// Pre-allocate Thomas Algorithm arrays for tridiagonal system: a*V_{i-1} + b*V_i + c*V_{i+1} = d
	a := make([]float64, gp.NS)
	b := make([]float64, gp.NS+1)
	c := make([]float64, gp.NS)
	d := make([]float64, gp.NS+1)
	var Vnew *mat.VecDense

	// 2. Time-stepping loop (Backwards from maturity to t=0)
	for j := gp.NT - 1; j >= 0; j-- {

		// 3. Construct the Tridiagonal System for internal nodes (i = 1 to NS-1)
		for i := 1; i < gp.NS; i++ {
			//S := float64(i) * dS

			// Continuous Black-Scholes coefficients
			sigmaSq := op.Sigma * op.Sigma
			alpha := 0.25 * dT * (sigmaSq*float64(i*i) - op.R*float64(i))
			beta := -0.5 * dT * (sigmaSq*float64(i*i) + op.R)
			gamma := 0.25 * dT * (sigmaSq*float64(i*i) + op.R*float64(i))

			// LHS Matrix components (implicit part)
			a[i] = -alpha
			b[i] = 1.0 - beta
			c[i] = -gamma

			// RHS Vector components (explicit part + boundaries)
			d[i] = alpha*V[i-1] + (1.0+beta)*V[i] + gamma*V[i+1]
		}

		// 4. Set Boundary Conditions
		// Node 0: S = 0 -> Option Value = 0
		b[0] = 1.0
		d[0] = 0.0

		// Node NS: S = S_max -> Option Value ~ S_max - K * e^(-r * t_remaining)
		tRemaining := op.T - (float64(j) * dT)
		b[gp.NS] = 1.0
		d[gp.NS] = sMax - op.K*math.Exp(-op.R*tRemaining)

		// 5. Solve the system using the Thomas Algorithm (Tridiagonal Matrix Solver)
		A := mat.NewTridiag(gp.NS+1, a, b, c)
		B := mat.NewVecDense(gp.NS+1, d)
		var err error
		Vnew, err = SolveTridiag(A, B, gp.NS+1)
		if err != nil {
			slog.Error("Failed to solve tridiagonal system", "error", err)
		}
	}

	return Interpolate(Vnew, dS, op.S0)
}

// interpolate maps the discrete grid back to the exact initial asset price S0
func Interpolate(V *mat.VecDense, dS, s0 float64) float64 {
	index := s0 / dS
	iLow := int(math.Floor(index))
	iHigh := int(math.Ceil(index))

	if iLow == iHigh {
		return V.At(iLow, 0)
	}

	// Linear interpolation formula
	weightHigh := index - float64(iLow)
	weightLow := 1.0 - weightHigh

	return weightLow*V.At(iLow, 0) + weightHigh*V.At(iHigh, 0)
}

func SolveTridiag(A *mat.Tridiag, b *mat.VecDense, n int) (*mat.VecDense, error) {
	var x mat.VecDense
	err := A.SolveVecTo(&x, false, b) // solves tridiagonal system
	if err != nil {
		return &x, err
	}

	return &x, nil
}
