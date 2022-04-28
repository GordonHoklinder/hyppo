package main

import (
	"math/rand"
	"strconv"
)

const max_int = int(^uint(0) >> 1)

func min(a int, b int) int {
	if (a < b) {
		return a
	}
	return b
}

func max(a int, b int) int {
	if (a > b) {
		return a
	}
	return b
}

type property struct {
	as_string string
	as_float float64
	as_int int
}

func property_from_string (value string) property {
	return property{value, 0.0, 0}
}

func property_from_int (value int) property {
	return property{"", 0.0, value}
}

func property_from_float (value float64) property {
	return property{"", value, 0}
}

type individual []property

type evaluated_individual struct {
	score float64
	data individual
}

func get_random_individual (variables []variable) evaluated_individual {
	newborn := make(individual, len(variables))
	for i, vari := range variables {
		if vari.format == string_format {
			newborn[i] = property_from_string(vari.options[rand.Intn(len(vari.options))])
		} else if vari.format == int_format {
			lower := int(variables[i].lower_boundary)
			upper := int(variables[i].upper_boundary)
			newborn[i] = property_from_int(rand.Intn(upper - lower + 1) + lower)
		} else {
			lower := variables[i].lower_boundary
			upper := variables[i].upper_boundary
			newborn[i] = property_from_float(lower + rand.Float64() * (upper - lower))
		}
	}
	return evaluated_individual{0.0, newborn}
}

func (this individual) to_string_slice (variables []variable) []string {
	output := make([]string, len(variables)) 
	for i, prop := range this {
		if variables[i].format == string_format {
			output[i] = prop.as_string
		} else if variables[i].format == int_format {
			output[i] = strconv.Itoa(prop.as_int)
		} else {
			output[i] = strconv.FormatFloat(prop.as_float, 'f', 12, 64)
		}
	}
	return output
}
