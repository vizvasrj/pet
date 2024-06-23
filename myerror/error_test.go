package myerror

import "testing"

func someError() (int, error) {
	return 0, WrapError(New("some new error"), "")
}

func someErrorTwo() error {
	_, err := someError()
	return WrapError(err, "error from SOMEERROR_TWO_2")
}

func TestError(t *testing.T) {
	err := someErrorTwo()
	t.Log(err)
}
