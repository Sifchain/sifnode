package main

import (
	"fmt"
	"math"
)

func getData() (float64, float64, float64, float64, float64, float64, float64, float64, float64) {

	//-----------------------------------------
	// Real looking data
	// var Y, X, P float64 = 94031960239698233561402555, 984870852540, 12285414478722416018888701 // USDC pool 8 June 15:03
	// var yPrice, xPrice, yDecimals, xDecimals float64 = 0.01033, 1, 18, 6                       // Rowan price $0.01033 8 June 15:03

	// //var y, x float64 = 0, 100000000 // asymmetric $100 USDC
	// //var y, x float64 = 10000 * math.Pow(10, 18), 0 // asymmetric $103 rowan
	// //var y, x float64 = 10000 * math.Pow(10, 18), 10000 * math.Pow(10, 18) * X / Y // symmetric
	// var y, x float64 = 0, 10_000_000_000 // asymmetric $10,000 USDC

	//------------------------------------------
	// Very shallow pool

	var Y, X, P float64 = 50, 50, 50
	var yPrice, xPrice, yDecimals, xDecimals float64 = 1, 1, 1, 1
	var y, x float64 = 10, 0

	return Y, X, P, y, x, yPrice, xPrice, yDecimals, xDecimals

}

func main() {
	var f, r float64 = 0.003, 0.05

	Y, X, P, y, x, yPrice, xPrice, yDecimals, xDecimals := getData()

	val := calcValue(y, x, yPrice, xPrice, yDecimals, xDecimals)
	fmt.Printf("Starting value: $%f\n", val)

	YAdd, XAdd, PAdd, LPUnits := addLiquidity(Y, X, y, x, f, r, P)

	xRemoved, yRemoved := removeAllLiquidity(YAdd, XAdd, PAdd, LPUnits)

	val = calcValue(yRemoved, xRemoved, yPrice, xPrice, yDecimals, xDecimals)

	fmt.Printf("End value: $%f\n", val)

	if xRemoved < x {
		// We have less x than we used to have
		// How much y would we get if we had just swapped this difference in x for y

		xDiff := x - xRemoved

		ySwap := calculateSwap(xDiff, X, Y, f, r) + y

		fmt.Printf("Amount of y if just swapped: %f\n", ySwap)
		fmt.Printf("Amount of y we got from dipping %f\n", yRemoved)
		fmt.Printf("y swapped / y dipping: %f\n", ySwap/yRemoved)

	} else if yRemoved < y {
		// We have less y than we used to have
		// How much x would we get if we had just swapped this difference in y for x

		yDiff := y - yRemoved

		xSwap := calculateSwap(yDiff, Y, X, f, r) + x

		fmt.Printf("Amount of x if just swapped: %f\n", xSwap)
		fmt.Printf("Amount of x we got from dipping %f\n", xRemoved)
		fmt.Printf("x swapped / x dipping: %f\n", xSwap/xRemoved)

	} else if (xRemoved == x) && (yRemoved == y) {
		fmt.Println("No change in removed x and y. Presumably this was a symmetric add?")
	} else {
		fmt.Println("More x and more y - unless this is rounding we have a problem")
	}

}

func calcValue(y, x, yPrice, xPrice, yDecimals, xDecimals float64) float64 {
	return x*xPrice/math.Pow(10, xDecimals) + y*yPrice/math.Pow(10, yDecimals)
}

func removeAllLiquidity(Y, X, P, LPUnits float64) (float64, float64) {
	frac := LPUnits / P

	x := frac * X
	y := frac * Y

	return x, y
}

func addLiquidity(Y, X, y, x, f, r, P float64) (float64, float64, float64, float64) {
	sellX, s := calculateSwapAmount(Y, X, y, x, f, r)

	fmt.Printf("Swap amount: %f\n", s)
	fmt.Printf("Sell x: %v\n", sellX)

	var xCorrected, yCorrected float64

	if sellX {
		xCorrected = x - s
		yCorrected = y + calculateSwap(s, X, Y, f, r)
	} else {
		xCorrected = x + calculateSwap(s, Y, X, f, r)
		yCorrected = y - s
	}

	fmt.Printf("x after swap: %f\n", xCorrected)
	fmt.Printf("Pool ratio: %f\n", Y/X)
	fmt.Printf("Asset ratio: %f\n", yCorrected/xCorrected)
	fmt.Printf("Pool ratio / asset ratio: %f\n", Y/X/(yCorrected/xCorrected))

	LPUnits := calculateLPPoolUnits(P, X, Y, xCorrected, yCorrected)
	P = P + LPUnits

	return Y + y, X + x, P, LPUnits
}

func calculateSwapAmount(Y, X, y, x, f, r float64) (bool, float64) {
	if Y/X > y/x {
		return true, math.Abs(1 / (2 * (r + 1) * Y) * (math.Sqrt(math.Pow((-1*f*X*Y-r*x*Y+r*X*y+r*X*Y-x*Y+X*y+2*X*Y), 2)-4*(r*Y+Y)*(-1*r*x*X*Y+r*X*X*y-x*X*Y+X*X*y)) + f*X*Y + r*x*Y - r*X*y - r*X*Y + x*Y - X*y - 2*X*Y))
	} else {
		tmp := X
		X = Y
		Y = tmp
		tmp = x
		x = y
		y = tmp
		//return false, math.Abs((math.Sqrt(math.Pow((-1*f*r*X*Y-f*X*Y+r*X*Y+x*Y-X*y+2*X*Y), 2)-4*X*(x*Y*Y-X*y*Y)) + f*r*X*Y + f*X*Y - r*X*Y - x*Y + X*y - 2*X*Y) / (2 * X))
		return false, math.Abs(1 / (2 * (r + 1) * Y) * (math.Sqrt(math.Pow((-1*f*X*Y-r*x*Y+r*X*y+r*X*Y-x*Y+X*y+2*X*Y), 2)-4*(r*Y+Y)*(-1*r*x*X*Y+r*X*X*y-x*X*Y+X*X*y)) + f*X*Y + r*x*Y - r*X*y - r*X*Y + x*Y - X*y - 2*X*Y))
	}
}

func calculateSwap(x, X, Y, f, r float64) float64 {
	return x * Y * (1 - f) / ((x + X) * (1 + r))
}

func calculateLPPoolUnits(P, X, Y, x, y float64) float64 {
	return x / X * P
}
