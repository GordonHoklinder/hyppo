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

func step_size(vari variable, splits int) float64 {
	return (vari.upper_boundary - vari.lower_boundary) / float64(splits - 1)
}

func possible_values (variables []variable) []int {
	output := make([]int, len(variables))
	for i, vari := range variables {
		if vari.format == string_format {
			output[i] = len(vari.options)
		} else if vari.lower_boundary == vari.upper_boundary {
			output[i] = 1
		} else if vari.format == float_format {
			output[i] = 0
		} else {
			max_splits := int(variables[i].upper_boundary - variables[i].lower_boundary + 1)
			if variables[i].splits != 0 {
				max_splits = min(max_splits, variables[i].splits)
			}
			output[i] = max_splits
		}
	}
	return output
}

type runs_update_function func (int, int) int

func runs_divide (runs, split int) int {
	return runs / split
}

func runs_subtract (runs, split int) int {
	return runs - split
}

func find_splits (variables []variable, possibilities []int, runs int, splits_cap int, update_function runs_update_function) []int {
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
			splits[i] = max(int(math.Pow(float64(runs), 1.0 / float64(i + 1))), 2)
			if possibilities[i] != 0 {
				splits[i] = max(splits[i], runs / possibilities[i])
			}
			if splits_cap != 0 && variables[i].splits == 0 {
				splits[i] = min(splits[i], splits_cap)
			}
			if variables[i].format == int_format {
				splits[i] = min(splits[i], max_splits)
			}
		}
		runs = update_function(runs, splits[i])
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

func compute_prefix (values []int, function prefix_function, initial int) []int {
	output := make([]int, len(values))
	output[0] = initial
	for i := 1; i < len(values); i++ {
		output[i] = function(output[i-1], values[i])
	}
	return output
}
