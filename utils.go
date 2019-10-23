package utils

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
	}
}

func BodyForm(args []string) string {
	var s string
	if len(args) < 2 || args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func BodyForm2(args []string) string {
	var s string
	if len(args) < 3 || args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func BodyForm3(args []string) int {
	var s string
	if len(args) < 2 || args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}

	n, err := strconv.Atoi(s)
	FailOnError(err, "Failed to convert arg to integer")

	return n
}

func SeverityForm(args []string) string {
	var s string
	if len(args) < 2 || args[1] == "" {
		s = "info"
	} else {
		s = args[1]
	}
	return s
}

func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes = append(bytes, byte(RandInt(65, 90)))
	}

	return string(bytes)
}
