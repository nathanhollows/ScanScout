package templates

import "strconv"

func intToString(i int) string {
	return strconv.Itoa(i)
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func filter[T any](items []T, fn func(T) bool) []T {
	var filtered []T
	for _, item := range items {
		if fn(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func percentage(num, den int) int {
	if num == 0 || den == 0 {
		return 0
	}
	if num > den {
		return 100
	}
	return (num * 100) / den
}
