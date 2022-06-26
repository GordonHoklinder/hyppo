package main

import (
	"testing"
	"os"
)

func communicator_for_log_testing(log_path string) communicator {
	return communicator{"", "", log_path, true}
}

func communicator_for_format_testing(pass_variables bool) communicator {
	return communicator{"", "", "", pass_variables}
}

func variable_for_testing(name string) variable {
	return string_variable(name, []string{}, false, "")
}

func Test_get_log_path(t *testing.T) {
	test_cases := [][]string{
		{"script", "", "script.hyppo-log"},
		{" 78934.hihoafsd", "", " 78934.hihoafsd.hyppo-log"},
		{"script", "log", "log"},
	}
	for _, test := range test_cases {
		log_path := get_log_path(test[1], test[0])
		if test[2] != log_path {
			t.Logf("Log path should be %s but is %s", test[2], log_path)
			t.Fail()
		}
	}
}

func Test_read_log(t *testing.T) {
	f, _ := os.CreateTemp("", "tmp")
	defer f.Close()
	defer os.Remove(f.Name())
	f.WriteString("Score: -37.9 Flags: 37.8 5 Score: 32 <- this one is not the real score\nScore: 1 Flags: 2")
	script_communicator := communicator_for_log_testing(f.Name())
	lines := script_communicator.read_log()
	if len(lines) != 2 {
		t.Log("Incorrect number of read lines.")
		t.Fail()
	}
	if lines[0].score != -37.9 {
		t.Logf("Score should be %f but is %f", -37.9, lines[0].score)
		t.Fail()
	}
}

func Test_log_integration(t *testing.T) {
	f, _ := os.CreateTemp("", "tmp")
	defer f.Close()
	defer os.Remove(f.Name())
	script_communicator := communicator_for_log_testing(f.Name())
	scores := []float64{32.65, 0.0, 5.6}
	for _, score := range scores {
		script_communicator.log_score(score, "")
	}
	script_communicator.sort_log()
	result := script_communicator.read_log()[2].score
	if result != 0.0 {
		t.Logf("The last line of sorted log is not 0.0 but %f", result)
		t.Fail()
	}
	highest := 42.0
	script_communicator.log_score(highest, "")
	result = script_communicator.global_best_score()
	if result != highest {
		t.Logf("The best result is %f but %f was expected", result, highest)
		t.Fail()
	}
}

func Test_format_parameter(t *testing.T) {
	// Test short format
	script_communicator := communicator_for_format_testing(false)
	vari := variable_for_testing("variable")
	value := "value"
	result := script_communicator.format_parameter(vari, value)
	if result != value {
		t.Logf("Formatted parameter should be %s but is %s", value, result)
		t.Fail()
	}
	// Test long format
	script_communicator = communicator_for_format_testing(true)
	result = script_communicator.format_parameter(vari, value)
	expected := "--variable=value"
	if result != expected {
		t.Logf("Formatted parameter should be %s but is %s", expected, result)
		t.Fail()
	}
}


