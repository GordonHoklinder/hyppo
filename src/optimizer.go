package main

type optimizer interface {
	find_optimal_hyperparameters(varibles []variable)
}
