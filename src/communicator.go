package main

import (
	"os"
	"os/exec"
	"log"
	"fmt"
	"strings"
	"strconv"
	"bufio"
	"sort"
	"math"
)

var current_best_score float64

type communicator struct {
	script string
	arguments string
	log_path string
	pass_names bool
}


func new_communicator (script string, arguments string, log_path string, pass_names bool) communicator {
	log_path = get_log_path(log_path, script)
	current_best_score = math.Inf(-1)
	return communicator{script, arguments, log_path, pass_names}
}

func get_log_path (log_path, script string) string {
	if log_path == "" {
		return script + ".hyppo-log"
	} else {
		return log_path
	}
}

type line_with_score struct {
	score float64
	line string
}


func (this communicator) log_score (score float64, flags string) {
	file, _ := os.OpenFile(this.log_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) 
	defer file.Close()
	line := fmt.Sprintf("Score: %f Flags: %s\n", score, flags)
	file.WriteString(line)
}

func (this communicator) read_log() []line_with_score {
	file, _ := os.OpenFile(this.log_path, os.O_RDWR, 0755)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := make([]line_with_score, 0)
	for scanner.Scan() {
		line := scanner.Text()
		score, _:= strconv.ParseFloat(strings.Split(line, " ")[1],64)
		lines = append(lines, line_with_score{score, line})
	}
	return lines
}

func (this communicator) sort_log () {
	lines := this.read_log()
	sort.Slice(lines, func(i, j int) bool {
    return lines[i].score > lines[j].score
	})
	file, _ := os.OpenFile(this.log_path, os.O_RDWR, 0755)
	file.Truncate(0)
	file.Seek(0, 0)
	for _, line := range lines {
		file.WriteString(line.line + "\n")
	}
}

func (this communicator) global_best_score () float64 {
	this.sort_log()
	return this.read_log()[0].score
}

func (this communicator) format_parameter(argument variable, value string) string {
	if this.pass_names {
		return "--" + argument.name + "=" + value
	} else {
		return value
	}
}

func (this communicator) format_flags(arguments []variable, arguments_values []string) []string {
	flags := make([]string, 0)
	for i := 0; i < len(arguments); i++ {
		flags = append(flags, this.format_parameter(arguments[i], arguments_values[i]))
	}
	return flags
}

func (this communicator) run_arguments (arguments []variable, arguments_values []string) float64 {
	flags := this.format_flags(arguments, arguments_values)
	var command *exec.Cmd
	if this.arguments != "" {
		flags = append([]string{this.arguments}, flags...)
	}
	command = exec.Command(this.script, flags...)
	out, execution_error := command.Output()
  if execution_error != nil {
  	log.Fatal(out, execution_error)
  }
	out_lines := strings.Split(strings.Trim(string(out), "\n"), "\n")
	score, parse_error := strconv.ParseFloat(out_lines[len(out_lines) - 1], 64)
  if parse_error != nil {
  	log.Fatal(parse_error)
  }
	this.log_score(score, strings.Join(flags, " "))
	current_best_score = math.Max(score, current_best_score)
	return score
}
