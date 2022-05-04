package main

import (
	"flag"
	"log"
	"fmt"
)

func main() {
	var script, optimizer_name, variables_path, logs_path, arguments string
	var runs, iterations, mutants, hybrids int
	var magic float64
	var pass_variable_names bool
	flag.StringVar(&script, "script", "./main", "Path to the script whose hyperparameters are optimized or to the interpreter.")
	flag.StringVar(&arguments, "arguments", "", "The arguements passed to the scrept before the hyperparameters. Usually used for passing script path to the interpreter.")
	flag.StringVar(&optimizer_name, "optimizer", "random", "The type of optimizer used.")
	flag.StringVar(&variables_path, "variables", "variables.tsv", "The path to file with variables.")
	flag.StringVar(&logs_path, "logs", "", "The path to file with logs.")
	flag.IntVar(&runs, "runs", 1000, "The approximate number of times the optimizer runs the program.")
	flag.IntVar(&iterations, "iterations", 0, "The number of iterations the program does.")
	flag.IntVar(&mutants, "mutants", 7, "The number of mutants used in each iteration.")
	flag.IntVar(&hybrids, "hybrids", 3, "The number of hybrids used in each iteration.")
	flag.Float64Var(&magic, "magic", 0.3, "A magical constant.")
	flag.BoolVar(&pass_variable_names, "pass_names", true, "Whether the variables to the script are passed with names.")
	flag.Parse();
	

	variables, err := load_variables(variables_path)
	if err != nil {
		log.Fatal(err)
	}
	script_communicator := communicator{script, arguments, logs_path, pass_variable_names}
	var used_optimizer optimizer
	switch optimizer_name {
	case "random":
		used_optimizer = random_optimizer{script_communicator, runs}
	case "grid":
		used_optimizer = grid_optimizer{script_communicator, runs}
	case "iterative":
		used_optimizer = iterated_grid_optimizer{script_communicator, runs, iterations}
	case "genetic":
		used_optimizer = genetic_optimizer{script_communicator, runs, magic, mutants, hybrids, iterations}
	case "coordinate":
		used_optimizer = coordinate_optimizer{script_communicator, runs}
	default:
		log.Fatalf("%s is not a supported optimizer.\n", optimizer_name)
	}
	used_optimizer.find_optimal_hyperparameters(variables)
	fmt.Printf("Score:\n%f\n", script_communicator.best_score())
}
