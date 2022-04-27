package main

import (
	"math"
	"strconv"
)

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

type grid_optimizer struct {
	script_communicator communicator
	runs int
}

func (this grid_optimizer) find_optimal_hyperparameters(variables []variable) {
	possibilities := make([]int, len(variables))
	possibilities[0] = 1
	for i := 1; i < len(variables); i++ {
		if possibilities[i - 1] == 0 || variables[i].format == float_format {
			possibilities[i] = 0
		} else if variables[i].format == int_format {
			max_splits := int(variables[i].upper_boundary - variables[i].lower_boundary + 1)
			if variables[i].splits != 0 {
				max_splits = min(max_splits, variables[i].splits)
			}
			possibilities[i] = possibilities[i-1] * max_splits
		} else {
			possibilities[i] = possibilities[i-1] * len(variables[i].options)
		}
	}
	
	runs_left := this.runs
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
			splits[i] = max(int(math.Pow(float64(runs_left), 1.0 / float64(i + 1))), 2)
			if possibilities[i] != 0 {
				splits[i] = max(splits[i], runs_left / possibilities[i])
			}
			if variables[i].format == int_format {
				splits[i] = min(splits[i], max_splits)
			}
		}
		runs_left /= splits[i];
	}

	configurations_ran := make([]int, len(variables))
	changed_index := len(variables) - 1
	for {
		variables_used := make([]string, len(variables))
		for i, configuration := range configurations_ran {
			if variables[i].format == string_format {
				variables_used[i] = variables[i].options[configuration]
			} else if splits[i] == 1 {
				variables_used[i] = strconv.FormatFloat(variables[i].lower_boundary, 'f', 12, 64) 
			} else {
				value := variables[i].lower_boundary + (variables[i].upper_boundary - variables[i].lower_boundary) / float64(splits[i] - 1) * float64(configuration)
				if variables[i].format == int_format {
					variables_used[i] = strconv.Itoa(int(value))
				} else {
					variables_used[i] = strconv.FormatFloat(value, 'f', 12, 64) 
				}

			}
		}
		this.script_communicator.run_arguments(variables, variables_used)
		for i := changed_index; i >= 0; i-- {
			if configurations_ran[changed_index] == splits[changed_index] - 1 {
				configurations_ran[changed_index] = 0
				changed_index--
			} else {
				configurations_ran[changed_index]++
				changed_index = len(variables) - 1
				break
			}
		}
		if changed_index < 0 {
			break
		}
	}	
}
