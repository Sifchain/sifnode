package main

import (
	"fmt"
	"math"
)

func main() {
	var Y, X, y, x, f, r float64 = 10000, 10000, 9, 1, 0.003, 0.05

	if Y/X > y/x {
		fmt.Println("more x in the pool than add")
		solution_for_s := math.Abs(1 / (2 * (r + 1) * Y) * (math.Sqrt(math.Pow((-1*f*X*Y-r*x*Y+r*X*y+r*X*Y-x*Y+X*y+2*X*Y), 2)-4*(r*Y+Y)*(-1*r*x*X*Y+r*X*X*y-x*X*Y+X*X*y)) + f*X*Y + r*x*Y - r*X*y - r*X*Y + x*Y - X*y - 2*X*Y))
		fmt.Println(solution_for_s)

	} else {
		fmt.Println("more x in the pool than add")
		solution_for_s := math.Abs((math.Sqrt(math.Pow((-1*f*r*X*Y-f*X*Y+r*X*Y+x*Y-X*y+2*X*Y), 2)-4*X*(x*Y*Y-X*y*Y)) + f*r*X*Y + f*X*Y - r*X*Y - x*Y + X*y - 2*X*Y) / (2 * X))
		fmt.Println(solution_for_s)
	}
}
