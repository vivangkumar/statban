package stats

import "github.com/google/go-github/github"

func arrayDifference(slice1, slice2 []string) (diff []string) {
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return
}

func addMissingLabels(missing []string, ss *[]SummarizedState) *[]SummarizedState {
	for _, l := range missing {
		*ss = append(*ss, NewSummarizedState(l, 0))
	}

	return ss
}

func getMilestone(m *github.Milestone) string {
	if m == nil {
		return ""
	} else {
		return *m.Title
	}
}
