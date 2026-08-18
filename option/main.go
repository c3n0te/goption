package main

import (
	"fmt"
	"log/slog"
)

func main() {
	// Configure simulation parameters
	params := OptionParams{
		S0:    100.0, // Stock Price
		K:     100.0, // Strike Price
		T:     1.0,   // 1 Year to Maturity
		R:     0.05,  // 5% Risk-free rate
		Sigma: 0.2,   // 20% Volatility
	}

	grid := GridParams{
		NS: 200, // Stock price steps
		NT: 100, // Time steps
	}

	price := CrankNicolson(params, grid)
	slog.Info(fmt.Sprintf("stock price = %.4f", price))
}
