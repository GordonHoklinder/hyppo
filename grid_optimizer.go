package main

import (
	"math"
	"strconv"
)

type iterated_grid_optimizer struct {
	script_communicator communicator
	runs int
	iterations int
}

func (this iterated_grid_optimizer) get_params (variables_count int) (int, int) {
	runs_per_variable := 1 + int(math.Pow(float64(this.runs), 1.0/ float64(variables_count)))
	iteration_evaluations := 9
	if this.iterations != 0 {
		return this.iterations, iteration_evaluations
	}
	iterations := runs_per_variable / iteration_evaluations + 1
	return iterations, iteration_evaluations	
}

type grid_optimizer struct {
	script_communicator communicator
	runs int
}

func step_size(vari variable, splits int) float64 {
	return (vari.upper_boundary - vari.lower_boundary) / float64(splits - 1)
}

func grid_iteration(variables []variable, script_communicator communicator, runs int, splits_cap int) []variable {
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
	
	runs_left := runs
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
			if splits_cap != 0 && variables[i].splits == 0 {
				splits[i] = min(splits[i], splits_cap)
			}
			if variables[i].format == int_format {
				splits[i] = min(splits[i], max_splits)
			}
		}
		runs_left /= splits[i];
	}

	configurations_ran := make([]int, len(variables))
	changed_index := len(variables) - 1
	var best_variables []string
	best_score := math.Inf(-1)
	for {
		variables_used := make([]string, len(variables))
		for i, configuration := range configurations_ran {
			if variables[i].format == string_format {
				variables_used[i] = variables[i].options[configuration]
			} else if splits[i] == 1 {
				variables_used[i] = strconv.FormatFloat(variables[i].lower_boundary, 'f', 12, 64) 
			} else {
				value := variables[i].lower_boundary + step_size(variables[i], splits[i]) * float64(configuration)
				if variables[i].format == int_format {
					variables_used[i] = strconv.Itoa(int(value))
				} else {
					variables_used[i] = strconv.FormatFloat(value, 'f', 12, 64) 
				}

			}
		}
		score, _ := script_communicator.run_arguments(variables, variables_used)
		if score > best_score {
			best_score = score
			best_variables = variables_used
		}
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

	new_variables := make([]variable, len(variables))
	for i, old_variable := range variables {
		new_variables[i] = variable(old_variable)
		if old_variable.format == string_format {
			new_variables[i].options = []string{best_variables[i]}
		} else if old_variable.lower_boundary != old_variable.upper_boundary {
			step := step_size(old_variable, splits[i])
			best_value, _ := strconv.ParseFloat(best_variables[i], 64)
			if int(step) == 1 && old_variable.format == int_format {
				new_variables[i].lower_boundary = best_value
				new_variables[i].upper_boundary = best_value
			} else {
				if best_value - step > new_variables[i].lower_boundary {
					new_variables[i].lower_boundary = best_value - step
				}
				if best_value + step < new_variables[i].upper_boundary {
					new_variables[i].upper_boundary = best_value + step
				}
			}
		}

	}

	return new_variables
}

func (this grid_optimizer) find_optimal_hyperparameters(variables []variable) {
	grid_iteration(variables, this.script_communicator, this.runs, 0)
}

func (this iterated_grid_optimizer) find_optimal_hyperparameters(variables []variable) {
	iterations, splits_cap := this.get_params(len(variables))
	for i := 0; i < iterations; i++ {
		variables = grid_iteration(variables, this.script_communicator, max_int, splits_cap)
	}
}
