package cli

import (
	"context"
	"fmt"

	"github.com/Zoisit/blog-aggregator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't find user: %w", err)

		}

		return handler(s, cmd, user)
	}
}
