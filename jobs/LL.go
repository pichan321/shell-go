package jobs

import (
	"fmt"
	"os"
)

type Job struct {
	Pid   int
	Cmd   string
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
		if job.Pid == pid {return &job}
	}
	return nil
}

func (jobs *Jobs) RemoveJob(jobToRemove Job) {
	for idx, job := range jobs.JobList {
		if job == jobToRemove {
			jobs.JobList = append(jobs.JobList[:idx], jobs.JobList[idx+1:]...)
			return
		}
	}	
}

func (job *Job) ChangeState(newState int) {
	// fmt.Println("HERE STTATE", job.State)
	job.State = newState
	// fmt.Println("HERE STTATE 2", job.State)
}

func (jobs *Jobs) GetForegroundJob() *Job {
	for _, job := range jobs.JobList {
		if job.State == 1 {return &job}
	}
	return nil
}

func (jobs *Jobs) PrintJobs() {
	fmt.Fprintf(os.Stderr, "Jobs:\n")
	for idx, job := range jobs.JobList {
		fmt.Fprintf(os.Stderr, "%d. %d %d %s\n", idx, job.State, job.Pid, job.Cmd)
	}
}
