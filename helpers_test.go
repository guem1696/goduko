package main

import (
	"testing"
)

func TestJoin(t *testing.T) {
	data := []string{
		"1",
		"2",
		"3",
	}

	res := Join(data, ";")

	if !Equals(res, []string{"1", ";", "2", ";", "3"}) {
		t.Fail()
	}
}

func Equals[T comparable](v1, v2 []T) bool {
	if len(v1) != len(v2) {
		return false
	}

	for idx := range v1 {
		if v1[idx] != v2[idx] {
			return false
		}
	}

	return true
}
