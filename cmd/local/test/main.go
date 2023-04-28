package main

import "errors"

type s struct {
}

func (s s) Error() string {
	return "sdfsdf"
}

func main() {
	var err error
	if errors.Is(err, s{}) {

	}
}
