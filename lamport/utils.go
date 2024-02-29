package main

type MyError struct {
	message string
}

// custom error
func (err MyError) Error() string {
	return err.message
}
