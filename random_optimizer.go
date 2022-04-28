package main

import (
	"math/rand"
	"time"
)

type random_optimizer struct {
	script_communicator communicator
	runs int
}

func (this random_optimizer) find_optimal_hyperparameters(variables []variable) {
	rand.Seed(time.Now().UnixNano())
	for run := 0; run < this.runs; run++ {
		variable_values := get_random_individual(variables).data.to_string_slice(variables)
		this.script_communicator.run_arguments(variables, variable_values)
	}
}
