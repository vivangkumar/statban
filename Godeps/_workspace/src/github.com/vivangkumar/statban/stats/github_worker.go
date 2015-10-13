package stats

import (
	r "github.com/statban/Godeps/_workspace/src/github.com/dancannon/gorethink"
	"github.com/statban/Godeps/_workspace/src/github.com/google/go-github/github"
	"github.com/statban/Godeps/_workspace/src/github.com/nu7hatch/gouuid"
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

func collector(ghConfig *GithubConfig, dbConfig *DbConfig) {
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
			extIssues := convertToStatbanIssues(issues, batchId.String())
			storeDailystate(ghConfig, dbConfig, extIssues)
		} else {
			log.Printf("No issues for label %v", lbl)
		}
	}
}

func storeDailystate(ghConfig *GithubConfig, dbConfig *DbConfig, issues []StatbanIssue) {
	for _, issue := range issues {
		_, err := r.DB(dbConfig.DbName).Table("hourly_state").Insert(issue).RunWrite(dbConfig.Session)
		if err != nil {
			log.Printf("Error inserting issue %v into table", issue)
		}
	}
}

func convertToStatbanIssues(issues []github.Issue, batchId string) []StatbanIssue {
	statbanIssues := make([]StatbanIssue, len(issues))

	for i, issue := range issues {
		statbanIssues[i] = StatbanIssue{
			IssueId:        *issue.Number,
			Title:          *issue.Title,
			Label:          *issue.Labels[0].Name,
			Username:       *issue.User.Login,
			Milestone:      getMilestone(issue.Milestone),
			IssueCreatedAt: *issue.CreatedAt,
			CreatedAt:      time.Now(),
			BatchId:        batchId,
		}
	}

	return statbanIssues
}

func RunCollector(dbConfig *DbConfig, ghConfig *GithubConfig) {
	log.Print("Running stats collector..")
	ticker := time.NewTicker(time.Duration(ghConfig.PollInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			collector(ghConfig, dbConfig)
		}
	}
}

func summarizeByInterval(ghConfig *GithubConfig) {

}

func getMilestone(m *github.Milestone) string {
	if m == nil {
		return ""
	} else {
		return *m.Title
	}
}
