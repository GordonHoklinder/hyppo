package main

import (
	"os"
	"os/exec"
	"log"
	"fmt"
	"strings"
	"strconv"
	"time"
	"bufio"
	"sort"
)

type communicator struct {
	script_path string
	language string
	pass_names bool
}

type line_with_score struct {
	score float64
	line string
}

func (this communicator) log_path () string {
	return this.script_path + ".hyppo-log"
}

func (this communicator) log_score (score float64, flags string) {
	file, _ := os.OpenFile(this.log_path(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) 
	defer file.Close()
	line := fmt.Sprintf("Score: %f Flags: %s\n", score, flags)
	file.WriteString(line)
}

func (this communicator) read_log() []line_with_score {
	file, _ := os.OpenFile(this.log_path(), os.O_RDWR, 0755)
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
	file, _ := os.OpenFile(this.log_path(), os.O_RDWR, 0755)
	file.Truncate(0)
	file.Seek(0, 0)
	for _, line := range lines {
		file.WriteString(line.line + "\n")
	}
}

func (this communicator) best_score () float64 {
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
	this.log_score(score, flags)
	return score, int64(execution_time)
}
