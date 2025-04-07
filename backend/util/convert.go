package util

import (
	"strconv"
	"time"
)

func AtoiOrZero(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func TodayDate() string {
	return time.Now().Format("2006-01-02")
}
