package jobs

import (
	"fmt"
	"os"
)

type Job struct {
	Pid   int
	Cmd   string
	State int
	Next  *Job
}

type Jobs struct {
	head *Job
}

func InitJobs() *Jobs {
	return new(Jobs)
}



func (jobs *Jobs) AddJob(job *Job) {
	if jobs.head == nil {
		jobs.head = job
		return
	}

	current := jobs.head
	for current.Next != nil {
		current = current.Next
	}

	current.Next = job
}

func (jobs *Jobs) DeleteJob(pid int) {
	if jobs.head.Pid == pid {
		jobs.head = jobs.head.Next
	}

	var prev *Job = nil
	current := jobs.head

	for current.Pid != pid {
		prev = current
		current = current.Next
	}

	prev.Next = current
}

func (jobs *Jobs) PrintJobs() {
	current := jobs.head

	for current.Next != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", current)
		current = current.Next
	}
}
