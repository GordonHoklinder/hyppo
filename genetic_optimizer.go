package main

import (
	"math/rand"
	"math"
	"time"
	"sort"
)

type genetic_optimizer struct {
	script_communicator communicator
	runs int
	mutation_deviation_coefficient float64
	mutants int
	hybrids int
	iterations int
}


func (this genetic_optimizer) get_iterations () int {
	if this.iterations == 0 {
		return int(math.Sqrt(float64(this.runs) / float64(this.mutants + this.hybrids)))
	} else {
		return this.iterations
	}
}

func (this genetic_optimizer) get_population_size () int {
	return int(0.5 + float64(this.runs) / float64(this.get_iterations()) / float64(this.hybrids + this.mutants))
}

func (this genetic_optimizer) kill_unfit (population []evaluated_individual) []evaluated_individual {
	sort.Slice(population, func(i, j int) bool {
		return population[i].score > population[j].score
	})
	return population[:this.get_population_size()]
}

func (this genetic_optimizer) mutate_string (mutant individual, index int, variables []variable) individual {
	mutant[index] = property_from_string(variables[index].options[rand.Intn(len(variables[index].options))])
	return mutant
}

func (this genetic_optimizer) get_mutated_value (original_value float64, index int, variables []variable) float64 {
	output := original_value + rand.NormFloat64() * (variables[index].upper_boundary - variables[index].lower_boundary) * this.mutation_deviation_coefficient
	output = math.Min(output, variables[index].upper_boundary)
	output = math.Max(output, variables[index].lower_boundary)
	return output
} 

func (this genetic_optimizer) mutate_float (mutant individual, index int, variables []variable) individual {
	mutant[index] = property_from_float(this.get_mutated_value(mutant[index].as_float, index, variables))
	return mutant
}

func (this genetic_optimizer) mutate_int (mutant individual, index int, variables []variable) individual {
	mutant[index] = property_from_int(int(this.get_mutated_value(float64(mutant[index].as_int), index, variables)))
	return mutant
}

func (this genetic_optimizer) get_mutant(parent individual, variables []variable) individual {
	mutant := make(individual, len(parent))
	copy(mutant, parent)
	index := rand.Intn(len(parent))
	if variables[index].format == string_format {
		this.mutate_string(mutant, index, variables)
	} else if variables[index].format == float_format {
		this.mutate_float(mutant, index, variables)
	} else {
		this.mutate_int(mutant, index, variables)
	}
	return mutant
}

func (this genetic_optimizer) get_hybrid (mother individual, father individual, variables []variable) individual {
	mutant := make(individual, len(mother))
	for i, father_property := range father {
		if rand.Int() % 2 == 0 {
			mutant[i] = father_property;
		}
	}
	return mutant
}


func (this genetic_optimizer) get_initial_population (variables []variable) []evaluated_individual {
	output := make([]evaluated_individual, 0)
	for i := 0; i < this.get_population_size(); i++ {
		output = append(output, get_random_individual(variables))
	}
	return output
} 

func (this genetic_optimizer) find_optimal_hyperparameters (variables []variable) {
	rand.Seed(time.Now().UnixNano())
	population := this.get_initial_population(variables)
	population_size := len(population)
	for i := 0; i < population_size; i++ {
			population[i] = population[i].evaluate_individual(variables, this.script_communicator)
	}
	iterations := this.get_iterations()
	for i := 0; i < iterations; i++ {
		for j := 0; j < population_size; j++ {
			for k := 0; k < this.mutants; k++ {
				population = append(population, evaluated_individual{0.0, this.get_mutant(population[j].data, variables)})
			}
			for k := 0; k < this.hybrids; k++ {
				second := rand.Intn(population_size)
				population = append(population, evaluated_individual{0.0, this.get_hybrid(population[j].data, population[second].data, variables)})
			}
		}
		for j := population_size; j < len(population); j++ {
			population[j] = population[j].evaluate_individual(variables, this.script_communicator)
		}
		population = this.kill_unfit(population)
	}
}


