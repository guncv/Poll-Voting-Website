package util

import "strconv"

func AtoiOrZero(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
