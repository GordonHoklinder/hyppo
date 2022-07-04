package main

type coordinate_optimizer struct {
	script_communicator communicator
	runs int
}

func (this coordinate_optimizer) find_optimal_hyperparameters(variables []variable) {
	possibilities := compute_prefix(possible_values(variables, 0), sum_no_negative)
	splits := find_splits(variables, possibilities, this.runs, get_addition_config())
	best := get_default_individual(variables)
	best = best.evaluate_individual(variables, this.script_communicator)
	for i := len(splits) - 1; i >= 0; i-- {
		for j := 0; j < splits[i]; j++ {
			current := best.make_copy()
			if variables[i].format == string_format {
				current.data[i] = property_from_string(variables[i].options[j])
			} else if splits[i] != 1 {
				step := step_size(variables[i], splits[i])
				if (variables[i].format == float_format) {
					current.data[i] = property_from_float(variables[i].lower_boundary + step * float64(j))
				} else {
					current.data[i] = property_from_int(int(variables[i].lower_boundary + step * float64(j)))
				}
			}
			current = current.evaluate_individual(variables, this.script_communicator)
			if (current.score > best.score) {
				best = current.make_copy()
			}
		}
	}
}
