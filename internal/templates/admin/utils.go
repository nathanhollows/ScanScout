package templates

import "strconv"

func intToString(i int) string {
	return strconv.Itoa(i)
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
