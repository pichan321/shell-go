package main

import (
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

func eval(line *string) {

	parsedLine, bg := parseline(line)

	handleBuiltIns(&parsedLine)
	dir := "/bin/"
	if bg {
		cmd := exec.Command(parsedLine[0], parsedLine...)
		cmd.Stdout = os.Stdout

		err := cmd.Start()

		pid := cmd.Process.Pid
		newJob := j.Job{
			Pid:   pid,
			Cmd:   parsedLine[0],
			State: 1,
		}
		jobs.AddJob(newJob)
		process, _ := os.FindProcess(pid)
		process.Wait()

		if err != nil {

		}
		return
	}

	cmd := exec.Command(dir+"/"+parsedLine[0], parsedLine[1:]...)
	cmd.Stdout = os.Stdout

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
	_, err := os.FindProcess(pid)
	if err != nil {
		return //not a valid process
	}
	job := jobs.GetJob(pid)
	if job == nil {
		return
	}

	for job.State == FG {
		time.Sleep(time.Millisecond * 100) //keep asking for update every 100 ms
		if job = jobs.GetJob(pid); job == nil {return}
	}
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
		jobs.PrintJobs()
	}

	if firstCmd == "history" {
	}
}

func handleSIGCHLD(sig os.Signal) {
	var status syscall.WaitStatus
	for {
		pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG|syscall.WUNTRACED, nil)
		if pid <= 0 || err != nil {
			return
		}
		job := jobs.GetJob(pid)

		if status.Exited() {
			jobs.RemoveJob(*job)

			fmt.Println("Status: ", status, "PID: ", pid)
			return
		}

		if status.Stopped() {
			job.ChangeState(STP)
			return
		}

		// if status.Continued() {
		// 	job.ChangeState(R)
		// 	return
		// }

	}

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
				handleSIGCHLD(incoming)

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
