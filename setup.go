package main

import (
	"github.com/google/go-github/github"
	s "github.com/vivangkumar/statban/stats"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

type Config struct {
	Port         string
	Env          string
	GithubClient *github.Client
	Db           *s.Db
}

var (
	StatbanConfig *Config
	GithubConfig  *s.GithubConfig
)

// TODO: Use envdecode
func initialize() {
	c := Config{}
	ghConfig := s.GithubConfig{}
	db := s.Db{}

	if addr := os.Getenv("PORT"); addr != "" {
		c.Port = addr
	} else {
		c.Port = ":" + "8083"
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		c.Env = env
		ghConfig.PollInterval = 3600
		ghConfig.DailyInterval = 86400
	} else {
		c.Env = "development"
		ghConfig.PollInterval = 10
		ghConfig.DailyInterval = 20
	}

	if rdb := os.Getenv("RETHINK_DB_ADDR"); rdb != "" {
		db.Address = rdb
	} else {
		db.Address = "localhost:28015"
	}

	if sdb := os.Getenv("STATBAN_DB"); sdb != "" {
		db.Name = sdb
	} else {
		db.Name = "statban"
	}
	c.Db = &db
	StatbanConfig = &c

	if at := os.Getenv("GITHUB_TOKEN"); at != "" {
		ghConfig.Client = createGithubClient(at)
	} else {
		panic("Github token not set")
	}

	if repo := os.Getenv("TARGET_REPOSITORY"); repo != "" {
		ghConfig.TargetRepo = repo
	} else {
		panic("Target repository not set")
	}

	if owner := os.Getenv("REPO_OWNER"); owner != "" {
		ghConfig.Owner = owner
	} else {
		panic("Owner not set")
	}

	if labels := os.Getenv("LABELS"); labels != "" {
		ghConfig.Labels = strings.Split(labels, ",")
	} else {
		ghConfig.Labels = []string{"ready", "development", "review", "release", "done"}
	}

	GithubConfig = &ghConfig
}

func createGithubClient(accessToken string) (client *github.Client) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = github.NewClient(tc)

	return
}
