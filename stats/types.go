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

type SummarizedGroup struct {
	Group     string
	Reduction int
}

type SummarizedBatch struct {
	BatchId   string             `gorethink:"batch_id,omitempty" json:"batch_id"`
	States    *[]SummarizedState `gorethink:"states",omitempty json:"states"`
	CreatedAt time.Time          `gorethink:"created_at" json:"created_at"`
}

type SummarizedState struct {
	Label string `gorethink:"label,omitempty" json:"label"`
	Count int    `gorethink:"count",omitempty" json:"count" `
}

type SummarizedDay struct {
	*SummarizedBatch
	Begining time.Time `gorethink:"begining" json:"begining"`
	End      time.Time `gorethink:"end" json:"end"`
}

func NewSummarizedState(label string, count int) SummarizedState {
	return SummarizedState{Label: label, Count: count}
}

func NewSummarizedDay(sb *SummarizedBatch, beg time.Time, end time.Time) *SummarizedDay {
	return &SummarizedDay{
		SummarizedBatch: sb,
		Begining:        beg,
		End:             end,
	}
}

func NewFromStatbanIssueFromGithubIssue(ghIssue *github.Issue, batchId string) *StatbanIssue {
	return &StatbanIssue{
		IssueId:        *ghIssue.Number,
		Title:          *ghIssue.Title,
		Label:          *ghIssue.Labels[0].Name,
		Username:       *ghIssue.User.Login,
		Milestone:      getMilestone(ghIssue.Milestone),
		IssueCreatedAt: *ghIssue.CreatedAt,
		CreatedAt:      time.Now(),
		BatchId:        batchId,
	}
}
