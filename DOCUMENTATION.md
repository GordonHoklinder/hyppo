# Programmer's documentation

This file provides an overview of how the project is structured and how are individual optimizers implemented.

All the code is in the subdirectory `src`.

## Controllers

### `main.go`

This is the entry point to the program. It defines and loads the flags and ensures all relevant controllers and the appropriate optimizer are called.

### `variable_loader.go`

Defines the `variable` struct and implements code for loading variables from the given .yaml file.

### `communicator.go`

This script is called from the optimizers (their description is below). It communicates with the underlaying script and stores the passed arguments and achieved results.

## Optimizers

### `optimizer.go`
This is the interface for all the optimizers.
The optimizers in general select several values to be passed to the script which is done via `communicator`.

### `random_optimizer.go`

Implements random search.

Tries all the valid combinations of values with uniform probability.

### `grid_optimizer.go`

Implements grid search and iterated grid search.

Grid search divides selects several (almost) equidistant values in each variable and tries all the combinations.

Iterated grid search does the same, but in several iterations. In each iteration, it finds the best combination of values and runs the next iteration on all the neighbouring cells.

Apart from that, these searches try to find an appropriate number of values in each variable.

Grid search always makes `splits` splits in variables where the value of `splits` is provided. In string variables it tries all the possibilities.

For other variables it computes the number of splits using the following algorithm.

At first it computer for each variable the number of possible combinations of values (with respect to `splits`, string variables and int variables) of previous variables (lets call it c). Then it assigns the number of splits from the last variable for a variable i, as: max (runs/c, i-th root of runs). And divides runs by this variable.

There were three main reasons for using this ad-hoc algorithm:

- If all the variables are float variables without provided split value, this result in equal (+-1) number of splits in each variable.
- With increasing `runs` the total number of calls increases with the same speed, even if there are variables with provided `splits`, string variables or int variables with a small range.
- The variables listed first are considered more important.

The iterated grid search respects the string variables and variables with `splits` provided, otherwise it does 9 splits.

The reason for this is that it is the value which minimizes the function (3/x) ^ (n/x), which is the length of the hypercube we end up with after the last iteration when running n times. The minimum of this function is in 3e, to which the closest natural number is 9.

### `coordinate_optimizer.go`

Works similarly to grid search but optimizes the variables one by one. It starts with a default values (or a random ones if defaults not provided). It selects the best value in the axis from those tried and continues by changing the best so far in the next axis.

The computation of the number of splits in each dimension is almost the same to the grid search. It only takes max(runs - c, runs/i) as the number of splits (for obvious reasons).

### `genetic_optimizer.go`

Maintatins a population of individuals, initialized randomly. It always leaves only those with highest achieved score and creates for each individual several mutants and hybrids (the concrete numbers can be specified using flags).

A mutant is created by mutating one variable (there should not be generally many variables so one should be sufficient).

A hybrid is created by randomly selecting second individual and chosing from one of the parents with the same probability.

### `annealing_optimizer.go`

Implements a simple version of the simulated annealing algorithm.

The temperature schedule for the simulated annealing is always taken to be linearly decreasing from initial temperature to zero.

The acceptance criterium is if a random variable between 0 and 1 is less than exp((new - old) / temperature).


## Other

### `individual.go`

Provides the struct and functions for easy manipulation with the values of variables, especially used in genetic optimizer.

### `miscellaneous.go`

Provides additional miscellaneous functions which are used all over the project.

## Tests

### `loader_test.go`

Provides a test for testing the variable loading in `variable_loader.go`.

### `communicator.go`

Provides tests for testing `communicator.go`.

