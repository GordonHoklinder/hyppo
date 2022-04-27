package main

import (
	"flag"
	"log"
	"fmt"
)

func main() {
	var path, optimizer_name, variables_path, language string
	var runs, iterations int
	var pass_variable_names bool
	flag.StringVar(&path, "path", "main.py", "The path to the algorithm whose hyperparameters are optimized.")
	flag.StringVar(&language, "language", "python", "The language shell command.")
	flag.StringVar(&optimizer_name, "optimizer", "random", "The type of optimizer used.")
	flag.StringVar(&variables_path, "variables", "variables.tsv", "The path to file with variables.")
	flag.IntVar(&runs, "runs", 1000, "The approximate number of times the optimizer runs the program.")
	flag.IntVar(&iterations, "iterations", 0, "The number of iterations the program does.")
	flag.BoolVar(&pass_variable_names, "pass_names", true, "Whether the variables to the script are passed with names.")
	flag.Parse();
	

	variables, err := load_variables(variables_path)
	if err != nil {
		log.Fatal(err)
	}
	script_communicator := communicator{path, language, pass_variable_names}
	var used_optimizer optimizer
	switch optimizer_name {
	case "random":
		used_optimizer = random_optimizer{script_communicator, runs}
	case "grid":
		used_optimizer = grid_optimizer{script_communicator, runs}
	case "iterative":
		used_optimizer = iterated_grid_optimizer{script_communicator, runs, iterations}
	default:
		log.Fatalf("%s is not a supported optimizer.\n", optimizer_name)
	}
	used_optimizer.find_optimal_hyperparameters(variables)
	fmt.Printf("Score:\n%f\n", script_communicator.best_score())
}
