package common

// SplitSlice 根据size分割原始切片并存储子切片
func SplitSlice[T any](slice []T, size int) [][]T {
	// 计算子切片的数量
	numSlices := len(slice) / size

	// 创建一个长度为 numSlices 的切片来存储子切片
	subSlices := make([][]T, numSlices)

	// 分割原始切片并存储子切片
	for i := 0; i < numSlices; i++ {
		start := i * size
		end := start + size
		subSlices[i] = slice[start:end]
	}

	return subSlices
}
