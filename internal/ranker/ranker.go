package ranker

import (
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
	"slices"
	"strings"
)

type JobWithPoints struct {
	job    jobvacancies.Job
	points int
}

type SimpleJobRanker struct {
	//pointsystem
}

func (SimpleJobRanker) RankJobs(keywords keywordextractor.KeyWordFormat, jobs []jobvacancies.Job) []jobvacancies.Job {
	var jobsWithPoints []JobWithPoints

	for i := range jobs {
		jobsWithPoints = append(jobsWithPoints, assignPoints(keywords, jobs[i]))
	}

	slices.SortFunc(jobsWithPoints, cmpJobPoints)
	var jobsWithoutPoints []jobvacancies.Job
	for i := 0; i < len(jobsWithPoints); i++ {
		jobsWithoutPoints = append(jobsWithoutPoints, jobsWithPoints[i].job)
	}

	return jobsWithoutPoints

}

func assignPoints(keywords keywordextractor.KeyWordFormat, job jobvacancies.Job) JobWithPoints {
	title := strings.ToLower(job.Title)
	desc := strings.ToLower(job.Description)

	points := 0

	points += scoreList(title, desc, keywords.JobTitles, 3)
	points += scoreList(title, desc, keywords.Keywords, 1)

	return JobWithPoints{job: job, points: points}
}

func scoreList(title, desc string, list []string, multiplier int) int {
	points := 0
	max := len(list)

	for i, k := range list {
		kw := strings.ToLower(k)
		weight := max - i

		if strings.Contains(title, kw) {
			points += multiplier * weight * 3
		}
		if strings.Contains(desc, kw) {
			points += multiplier * weight
		}
	}

	return points
}

func cmpJobPoints(a JobWithPoints, b JobWithPoints) int {
	return b.points - a.points
}
