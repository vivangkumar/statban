# statban

A simple Kanban workflow stats collection tool.

Statban simply polls a particular Github repository and looks for issues
with certain labels and aggregates them into an hourly and daily view.

## Components

- **Database**: Statban is backed by RethinkDB. There is no reason for
choosing RethinkDB. Merely a way to play around with new technologies.
JSON based storage was preferred since the data is simply pushed to a
database with minimal aggregation.

- **HTTP Server**: Statban runs a simple HTTP server with a couple of endpoints
that serves JSON. It is simply a representation of the RethinkDB data.

- **Collectors**: Simple functions that runs as Go routines that perform basic
aggregation and push to the database. One which runs every hour and another every
24 hours. These are configurable using Environment variables, mostly for development.

## Environment variables

- `HTTP_ADDR`: The host name and port at which the server is run. Defaults to `localhost:8083`.
- `ENVIRONMENT`: Set the environment for development or production.
- `RETHINK_DB_ADDR`: Hostname and port of the RethinkDB server. Defaults to `localhost:28015`.
- `STATBAN_DB`: Name of the database. Defaults to `statban`.
- `GITHUB_TOKEN`: Github token to make API calls. This is required.
- `TARGET_REPOSITORY`: Repository to poll. This is required.
- `REPO_OWNER`: Owner of the repository. Required.
- `LABELS`: Issue labels. Defaults to `ready, development, review, release, done`.

## Dependencies

Statban uses Godep to manage dependencies. Dependencies should be vendored in as part of the repo.

You can use `go get` to fetch them for local development.

## Building and running

To build, use `go build`

This should build you an executable, which you can run using `./statban`

