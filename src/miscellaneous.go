package main

import "math"

const max_int = int(^uint(0) >> 1)

func min(a int, b int) int {
	if (a < b) {
		return a
	}
	return b
}

func max(a int, b int) int {
	if (a > b) {
		return a
	}
	return b
}

// Find the size of the step for a given variable with given number of splits.
func step_size(vari variable, splits int) float64 {
	return (vari.upper_boundary - vari.lower_boundary) / float64(splits - 1)
}

// Compute the maximum number of splits for each variable.
func possible_values (variables []variable, splits_cap int) []int {
	output := make([]int, len(variables))
	for i, vari := range variables {
		if vari.format == string_format {
			output[i] = len(vari.options)
		} else if vari.lower_boundary == vari.upper_boundary {
			output[i] = 1
		} else if vari.format == float_format {
			output[i] = variables[i].splits
		} else {
			max_splits := int(variables[i].upper_boundary - variables[i].lower_boundary + 1)
			if variables[i].splits != 0 {
				max_splits = min(max_splits, variables[i].splits)
			} else if splits_cap != 0 {
				max_splits = min(max_splits, splits_cap)
			}
			output[i] = max_splits
		}
		if output[i] == 0 {
			output[i] = splits_cap
		}
	}
	return output
}

type runs_update_function func (int, int) int
type expected_compute_function func (int, int) int
type splits_config struct {
	update_function runs_update_function
	compute_function expected_compute_function
}

func runs_divide (runs, split int) int {
	return runs / split
}

func runs_subtract (runs, split int) int {
	return runs - split
}

func expected_root (runs, remaining int) int {
	return int(math.Pow(float64(runs), 1.0 / float64(remaining)))
}

func expected_division(runs, remaining int) int {
	return runs / remaining
}

func get_multiplication_config () splits_config  {
	return splits_config{runs_divide, expected_root}
}	

func get_addition_config () splits_config  {
	return splits_config{runs_subtract, expected_division}
}	

// Find the number of splits for each variable in a way that the total number
// of calls as described by `config` is about `runs` and the number of splits
// does not exceed the maximum number of splits in each variable.
func find_splits (variables []variable, possibilities []int, runs int, config splits_config) []int {
	splits := make([]int, len(variables))
	for i := len(variables) - 1; i >= 0; i-- {
		max_splits := int(variables[i].upper_boundary - variables[i].lower_boundary + 1)
		if variables[i].format == string_format {
			splits[i] = len(variables[i].options)
		} else if variables[i].splits != 0 {
			splits[i] = min(variables[i].splits, max_splits)
		} else if variables[i].lower_boundary == variables[i].upper_boundary {
			splits[i] = 1
		} else {
			splits[i] = max(config.compute_function(runs, i + 1), 2)
			if i > 0 && possibilities[i-1] != 0 {
				splits[i] = max(splits[i], config.update_function(runs, possibilities[i-1]))
			}
			if variables[i].format == int_format {
				splits[i] = min(splits[i], max_splits)
			}
		}
		runs = config.update_function(runs, splits[i])
	}
	return splits
}

type prefix_function func(int, int) int

func product (a, b int) int {
	return a * b
}

func sum_no_negative(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return a + b
}

// Compute the generalized prefix sum array where summation is replaced by `prefix_function`.
func compute_prefix (values []int, function prefix_function) []int {
	output := make([]int, len(values))
	output[0] = values[0]
	for i := 1; i < len(values); i++ {
		output[i] = function(output[i-1], values[i])
	}
	return output
}
