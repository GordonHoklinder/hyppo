package main

import (
	"os/exec"
	"log"
	"fmt"
	"strings"
	"strconv"
	"time"
)

type communicator struct {
	script_path string
	language string
	pass_names bool
}

func (this communicator) format_parameter(argument variable, value string) string {
	if this.pass_names {
		return "--" + argument.name + "=" + value
	} else {
		return value
	}
}


func (this communicator) format_command() string {
	if this.language == "" {
		return this.script_path
	} else {
		return fmt.Sprintf("%s %s", this.language, this.script_path)
	}
}

func (this communicator) format_flags(arguments []variable, arguments_values []string) string {
	flags := ""
	for i := 0; i < len(arguments); i++ {
		flags += " " + this.format_parameter(arguments[i], arguments_values[i])
	}
	return flags
}

func (this communicator) run_arguments (arguments []variable, arguments_values []string) (float64, int64) {
	flags := this.format_flags(arguments, arguments_values)
	command := this.format_command()
	execution_start := time.Now()
	out, execution_error := exec.Command(command, flags).Output()
	execution_time := time.Since(execution_start)
  if execution_error != nil {
  	log.Fatal(execution_error)
  }
	out_lines := strings.Split(strings.Trim(string(out), "\n"), "\n")
	score, parse_error := strconv.ParseFloat(out_lines[len(out_lines) - 1], 64)
  if parse_error != nil {
  	log.Fatal(parse_error)
  }
	fmt.Printf("Score: %f\n", score)
	return score, int64(execution_time)
}
