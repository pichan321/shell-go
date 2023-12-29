package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Job struct {
	Pid   int
	Cmd   *exec.Cmd
	State int
}


type Jobs struct {
	JobList []Job
}

func InitJobs() *Jobs {
	return new(Jobs)
}

func (jobs *Jobs) AddJob(job Job) {
	jobs.JobList = append(jobs.JobList, job)
}

func (jobs *Jobs) GetJob(pid int) *Job {
	for _, job := range jobs.JobList {
		if job.Pid == pid {
			return &job
		}
	}
	return nil
}

func (jobs *Jobs) RemoveJob(jobToRemove *Job) {
	for idx, job := range jobs.JobList {
		if job == *jobToRemove {
			jobs.JobList = append(jobs.JobList[:idx], jobs.JobList[idx+1:]...)
			return
		}
	}
}

func (jobs *Jobs) ChangeState(jobToUpdate *Job, newState int) {
	for idx, job := range jobs.JobList {
		if job.Pid == (*jobToUpdate).Pid {
			jobs.JobList[idx].State = newState
		}
	}
}

func (jobs *Jobs) GetForegroundJob() *Job {
	for _, job := range jobs.JobList {
		if job.State == 1 {
			return &job
		}
	}
	return nil
}

func (jobs *Jobs) PrintJobs() {
	if len(jobs.JobList) <= 0 {return}
	fmt.Fprintf(os.Stderr, fmt.Sprintf("No.\tState\tPID\tCommand\n%s\n", strings.Repeat("-", 100)))
	for idx, job := range jobs.JobList {
		fmt.Fprintf(os.Stderr, "%d.\t%d\t[%d]\t%s\n", idx, job.State, job.Pid, job.Cmd)
	}
}
