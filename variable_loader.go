package main

import (
    "bufio"
    "os"
		"strconv"
		"strings"
)

type variable_format string

const (
	int_format variable_format = "int"
	float_format               = "float"
	string_format              = "string"
)

type variable struct {
	name string
	format variable_format
	lower_boundary float64
	upper_boundary float64
	splits int
	options []string
}

func load_variables(path string) ([]variable, error) {
	read_file, err := os.Open(path)
	result := make([]variable, 0);
	if err != nil {
		return result, err;
	}
	file_scanner := bufio.NewScanner(read_file)
	file_scanner.Split(bufio.ScanLines)

	for file_scanner.Scan() {
		line := strings.Split(file_scanner.Text(), "\t")
		if len(line) > 1 {
			name := line[0]
			format := variable_format(line[1])
			var lower_boundary, upper_boundary float64
			var splits int
			var err error
			var options []string
			if format == string_format {
				options = line[2:]
			} else {
				lower_boundary, err = strconv.ParseFloat(line[2] , 64)
				if err != nil {
					return result, err;
				}
				upper_boundary, err = strconv.ParseFloat(line[3] , 64)
				if err != nil {
					return result, err;
				}
				if len(line) > 4 {
					splits, err = strconv.Atoi(line[4])
					if err != nil {
						return result, err;
					}
				}
			}
			result = append(result, variable{name, format, lower_boundary, upper_boundary, splits, options})
		}
	}
	return result, nil
}
