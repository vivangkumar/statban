package stats

import (
	"github.com/google/go-github/github"
	"log"
	"time"
)

type GithubConfig struct {
	Client        *github.Client
	TargetRepo    string
	Owner         string
	Labels        []string
	PollInterval  int
	DailyInterval int
}

func collector(ghConfig *GithubConfig, db *Db) {
	log.Printf("Getting issues for %v", ghConfig.TargetRepo)

	for _, lbl := range ghConfig.Labels {
		log.Printf("Getting issues with label %v", lbl)

		opt := github.IssueListByRepoOptions{
			State:  "open",
			Labels: []string{lbl},
		}
		issues, _, err := ghConfig.Client.Issues.ListByRepo(ghConfig.Owner, ghConfig.TargetRepo, &opt)
		if err != nil {
			log.Printf("Github error: %v", err.Error())
			return
		}

		if len(issues) > 0 {
			statbanIssues := make([]*StatbanIssue, len(issues))
			for i, issue := range issues {
				statbanIssues[i] = NewStatbanIssueFromGithubIssue(&issue)
			}
			db.StoreHourlyState(statbanIssues)
		} else {
			log.Printf("No issues for label %v", lbl)
		}
	}

	go db.SummarizeByHour(ghConfig)
}

func RunCollector(db *Db, ghConfig *GithubConfig) {
	log.Print("Running stats collector..")
	intervalTicker := time.NewTicker(time.Duration(ghConfig.PollInterval) * time.Second)
	dayTicker := time.NewTicker(time.Duration(ghConfig.DailyInterval) * time.Second)

	// TODO: This will leak memory, if the ticker is never closed
	// Maybe write a ticker manager which spawns new tickers?
	for {
		select {
		case <-intervalTicker.C:
			go collector(ghConfig, db)
		case <-dayTicker.C:
			go db.SummarizeByDay()
		}
	}
}
