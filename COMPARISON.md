# The comparison of optimizers

This comparison was done on our heuristic algorithm (which is not public at this moment) for the Reply Code Challenge. It is not intended as a rigorous evaluation of the optimizers, but rather to provide an idea how well do the models behave on a particular task.

The values for nondeterministic algorithms were taken as the average of 3 runs.

# A lot of runs
The first (01, not the sample) input of reply code challenge was quite small and our solution run there in less than a second, thus it was possible to run it many times and it still finished in a few minutes. The solution had 6 hyperparameters which were optimized.

All optimizers were run with `runs=3600`. Iterative grid search was not ran as it as a single iteration would need much more runs in the default settings.

The results were following:
- Random search ........... 102431
- Genetic algorithm ....... 102117
- Simulated annealing ..... 101330
- Grid search .............  96242
- Coordinate search .......  95832

# One variable
The same input, but only one variable was optimized. It was ran with `runs=730`.

The results were following:
- Simulated annealing ..... 100896
- Genetic algorithm ....... 100643
- Random search ........... 100616
- Grid search ............. 100463
- Coordinate search ....... 100463
- Iterative grid search ...  99753

# A few runs
The second input was a bit larger and the optimizers were run with `runs=80`. The iterative grid search was not ran for the same reason as in the first comparison.

- Grid search ............. 5610129
- Genetic algorithm ....... 5493290
- Random search ........... 5342324
- Coordinate search ....... 5264229
- Simulated annealing ..... 5260585

# Conclusion

Note again that this comparison is by no means comperhensive - the algorithms were ran on a very specific task and the conditions were not even fair (e.g. for `runs=80` grid search did only 64 calls). Nevertheless I would like to conclude the results anyway.

- Random search: Very consistent in the performance, always being among the best 3.
- Genetic algorithm: Seems to generally work well. It's advantage over random search is that it's able to exploit local optima better - this results in less balanced score, sometimes scoring poorly, but the maximum of the 3 runs was always higher than in random search.
- Simulated annealing: Appers to be very bad in situation where there is a vast parameter space and only a few runs. In the situation with only a small space and many runs it was superior.
- Grid search: Not great generally, if there are only a few runs, can benefit from trying the extreme values for the parameters.
- Iterated grid search & coordinate search: Seem to be bad.
