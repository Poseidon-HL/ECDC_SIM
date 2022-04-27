package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(begin, end int) int {
	return begin + rand.Intn(end-begin+1)
}

func GenerateListSample(numRange int, num int) []int {
	var originList []int
	for i := 0; i < numRange; i++ {
		originList = append(originList, i)
	}
	return Sample(originList, num)
}

// Sample 从列表中随机取出num个元素
func Sample(sample []int, num int) []int {
	var result []int
	for num != 0 {
		idx := rand.Intn(len(sample))
		result = append(result, sample[idx])
		num--
		sample = append(sample[:idx], sample[idx+1:]...)
	}
	return result
}
