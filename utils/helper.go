package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func MapToString(m map[string]int) string {
	var s string
	for k, v := range m {
		if s == "" {
			s = fmt.Sprintf("%s,%d", k, v)
		} else {
			s = fmt.Sprintf("%s/%s,%d", s, k, v)
		}
	}
	return s
}

func StringToMap(s string) map[string]int {
	m := make(map[string]int)
	if s == "" {
		return m
	}
	arr := strings.Split(s, "/")
	for _, val := range arr {
		keyval := strings.Split(val, ",")
		v, _ := strconv.Atoi(keyval[1])
		m[keyval[0]] = v
	}
	return m
}
