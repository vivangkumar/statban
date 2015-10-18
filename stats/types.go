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
}

type SummarizedHour struct {
	States    *[]SummarizedState `gorethink:"states,omitempty" json:"states"`
	CreatedAt time.Time          `gorethink:"created_at,omitmepty" json:"created_at"`
	HourStart time.Time          `gorethink:"hour_start,omitempty json:"hour_start"`
	HourEnd   time.Time          `gorethink:"hour_end,omitempty json:"hour_end"`
}

type SummarizedDay struct {
	*SummarizedHour
	DayStart time.Time `gorethink:"day_start,omitempty" json:"day_start"`
	DayEnd   time.Time `gorethink:"day_start,omitempty" json:"day_end"`
}

type SummarizedState map[string]int

type SummarizedGroup struct {
	Group     string
	Reduction int
}

func NewSummarizedState(label string, count int) SummarizedState {
	s := make(SummarizedState)
	s[label] = count
	return s
}

func NewSummarizedDay(sb *SummarizedHour, beg time.Time, end time.Time) *SummarizedDay {
	return &SummarizedDay{
		SummarizedHour: sb,
		DayStart:       beg,
		DayEnd:         end,
	}
}

func NewStatbanIssueFromGithubIssue(ghIssue *github.Issue) *StatbanIssue {
	return &StatbanIssue{
		IssueId:        *ghIssue.Number,
		Title:          *ghIssue.Title,
		Label:          *ghIssue.Labels[0].Name,
		Username:       *ghIssue.User.Login,
		Milestone:      getMilestone(ghIssue.Milestone),
		IssueCreatedAt: *ghIssue.CreatedAt,
		CreatedAt:      time.Now(),
	}
}
