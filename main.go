package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

func handleSIGCHLD() {
	
}

func handleSig(sig chan os.Signal) {
	for {
		select {
		case incoming := <-sig:
			switch incoming {
				case syscall.SIGINT: os.Exit(0)
				case syscall.SIGKILL: os.Exit(0)
				case syscall.SIGCHLD: handleSIGCHLD()
			}
		
	
		}
	}
}

func main() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGCHLD, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTSTP)
	go handleSig(sig)

	for {
		line := ""
		fmt.Scanf("%s", &line)

		eval(&line)
	}
}
