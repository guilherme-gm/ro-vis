package utils

func ConvertIntSliceToInt32(src []int) []int32 {
	out := make([]int32, len(src))

	for i, v := range src {
		out[i] = int32(v)
	}

	return out
}
