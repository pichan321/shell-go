package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	j "shell/jobs"
	"strings"
	"syscall"
	"time"
)

var jobs *j.Jobs = j.InitJobs()

const (
	FG  = 1
	BG  = 2
	STP = 3
)

func overwriteFileDescriptor(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
}

func eval(line *string) {
	parsedLine, bg := parseline(line)

	isBuiltIn := handleBuiltIns(&parsedLine)
	if isBuiltIn {
		return
	}

	if bg {
		cmd := exec.Command("/Users/pichan/Desktop/projects/shell/hello", parsedLine...)

		err := cmd.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to initiate background job | %s\n", cmd.String())
		}
		pid := cmd.Process.Pid
		newJob := j.Job{
			Pid:   pid,
			Cmd:   cmd.String(),
			State: BG,
		}

		jobs.AddJob(newJob)
		fmt.Fprintf(os.Stderr, "+1 [%d] %s\n", pid, cmd.String())
		return
	}

	cmd := exec.Command("/Users/pichan/Desktop/projects/shell/hello", parsedLine[1:]...) //dir+"/"+parsedLine[0]
	overwriteFileDescriptor(cmd)

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	pid := cmd.Process.Pid

	jobs.AddJob(j.Job{
		Pid:   pid,
		Cmd:   cmd.String(),
		State: FG,
	})

	waitForeground(pid)
}

func waitForeground(pid int) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return //not a valid process
	}
	job := jobs.GetForegroundJob()
	if job.Pid != process.Pid {
		return
	}

	for job.State == FG && job.Pid == pid {
		time.Sleep(time.Millisecond * 500) //keep asking for update every 500 ms

		job = jobs.GetForegroundJob()

		if job == nil {
			return
		}
		if job.State == STP {
			return
		}
	}
}

func parseline(line *string) ([]string, bool) {
	splits := strings.Fields(*line)

	if splits[len(splits)-1] == "&" {
		return splits, true
	}

	return splits, false
}

func handleBuiltIns(parsedLine *[]string) bool {
	if len(*parsedLine) <= 0 {
		return false
	}

	firstCmd := string((*parsedLine)[0])
	if firstCmd == "quit" {
		os.Exit(0)
		return true
	}

	if firstCmd == "jobs" || firstCmd == "ps" {
		jobs.PrintJobs()
		return true
	}

	if firstCmd == "fg" || firstCmd == "bg" {
		startFgBg(parsedLine)
		return true
	}

	if firstCmd == "history" {
	}

	return false
}

func startFgBg(parsedLine *[]string) {
	if (*parsedLine)[0] == "fg" {
		
	}
	if (*parsedLine)[0] == "bg" {
		
	}
}

func handleSIGCHLD(sig os.Signal) {
	var status syscall.WaitStatus
	var job *j.Job
	for {
		pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG|syscall.WUNTRACED, nil)
		if pid <= 0 || err != nil {
			return
		}
		job = jobs.GetJob(pid)
		if job == nil {return}

		if status.Exited() {
			jobs.RemoveJob(job)
			return
		}

		if status.Signaled() {
			fmt.Fprintf(os.Stderr, "Terminated by signal\n")
			jobs.RemoveJob(job)
			return
		}

		if status.Stopped() {
			jobs.ChangeState(job, BG)
			return
		}

		if status.Continued() {
			return
		}

	}

}

func handleSigInt() {
	job := jobs.GetForegroundJob()
	if job == nil {
		if len(jobs.JobList) <= 0 {
			os.Exit(0)
		}
		return
	}

	process, err := os.FindProcess(job.Pid)
	if err != nil {
		return
	}
	syscall.Kill(process.Pid, syscall.SIGKILL)
}

func handleSigStop() {
	job := jobs.GetForegroundJob()
	if job == nil {return}

	process, err := os.FindProcess(job.Pid)
	if err != nil {
		return
	}

	process.Signal(syscall.SIGTSTP)
}

func handleSig(sig chan os.Signal) {
	for {
		select {
		case incoming := <-sig:
			switch incoming {
			case syscall.SIGINT:
				handleSigInt()
			case syscall.SIGTSTP:
				handleSigStop()
			case syscall.SIGCHLD:
				handleSIGCHLD(incoming)

			}
		}
	}
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGCHLD, syscall.SIGINT, syscall.SIGTSTP)
	go handleSig(sig)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "> ")
		scanner.Scan()
		line := scanner.Text()
		if line == "" {
			continue
		}
		eval(&line)
	}
}
