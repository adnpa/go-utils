package utils

import (
	"math/rand"
)

// RandIntN 从 [0, n) 随机获取一个整数
func RandIntN(n int) int {
	return rand.Intn(n)
}

// RandIntRange 从 [min, max) 中随机获取一个整数
func RandIntRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// RandFloat64 从 [0.0, 1.0] 中获取随机数
func RandFloat64() float64 {
	return rand.Float64()
}

// RandNormFloat64 从均值为0, 标准差为1的正态分布获取随机浮点数
func RandNormFloat64() float64 {
	return rand.NormFloat64()
}

// RandNormFloat64WithMeanStddev 从均值为mean, 标准差为stddev 的正态分布中获取随机浮点数
func RandNormFloat64WithMeanStddev(mean, stddev float64) float64 {
	return rand.NormFloat64()*stddev + mean
}

// Perm 返回半开区间 [0,n) 内整数伪随机排列
func Perm(n int) []int {
	return rand.Perm(n)
}

// Shuffle 打乱数组中元素的顺序
func Shuffle[T any](slice []T) {
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)                   // 生成随机索引
		slice[i], slice[j] = slice[j], slice[i] // 交换元素
	}
}
