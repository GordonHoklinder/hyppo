package main

import (
	"math/rand"
	"math"
	"strconv"
)

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

type property struct {
	as_string string
	as_float float64
	as_int int
}

func property_from_string (value string) property {
	return property{value, 0.0, 0}
}

func property_from_int (value int) property {
	return property{"", 0.0, value}
}

func property_from_float (value float64) property {
	return property{"", value, 0}
}

type individual []property

type evaluated_individual struct {
	score float64
	data individual
}

func get_random_individual (variables []variable) evaluated_individual {
	newborn := make(individual, len(variables))
	for i, vari := range variables {
		if vari.format == string_format {
			newborn[i] = property_from_string(vari.options[rand.Intn(len(vari.options))])
		} else if vari.format == int_format {
			lower := int(variables[i].lower_boundary)
			upper := int(variables[i].upper_boundary)
			newborn[i] = property_from_int(rand.Intn(upper - lower + 1) + lower)
		} else {
			lower := variables[i].lower_boundary
			upper := variables[i].upper_boundary
			newborn[i] = property_from_float(lower + rand.Float64() * (upper - lower))
		}
	}
	return evaluated_individual{0.0, newborn}
}

func (this individual) to_string_slice (variables []variable) []string {
	output := make([]string, len(variables)) 
	for i, prop := range this {
		if variables[i].format == string_format {
			output[i] = prop.as_string
		} else if variables[i].format == int_format {
			output[i] = strconv.Itoa(prop.as_int)
		} else {
			output[i] = strconv.FormatFloat(prop.as_float, 'f', 12, 64)
		}
	}
	return output
}

func (this evaluated_individual) evaluate_individual (variables []variable, script_communicator communicator) evaluated_individual {
	this.score, _ = script_communicator.run_arguments(variables, this.data.to_string_slice(variables))
	return this
}

