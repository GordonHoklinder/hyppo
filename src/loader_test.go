package main

import (
	"testing"
)

func assert (t *testing.T, statement bool, message string) {
	if !statement {
		t.Log(message)
		t.Fail()
	}
}

func assert_slices_equal(t *testing.T, a, b []int) {
	assert(t, len(a) == len(b), "Arrays have different length.")
	for i, a_i := range a {
		t.Log(a_i)
		t.Log(b[i])
		assert(t, a_i == b[i], "The arrays differ.")
	}
}

func Test_load_variables(t *testing.T) {
	variables, _ := load_variables("../variables.yaml")
	assert(t, len(variables) == 3, "Variables do not have correct length.")
	assert(t, variables[1].name == "learning_rate", "Incorrect variable name.")
	assert(t, variables[2].format == int_format, "Incorrect variable format.")
	assert(t, len(variables[0].options) == 2, "Incorrect length of options.")
	assert(t, len(variables[1].options) == 0, "Options should be by default empty.")
	assert(t, variables[0].has_default_value, "has_default_value should be true for set value.")
	assert(t, !variables[1].has_default_value, "has_default_value should be false by default.")
	assert(t, variables[0].default_value.as_string == "Adam", "has_default_value should be true for set value.")
	assert(t, variables[1].lower_boundary == 0.001, "Incorrect lower boundary.")
	assert(t, variables[1].upper_boundary == 1.0, "Incorrect lower boundary.")
	assert(t, variables[2].lower_boundary == 10.0, "Incorrect converted lower boundary.")
	assert(t, variables[1].splits == 0, "Splits should by default be 0.")
	assert(t, variables[2].splits == 10, "Incorrect splits.")
}

