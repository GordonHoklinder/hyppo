package main

import (
	"math/rand"
	"time"
	"strconv"
)

type random_optimizer struct {
	script_communicator communicator
	runs int
}

func (this random_optimizer) find_optimal_hyperparameters(variables []variable) {
	rand.Seed(time.Now().UnixNano())
	for run := 0; run < this.runs; run++ {
		variable_values := make([]string, len(variables))
		for i := 0; i < len(variables); i++ {
			if variables[i].format == string_format {
				variable_values[i] = variables[i].options[rand.Intn(len(variables[i].options))]
			} else if variables[i].format == int_format {
				lower := int(variables[i].lower_boundary)
				upper := int(variables[i].upper_boundary)
				variable_values[i] = strconv.Itoa(rand.Intn(upper - lower) + lower)
			} else {
				lower := variables[i].lower_boundary
				upper := variables[i].upper_boundary
				variable_values[i] = strconv.FormatFloat(lower + rand.Float64() * (upper - lower), 'f', 12, 64)
			}
		}
		this.script_communicator.run_arguments(variables, variable_values)
	}
}
