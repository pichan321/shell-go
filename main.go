package main

import (
	"fmt"
	"strings"
)

func eval(line *string) {

	_, bg := parseline(line)

	if bg {

	}
}

func parseline(line *string) ([]string, bool) {
	splits := strings.Fields(*line)

	if splits[len(splits)-1] == "&" {
		return splits, true
	}

	return splits, false
}

func main() {
	for {
		line := ""
		fmt.Scanf("%s", &line)

		eval(&line)
	}
}
