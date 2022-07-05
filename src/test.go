package main

import "testing"

func assert (t *testing.T, statement bool, message string) {
	if !statement {
		t.Log(message)
		t.Fail()
	}
}

func assert_slices_equal(t *testing.T, a, b []int) {
	assert(t, len(a) == len(b), "Arrays have different length.")
	for i, a_i := range a {
		t.Logf("At position %d the values are %d and %d", i, a_i, b[i])
		assert(t, a_i == b[i], "The arrays differ.")
	}
}

