package util

import (
	"strconv"
	"strings"
)

func ParseMilestones(raw string) map[int]string {
	result := make(map[int]string)
	pairs := strings.Split(raw, ",")
	for _, p := range pairs {
		split := strings.Split(p, ":")
		if len(split) == 2 {
			k, _ := strconv.Atoi(split[0])
			result[k] = split[1]
		}
	}
	return result
}
