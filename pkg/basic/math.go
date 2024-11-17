package basic

import "math"

// Sign 符号函数
func Sign(num float64) int {
	if num > 0 {
		return 1
	} else if num == 0 {
		return 0
	} else {
		return -1
	}
}

// Threshold 阈值函数
func Threshold(x float64) float64 {
	if math.Abs(x) >= 1 {
		return math.Abs(x)
	}
	return 1
}
