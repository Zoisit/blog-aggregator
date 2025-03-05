package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Zoisit/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type Command struct {
	Name      string
	Arguments []string
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("username required")
	}

	_, err := s.DB.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return fmt.Errorf("couldn't find user: user is not registered in database")
		}
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.Config.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("name required")
	}

	arg := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
	}

	u, err := s.DB.CreateUser(context.Background(), arg)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_name_key"` {
			return fmt.Errorf("couldn't create user: user already exists")
		}
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.Config.SetUser(u.Name)
	if err != nil {
		return err
	}
	fmt.Printf("Created user %s:\nID:%v\nCreated at:%v,\nUpate at:%v", u.Name, u.ID, u.CreatedAt, u.UpdatedAt)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.DB.DeleteAllUsers(context.Background())

	if err != nil {
		return fmt.Errorf("couldn't delete all user: %w", err)
	}

	fmt.Println("Successfully deleted all users")
	return nil
}
