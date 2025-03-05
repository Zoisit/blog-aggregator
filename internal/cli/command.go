package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Zoisit/blog-aggregator/internal/database"
	"github.com/Zoisit/blog-aggregator/internal/rss"
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

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())

	if err != nil {
		return fmt.Errorf("couldn't get list of users: %w", err)
	}

	fmt.Println("Users in database:")
	for _, u := range users {
		if s.Config.CurrentUserName == u.Name {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	rssFeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't get rss feed: %w", err)
	}

	fmt.Println(rssFeed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("name and url of the feed required")
	}

	u, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	arg := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    u.ID,
	}

	f, err := s.DB.CreateFeed(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	fmt.Printf("Created feed %s:\nID:%v\nCreated at:%v,\nUpate at:%v\nURL:%v\nuser_id:%v\n", f.Name, f.ID, f.CreatedAt, f.UpdatedAt, f.Url, f.UserID)

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.DB.GetFeedsWithUsername(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get list of feeds: %w", err)
	}

	fmt.Println("Feeds in database:")
	for _, f := range feeds {
		fmt.Printf("Name: %s\nURL: %s\nUser: %s\n\n", f.Name, f.Url, f.UserName)
	}
	return nil
}
