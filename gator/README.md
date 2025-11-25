# Gator — RSS Feed Aggregator CLI

Gator is a command-line tool for subscribing to RSS feeds, aggregating posts, and browsing them directly in your terminal.

## Requirements

Before installing, ensure you have:

- **Go 1.22+**
- **PostgreSQL 14+**
- A Unix-like terminal (macOS, Linux, WSL)

## Installation

Clone the repo:

```bash
git clone https://github.com/ckm54/go-projects/gator.git
cd gator
```

Install the CLI:

```bash
go install .
```

> This places the gator binary in your $GOPATH/bin—make sure it’s on your PATH.

## Database Setup

Create a Postgres database:

```bash
createdb gator
```

Run migrations:
Navigate into the `/sql/schema` folder and run:

```bash
goose postgres postgres://postgres:@localhost:5432/gator up
```

## Configuration

Gator stores configuration in a JSON file located in your home directory `~/.gatorconfig.json`:

```bash
Example:
{
"db_url": "postgres://localhost:5432/gator?sslmode=disable",
"current_user_name": "collins"
}
```

Make sure the DB URL matches your local setup.

## Usage

Initialize the database and register a user:

```bash
gator register <username>
```

Log in:

```bash
gator login <username>
```

Add a feed:

```bash
gator add https://example.com/feed.xml
```

Follow a feed:

```bash
gator follow https://example.com/feed.xml
```

Browse recent posts (with optional limit):

```bash
gator browse 5
```

Aggregate feeds manually:

```bash
gator agg 1m
```

(Unattended scraping is usually handled via cron or background jobs.)

## Development

Generate sqlc code:

```bash
sqlc generate
```

Project Structure

```bash
.
├── internal
│ ├── commands // CLI commands
│ └── config // Setup user
│ └── database // generated sqlc queries
├── sql
│ ├── queries // SQL queries
│ ├── schema // Migrations
├── go.mod
├── go.sum
├── main.go
├── middleware.go
├── sqlc.yaml
└── README.md
```
