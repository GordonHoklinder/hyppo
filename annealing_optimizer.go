package main

import (
	"math/rand"
	"time"
	"math"
)

type annealing_optimizer struct {
	script_communicator communicator
	runs int
	initial_temperature float64
	mutation_deviation_coefficient float64
}

func (this annealing_optimizer) temperature_schedule (step int) float64 {
	return this.initial_temperature * (1.0 - float64(step) / float64(this.runs))
}

func (this annealing_optimizer) accepted (old_score float64, new_score float64, temperature float64) bool {
	treshold := math.Exp((new_score - old_score) / temperature)
	return rand.Float64() < treshold
}

func (this annealing_optimizer) find_optimal_hyperparameters(variables []variable) {
	rand.Seed(time.Now().UnixNano())
	current := get_random_individual(variables)
	current = current.evaluate_individual(variables, this.script_communicator)
	basis := genetic_basis{this.mutation_deviation_coefficient}
	for i := 1; i < this.runs; i++ {
		temperature := this.temperature_schedule(i)
		next := evaluated_individual{0.0, basis.get_mutant(current.data, variables)}
		next = next.evaluate_individual(variables, this.script_communicator)
		if this.accepted(current.score, next.score, temperature) {
			current = next
		} 
	}
}
