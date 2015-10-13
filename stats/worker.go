package stats

import (
	"github.com/google/go-github/github"
	"github.com/nu7hatch/gouuid"
	"log"
	"time"
)

type GithubConfig struct {
	Client       *github.Client
	TargetRepo   string
	Owner        string
	Labels       []string
	PollInterval int
}

func collector(ghConfig *GithubConfig, db *Db) {
	log.Printf("Getting issues for %v", ghConfig.TargetRepo)

	batchId, _ := uuid.NewV4()
	for _, lbl := range ghConfig.Labels {
		log.Printf("Getting issues with label %v", lbl)

		opt := github.IssueListByRepoOptions{
			State:  "open",
			Labels: []string{lbl},
		}
		issues, _, err := ghConfig.Client.Issues.ListByRepo(ghConfig.Owner, ghConfig.TargetRepo, &opt)
		if err != nil {
			log.Printf("Github error: %v", err.Error())
		}

		if len(issues) > 0 {
			statbanIssues := make([]*StatbanIssue, len(issues))
			for i, issue := range issues {
				statbanIssues[i] = NewFromGithubIssue(&issue, batchId.String())
			}
			db.StoreDailyState(statbanIssues)
		} else {
			log.Printf("No issues for label %v", lbl)
		}
	}
}

func RunCollector(db *Db, ghConfig *GithubConfig) {
	log.Print("Running stats collector..")
	ticker := time.NewTicker(time.Duration(ghConfig.PollInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			collector(ghConfig, db)
		}
	}
}
