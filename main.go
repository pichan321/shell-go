package main

import (
	"fmt"
	"os"
	"os/signal"
	j "shell/jobs"
	"strings"
	"syscall"
)

var jobs *j.Jobs = j.InitJobs()

func eval(line *string) {

	parsedLine, bg := parseline(line)

	handleBuiltIns(&parsedLine)

	if bg {
		go func() {
			_, err := syscall.ForkExec("./"+parsedLine[0], parsedLine, &syscall.ProcAttr{Files: []uintptr{0, 1, 2}})
			if err != nil {
				fmt.Println(err.Error())
			}

		}()
		return
	}

	pid, err := syscall.ForkExec("./"+parsedLine[0], parsedLine, &syscall.ProcAttr{Files: []uintptr{0, 1, 2}})
	if err != nil {
		fmt.Println(err.Error())
	}

	newJob := j.Job{
		Pid: pid,
		Cmd: parsedLine[0],
		State: 1,
		Next: nil,

	}

	if pid != 0 {
		jobs.AddJob(&newJob)
		process, _ := os.FindProcess(pid)
		process.Wait()
		jobs.PrintJobs()
	}
	

	// if pid != 0 {
	// 	fmt.Println(os.Getpid(), pid)
	// 	childProcess, err := os.FindProcess(pid)
	// 	if err != nil {return}

	// 	if bg {

	// 	}
	// 	_, err = childProcess.Wait()

	// }
}

func waitForeground(pid int) {

}

func parseline(line *string) ([]string, bool) {
	splits := strings.Fields(*line)

	if splits[len(splits)-1] == "&" {
		return splits, true
	}

	return splits, false
}

func handleBuiltIns(parsedLine *[]string) {
	if len(*parsedLine) <= 0 {
		return
	}

	firstCmd := string((*parsedLine)[0])
	if firstCmd == "quit" {
		os.Exit(0)
	}
	if firstCmd == "jobs" {

	}

	if firstCmd == "history" {
	}
}

func handleSIGCHLD() {

}

func handleSig(sig chan os.Signal) {
	for {
		select {
		case incoming := <-sig:
			switch incoming {
			case syscall.SIGINT:
				os.Exit(0)
			case syscall.SIGKILL:
				os.Exit(0)
			case syscall.SIGCHLD:
				handleSIGCHLD()

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
		fmt.Fprintf(os.Stderr, "> ")
		fmt.Scanln(&line)

		if line == "" {
			continue
		}
		eval(&line)
	}
}
