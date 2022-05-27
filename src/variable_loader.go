package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v3"
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
	default_value property
	has_default_value bool
	options []string
}


func int_variable (name string, lower_boundary float64, upper_boundary float64, splits int, has_default_value bool, default_value int) variable {
	return variable {
		name,
		int_format,
		lower_boundary,
		upper_boundary,
		splits,
		property_from_int(default_value),
		has_default_value,
		make([]string, 0),
	}
}

func float_variable (name string, lower_boundary float64, upper_boundary float64, splits int, has_default_value bool, default_value float64) variable {
	return variable {
		name,
		float_format,
		lower_boundary,
		upper_boundary,
		splits,
		property_from_float(default_value),
		has_default_value,
		make([]string, 0),
	}
}

func string_variable (name string, options []string, has_default_value bool, default_value string) variable {
	return variable {
		name,
		string_format,
		0.0,
		0.0,
		0,
		property_from_string(default_value),
		has_default_value,
		options,
	}
}

func load_variables(path string) ([]variable, error) {
	result := make([]variable, 0);

	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		return result, err
	}

	data := make([]map[string]interface{}, 0)

	err = yaml.Unmarshal(yfile, &data)
	if err != nil {
		 return result, err
	}
	
	for _, current_variable := range data {
		name := current_variable["name"].(string)
		format := current_variable["format"].(string)
		default_value, has_default_value := current_variable["default"] 
		if format == string_format {
			options_unsafe := current_variable["options"].([]interface{})
			options := make([]string, len(options_unsafe))
			for i, option := range options_unsafe {
				options[i] = option.(string)
			}
			def := ""
			if has_default_value {
				def = default_value.(string)
			}
			result = append(result, string_variable(name, options, has_default_value, def))
		} else {
			var lower_boundary float64
			if variable_format(format) == int_format {
				lower_boundary = float64(current_variable["lower"].(int))
			} else {
				lower_boundary = current_variable["lower"].(float64)
			}
			upper_boundary := lower_boundary
			if upper_boundary_value, ok := current_variable["upper"]; ok {
				if variable_format(format) == int_format {
					upper_boundary = float64(upper_boundary_value.(int))
				} else {
					upper_boundary = upper_boundary_value.(float64)
				}
			}
			splits := 0
			if splits_value, ok := current_variable["splits"]; ok {
				splits = splits_value.(int)
			}

			if variable_format(format) == int_format {
				def := 0
				if has_default_value {
					def = default_value.(int)
				}
				result = append(result, int_variable(name, lower_boundary, upper_boundary, splits, has_default_value, def))
			} else {
				def := 0.0
				if has_default_value {
					def = default_value.(float64)
				}
				result = append(result, float_variable(name, lower_boundary, upper_boundary, splits, has_default_value, def))

			}
			

		}
		 
	}


	return result, nil
}
