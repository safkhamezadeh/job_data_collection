package ranker

import (
	jobvacancies "job_vacancies/internal/job_vacancies"
	"job_vacancies/internal/keywordextractor"
	"slices"
	"strings"
)

type JobWithPoints struct {
	job    *jobvacancies.Job
	points int
}

type SimpleJobRanker struct {
	//pointsystem
}

func (SimpleJobRanker) RankJobs(keywords keywordextractor.KeyWordFormat, jobs []jobvacancies.Job) []jobvacancies.Job {
	var jobsWithPoints []JobWithPoints

	for i := range jobs {
		jobsWithPoints = append(jobsWithPoints, assignPoints(keywords, &jobs[i]))
	}

	slices.SortFunc(jobsWithPoints, cmpJobPoints)
	var jobsWithoutPoints []jobvacancies.Job
	for i := 0; i < len(jobsWithPoints); i++ {
		jobsWithoutPoints = append(jobsWithoutPoints, *jobsWithPoints[i].job)
	}

	return jobsWithoutPoints

}

func assignPoints(keywords keywordextractor.KeyWordFormat, job *jobvacancies.Job) JobWithPoints {
	var currPoints int
	maxPoints := getMaxPoints(keywords.JobTitles)
	for i := 0; i < len(keywords.JobTitles); i++ {
		if strings.Contains(job.Description, keywords.JobTitles[i]) {
			currPoints = currPoints + (maxPoints - i)
		}
	}
	maxPoints = getMaxPoints(keywords.Keywords)
	for i := 0; i < len(keywords.Keywords); i++ {
		if strings.Contains(job.Description, keywords.JobTitles[i]) {
			currPoints = currPoints + (maxPoints - i)
		}
	}
	return JobWithPoints{job: job, points: currPoints}
}

func getShortest[v any](sl1 []v, sl2 []v) int {
	if len(sl1) < len(sl2) {
		return len(sl1)
	}
	return len(sl2)
}

func getMaxPoints[v any](sl []v) int {
	return len(sl) * 2
}

func cmpJobPoints(a JobWithPoints, b JobWithPoints) int {
	if a.points < b.points {
		return -1
	} else if a.points > b.points {
		return 1
	} else {
		return 0
	}
}

//per job assign points to it
//assign points based on if description matches a word in keywords
//first index most points, last index least points
//so check len of index to decide how many points it should get

//point system :maxpoints = shortest len * 2?
//or check per slice and give maxpoints seperately?
