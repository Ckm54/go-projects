package main

import (
	"context"
	"fmt"

	"github.com/ckm54/go-projects/gator/internal/commands"
	"github.com/ckm54/go-projects/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *commands.State, cmd commands.Command, user database.User) error) func(*commands.State, commands.Command) error {
	return func(s *commands.State, c commands.Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("not logged in: %w", err)
		}

		return handler(s, c, user)
	}
}
