package main

import "strconv"

type MyError struct {
	message string
}

// custom error
func (err MyError) Error() string {
	return err.message
}

func parseIntOrPanic(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return num
}
