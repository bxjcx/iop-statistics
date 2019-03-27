package util

import (
	"strconv"
	"time"
)

func GetDate() int {
	n, err := strconv.Atoi(time.Now().Format("20060102"))
	if err != nil {
		panic(err) // 不废话，直接panic
	}
	return n
}
