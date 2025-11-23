package commands

import (
	"github.com/ckm54/go-projects/gator/internal/config"
	"github.com/ckm54/go-projects/gator/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	handlers map[string]func(*State, Command) error
}
