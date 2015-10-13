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

func NewFromGithubIssue(ghIssue *github.Issue, batchId string) *StatbanIssue {
	issue := &StatbanIssue{
		IssueId:        *ghIssue.Number,
		Title:          *ghIssue.Title,
		Label:          *ghIssue.Labels[0].Name,
		Username:       *ghIssue.User.Login,
		Milestone:      getMilestone(ghIssue.Milestone),
		IssueCreatedAt: *ghIssue.CreatedAt,
		CreatedAt:      time.Now(),
		BatchId:        batchId,
	}

	return issue
}

func getMilestone(m *github.Milestone) string {
	if m == nil {
		return ""
	} else {
		return *m.Title
	}
}
