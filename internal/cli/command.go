package cli

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("time between feed checks required")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("time between requests must be a duration")
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *State) error {
	dbFeed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	err = s.DB.MarkFeedFetched(context.Background(), dbFeed.ID)
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	rssFeed, err := rss.FetchFeed(context.Background(), dbFeed.Url)
	if err != nil {
		return fmt.Errorf("couldn't get rss feed: %w", err)
	}

	fmt.Printf("%d posts in feed %s:\n", len(rssFeed.Channel.Item), dbFeed.Name)
	for _, item := range rssFeed.Channel.Item {
		fmt.Println("* " + item.Title)

		titleValid := len(item.Title) > 0
		descriptionValid := len(item.Description) > 0

		dur, err := time.Parse("2006-01-02 15:04:05", item.PubDate)
		pubTime := sql.NullTime{}
		if err != nil {
			pubTime.Valid = false
		} else {
			pubTime.Valid = true
			pubTime.Time = dur
		}

		arg := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: titleValid},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: descriptionValid},
			PublishedAt: pubTime,
			FeedID:      dbFeed.ID,
		}

		_, err = s.DB.CreatePost(context.Background(), arg)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				continue
			}
			return fmt.Errorf("couldn't save post: %w", err)
		} else {
			fmt.Printf("Saved post %s\n", item.Link)
		}
	}
	fmt.Println()
	return nil
}

func HandlerAddFeed(s *State, cmd Command, u database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("name and url of the feed required")
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

	arg2 := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    f.ID,
	}

	_, err = s.DB.CreateFeedFollow(context.Background(), arg2)
	if err != nil {
		return fmt.Errorf("couldn't create feed following: %w", err)
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

func HandlerFollow(s *State, cmd Command, u database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("feed url required")
	}

	f, err := s.DB.GetFeedByUrl(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}

	arg := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    f.ID,
	}

	ff, err := s.DB.CreateFeedFollow(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("user couldn't follow feed: %w", err)
	}

	fmt.Printf("User %s now follows feed %s\n", ff.UserName, ff.FeedName)

	return nil
}

func HandlerFollowing(s *State, cmd Command, u database.User) error {
	ffu, err := s.DB.GetFeedFollowsForUser(context.Background(), u.Name)
	if err != nil {
		return fmt.Errorf("couldn't find user's feeds: %w", err)
	}

	if len(ffu) == 0 {
		fmt.Println("User does not follow any feeds")
		return nil
	}

	fmt.Printf("Feeds of %s:\n", u.Name)
	for _, ff := range ffu {
		fmt.Printf("* %s\n", ff.FeedName)

	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, u database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("feed url required")
	}

	arg := database.DeleteFeedFollowParams{
		UserID: u.ID,
		Url:    cmd.Arguments[0],
	}

	err := s.DB.DeleteFeedFollow(context.Background(), arg)

	if err != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	fmt.Println("Successfully unfollowed feed")
	return nil
}

func HandlerBrowse(s *State, cmd Command) error {
	limit := 2
	if len(cmd.Arguments) > 0 {
		l, err := strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return fmt.Errorf("give limit or no argument (defaults to 2): %w", err)
		} else {
			limit = l
		}
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), int32(limit))

	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	fmt.Printf("%d mot recent pots in database:\n", limit)
	for _, post := range posts {
		fmt.Printf("* %v\n", post.Title)
	}
	return nil
}
