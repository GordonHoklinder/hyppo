# Hyppo (HYPerParameter Optimizer)

Hyppo is intended to serve as a lightweight, command-line optimizer of various parameters of heuristic algoritms and neural networks or other machine learning algoritms.

## How to use it

### Instalation

To install hyppo on your computer, you need to have `git` and `go` installed. Then run the following commands:

```
git clone https://github.com/GordonHoklinder/hyppo
```

```
cd hyppo/src/; go build
```

To allow for running `hyppo` from any folder:
```
sudo ln -sr hyppo /bin/
```

### Typical usage

There are two things that need to be done before running hyppo on your script.

First, **the last line of your program's output must contain a score only,** where the score is a metric which the optimizer tries to maximize. Note that hyppo takes the last line of standard output and error output combined.

For example, when training a neural network, the format of a correct output could look like this:

```
...
Dev loss after epoch 9 is 0.8791
Dev accuracy after epoch 10 is 94.87
Dev loss after epoch 10 is 0.8502
94.87
```

If you wanted to minimize the dev loss instead of maximizing accuracy,
the lass line would contain `-0.8502` (the minus sign is there because hyppo only maximizes).

Second, **a file with hyperparameters and their ranges needs to be created.**
Its format is described in the subsection below.

When this initial setting is done, the script can be run as `hyppo --script=./heuristic_algorithm --optimizer=genetic --runs=100`. There are other flags which can be used. Meaning of all of them is described below.

The program prints the best score on the particular script so far and the best score achieved during this execution. Notice that the score is printed on a separate line so that the hyppo parameters for a given problem can be optimized using hyppo as well.

It writes all the parameters which were tried and the resulting score for each into a log file. The default path is `*name of the script*.hyppo-log`. This is done so even if hyppo is terminated in a middle of a run.

### Control flags

Hyppo recogizes the following flags.

- `script`: The path to the script whose hyperparameters are optimized. If the script needs to be run via interpreter, here should be the name of the interpreter (e.g. `--script=python3`).
- `arguments`: Arguments to the script which are passed before the arguments from optimizer. A possible usage is when using interpreter to enter the path to the program (e.g. `--script=python3 --arguments=main.py`). By default is empty. Note that at this point all the arguments are passed as one to the underlying script.
- `optimizer`: The type of optimizer used. They are described below. The default value is `genetic`.
- `runs`: The number of times the script should be run. Note that for some of the optimizers, this is only approximate.
- `variables`: The path to the file with hyperparameter names and their ranges. Default is `./variables.yaml`.
- `logs`: The path to file with logs. By default hyppo logs into a file `*script*.hyppo-log`, where `*script*` is the value passed to the script flag. If the file already exists, hyppo appends to it.
- `pass_names`: If true, the hyperparameters are passed to the script with their names (e.g. `./script --a=42 --b=37`), otherwise only the values are passed (e.g. `./script 42 37`).

There are other flags which are specific only for some optimizers. These are discussed in the subsection *Optimizers*.

### Format of variables file

The variables passed to the script are provided to hyppo via .yaml file passed as a flag `variables`.

The content of the file should be a list of dictionaries, each discribing one variable. The order in which they appear in the list is the same as the order in which they are passed to the underlaying script.

Each variable should contain the following:

- `name` (mandatory): The name of the variable.
- `format` (mandatory): The format of the variable - either `string`, `float` or `int`.
- `default` (optional): The default value for some optimizers. Should correspond to the format of variable.
- `options` (only for string variables, mandatory): A list of options.
- `lower` (only for nonstring variables, mandatory): The lower boundary of the interval from which the values are taken.
- `upper` (only for nonstring variables, optional): The upper boundary of the interval from which the values are taken. If not provided, the value of `lower` is taken.
- `splits` (only for nonstring variables, optional): For some optimizers, the number of splits in the given variable.

To get a better idea see the file `variables.yaml`.

### Optimizers

The following optimizers are available.
For a better idea how the optimizers are implemented, see the `DOCUMENTATION.md`.
If you would like to know the performance of the given algoritmhs, `COMPARISON.md` might be helpful - it contains the comparison of the optimizers on a concrete task.

#### Random Search
`--optimizer=random`

All the possibilities are tried with uniform probability. There are in total `runs` runs of the script and this is the only argument random search recognizes.

Note that the probability distribution in int and float variables is uniform. This may be inconvinient when searching for optimal value which may span multiple orders (for example the learning rate). Although there is no exponential probability distribution, there are two possible workarounds. Either provide the variable as a string variable, or provide another variable and comupute the desired in your script as an exponentioal of the one provided.

#### Grid Search
`--optimizer=grid`

Grid search recognizes only the flag `runs`. The number of call to the underlaying script should be close to this value, but are in general different.

If you want to have a better control over the number of splits in each variable, provide the argument `splits` for float/int variables in the .yaml config file. Note that for string variables all possibilities are always tried.

#### Iterative Grid Search
`--optimizer=iterative`

Similarly to grid search the `runs` flag determines the approximate number of calls to the script.

The number of iterations grid search does is computed automatically, it could be forced to a given value by passing the flag `iterations`. As in grid search one can use the `splits` atribute for a given variable to control the number of splits in each iteration.

#### Coordinate Search
`--optimizer=coordinate`

As in the previous two algorithms, one can use `runs` flag to determine approximate number of runs and provide `splits` attributes in each variable.

By default coordinate search starts with a random value for each variable. You can override this by passing a `default` in variable specification.

#### Genetic Algorithm
`--optimizer=genetic`

The number of iterations for this optimizer can be either specified as `iterations` or is computed roughly as a square root of the `runs`. The number of new hybrids and mutants in each generation can be passed as `hybrids` or `mutants` respectively. The size of the population is computed as runs / ((hybrids + mutants + 1) * iterations).

Last parameter is `magic`. It determines the standard deviation of a change when mutating a non-string variable. The deviation is computed as `magic` times the range of the variable.

#### Simulated Annealing
`--optimizer=annealing`

Simulated annealing recognizes the following flags:

- `runs`: There are in total `runs` runs.
- `temperature`: The initial temperature. Default `1.0`.
- `magic`: Determines the standard deviation of a change when mutating a non-string variable. The deviation is computed as `magic` times the range of the variable.

By default simulated annealing starts with a random value for each variable. You can override this by passing a `default` in variable specification.

### Examples

The folder `examples` contatins examples of typical usage.

The python script would typically be run as:

```
hyppo --script=python --arguments=python_example.py
```

Note that the parameters are passed with their names as `--pass_names` is by default true.

The bash script would typically be run as:

```
hyppo --script=./bash_example.sh --pass_names=false
```




## Credits

This project was created as a seminary project for the course Programming 2 in the Bachelor's programme Computer Science at Charles University.

You can learn more about the structure of the project and the used algorithms in `DOCUMENTATION.md`. In `COMPARISON.md` you can see a non-formal comparison of the optimizers.

The code in this repository can be used as described in LICENSE.

If you've found any issues or there are any features missing, do not hasitate to raise an issue or create a pull request.
