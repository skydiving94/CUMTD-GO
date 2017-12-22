package main

import (
	"time"
	"strings"
	"strconv"
)

func TimeToSecond(time time.Time) (int) {
	return time.Hour() * 3600 + time.Minute() * 60 + time.Second()
}

func StringToSecond(str string) (int) {
	time := strings.Split(str, ":")
	hour, _ := strconv.Atoi(time[0])
	min, _ := strconv.Atoi(time[1])
	sec, _ := strconv.Atoi(time[2])
	return hour * 3600 + min * 60 + sec
}
