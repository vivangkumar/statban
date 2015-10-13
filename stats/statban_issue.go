package stats

import (
	"github.com/google/go-github/github"
	"time"
)

type StatbanIssue struct {
	IssueId        int       `gorethink:"issue_id,omitempty" json:"id"`
	Title          string    `gorethink:"title,omitempty" json:"title"`
	Label          string    `gorethink:"label,omitempty" json:"label"`
	Username       string    `gorethink:"user_name,omitempty" json:"user_name"`
	Milestone      string    `gorethink:"milestone",omitempty json:"milestone"`
	IssueCreatedAt time.Time `gorethink:"issue_created_at,omnitempty" json:"created_at"`
	CreatedAt      time.Time `gorethink:"created_at,omitempty" json:"recorded_at"`
	BatchId        string    `gorethink:"batch_id,omniempty" json:"batch_id"`
}

func NewFromGithubIssue(i *StatbanIssue, ghIssue *github.Issue, batchId string) *StatbanIssue {
	i.IssueId = *ghIssue.Number
	i.Title = *ghIssue.Title
	i.Label = *ghIssue.Labels[0].Name
	i.Username = *ghIssue.User.Login
	i.Milestone = getMilestone(ghIssue.Milestone)
	i.IssueCreatedAt = *ghIssue.CreatedAt
	i.CreatedAt = time.Now()
	i.BatchId = batchId

	return i
}

func getMilestone(m *github.Milestone) string {
	if m == nil {
		return ""
	} else {
		return *m.Title
	}
}
