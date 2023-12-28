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

	isBuiltIn := handleBuiltIns(&parsedLine)
	if isBuiltIn {return}
	// dir := "/bin/"
	if bg {
		// cmd := exec.Command(parsedLine[0], parsedLine...)
		// cmd.Stdout = os.Stdout

		// err := cmd.Start()r

		// pid := cmd.Process.Pid
		// newJob := j.Job{
		// 	Pid:   pid,
		// 	Cmd:   parsedLine[0],
		// 	State: 1,
		// }
		// jobs.AddJob(newJob)
		// process, _ := os.FindProcess(pid)

		// if err != nil {

		// }
		// return
	}

	cmd := exec.Command("/Users/pichan/Desktop/projects/shell/hello", parsedLine[1:]...) //dir+"/"+parsedLine[0]
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Setpgid: true,
	//     Pgid: 0,
	// }

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	pid := cmd.Process.Pid

	// go func (exec *exec.Cmd) {
	// 	err := cmd.Wait()
	// 	if err != nil {fmt.Println(err)}
	// }(cmd)

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
		time.Sleep(time.Second * 1) //keep asking for update every 500 ms

		job = jobs.GetForegroundJob()

		if job == nil {return}
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

	if firstCmd == "jobs" {
		jobs.PrintJobs()
		return true
	}

	if firstCmd == "history" {
	}

	return false
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
	// if job == nil {
	// 	return
	// }

	process, err := os.FindProcess(job.Pid)
	if err != nil {
		return
	}
	syscall.Kill(process.Pid, syscall.SIGKILL)
}

func handleSigStop() {
	job := jobs.GetForegroundJob()
	// if job == nil {
	// 	return
	// }
	fmt.Println("GOT SIGSTOP")
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
