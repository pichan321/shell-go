package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	j "shell/jobs"
	"strings"
	"sync"
	"syscall"
	"time"
)

var jobs *j.Jobs = j.InitJobs()

const (
	FG  = 1
	BG  = 2
	STP = 3
)

var mutex sync.Mutex

func eval(line *string) {

	parsedLine, bg := parseline(line)

	handleBuiltIns(&parsedLine)
	// dir := "/bin/"
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

	cmd := exec.Command("/Users/pichan/Desktop/projects/shell/hello", parsedLine[1:]...) //dir+"/"+parsedLine[0]
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	pid := cmd.Process.Pid

	// if pid == 0 {
	
	// }

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

	for job.State == FG && job.Pid == pid {
		time.Sleep(time.Second * 1) //keep asking for update every 500 ms

		job := jobs.GetJob(pid)
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

		if status.Exited() || status.Signaled() {
			jobs.RemoveJob(*job)
			fmt.Println("removed", jobs)
			return
		}

		if status.Stopped() {
			fmt.Println("Before", job)
			// job.State = 3
			job.ChangeState(STP)
			fmt.Println("After", job)
			return
		}

		if status.Continued() {
			return
		}

	}

}

func handleSigIntSigKill() {
	job := jobs.GetForegroundJob()
	if job == nil {
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
	if job == nil {
		return
	}

	process, err := os.FindProcess(job.Pid)
	if err != nil {
		return
	}
	syscall.Kill(process.Pid, syscall.SIGTSTP)
}

func handleSig(sig chan os.Signal) {
	for {
		select {
		case incoming := <-sig:
			switch incoming {
			case syscall.SIGINT:
				handleSigIntSigKill()
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
