package main

import (
	"github.com/google/go-github/github"
	s "github.com/vivangkumar/statban/stats"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	HttpAddress  string
	Env          string
	GithubClient *github.Client
	DbConfig     *s.DbConfig
}

var (
	StatbanConfig *Config
	GithubConfig  *s.GithubConfig
)

func init() {
	c := Config{}
	ghConfig := s.GithubConfig{}
	dbConfig := s.DbConfig{}

	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		c.HttpAddress = addr
	} else {
		c.HttpAddress = "localhost:8083"
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		c.Env = env
		ghConfig.PollInterval = 86400
	} else {
		c.Env = "development"
		ghConfig.PollInterval = 10
	}

	if rdb := os.Getenv("RETHINK_DB_ADDR"); rdb != "" {
		dbConfig.Address = rdb
	} else {
		dbConfig.Address = "localhost:28015"
	}

	dbConfig.Tables = []string{"hourly_state", "daily_state"}

	if db := os.Getenv("STATBAN_DB"); db != "" {
		dbConfig.DbName = db
	} else {
		dbConfig.DbName = "statban"
	}
	c.DbConfig = &dbConfig
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

func main() {
	cfg, err := StatbanConfig.DbConfig.Setup()
	if err != nil {
		panic(err.Error())
	}

	go s.RunCollector(cfg, GithubConfig)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Statban Server"))
	})

	log.Printf("Statban server running on %v", StatbanConfig.HttpAddress)
	http.ListenAndServe(StatbanConfig.HttpAddress, nil)
}

func createGithubClient(accessToken string) (client *github.Client) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = github.NewClient(tc)

	return
}
