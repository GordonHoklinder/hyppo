package main

import (
	"testing"
	"math"
)

func Test_min(t *testing.T) {
	assert(t, min(9, -5) == -5, "-5 should be less than 9.")
	assert(t, min(15156, 15157) == 15156, "15156 should be less than 15157.")
}

func Test_max(t *testing.T) {
	assert(t, max(37, -42) == 37, "-42 should be less than 37.")
	assert(t, max(0, 1516) == 1516, "1516 should be more than 0.")
}

func Test_step_size(t *testing.T) {
	vari := float_variable("", 2.5, 5.5, 7, false, 0.0)
	assert(t, math.Abs(step_size(vari, vari.splits) - 0.5) < 10e-7, "Step size is not correct")
}

func Test_possible_values(t *testing.T) {
	variables := []variable{
		string_variable("", []string{"prvni", "druhy", "treti"}, false, ""),
		float_variable("", 0.0, 0.0, 0, false, 0.0),
		float_variable("", 0.0, 1.0, 7, false, 0.0),
		float_variable("", 0.0, 8.0, 2, false, 0.0),
		float_variable("", 0.0, 1.0, 0, false, 0.0),
		int_variable("", 0.0, 8.0, 0, false, 0.0),
		int_variable("", 0.0, 8.0, 6, false, 0.0),
		int_variable("", 0.0, 8.0, 11, false, 0.0),
	}
	values := possible_values(variables)
	assert(t, len(values) == 8, "values have inccorect length.")
	assert(t, values[0] == 3, "Incorrect splits for string variables.")
	assert(t, values[1] == 1, "Incorrect splits for float variables with one option.")
	t.Logf("%d", values[2])
	assert(t, values[2] == 7, "Incorrect splits for float variables with splits.")
	assert(t, values[3] == 2, "Incorrect splits for float variables with splits.")
	assert(t, values[4] == 0, "Incorrect splits for float variables without splits.")
	assert(t, values[5] == 9, "Incorrect splits for int variables without splits.")
	assert(t, values[6] == 6, "Incorrect splits for int variables with splits.")
	assert(t, values[7] == 9, "Incorrect splits for int variables with splits and limited range.")
}

func Test_runs_divide(t *testing.T) {
	assert(t, runs_divide(7, 3) == 2, "Incorrect division in runs_divide.")
}

func Test_runs_subtract(t *testing.T) {
	assert(t, runs_subtract(7, 3) == 4, "Incorrect subtraction in runs_subtract.")
}

func Test_product(t *testing.T) {
	assert(t, product(2, 3) == 6, "Incorrect product is not correct.")
	assert(t, product(2, 0) == 0, "Product should be zero.")
}

func Test_sum_no_negative(t *testing.T) {
	assert(t, sum_no_negative(3, 5) == 8, "Incorrect sum.")
	assert(t, sum_no_negative(0, 5) == 0, "Zero at first element should imply zero result.")
	assert(t, sum_no_negative(3, 0) == 0, "Zero at second element should imply zero result.")
}

func Test_compute_prefix(t *testing.T) {
	prefix_sums := compute_prefix([]int{5, 3, 0, 7}, sum_no_negative)
	assert_slices_equal(t, prefix_sums, []int{5, 8, 0, 0})
	prefix_products := compute_prefix([]int{1, 5, 3, 0, 7}, product)
	assert_slices_equal(t, prefix_products, []int{1, 5, 15, 0, 0})
}
