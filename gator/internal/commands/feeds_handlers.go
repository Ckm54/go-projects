package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ckm54/go-projects/gator/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandlerAggregate(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: gator agg <duration> (e.g. 30s, 5m, 1h)")
	}

	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid time duration: %w", err)
	}

	fmt.Printf("⏳ Collecting feeds every %s\n", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err = scrapeFeeds(s); err != nil {
			fmt.Println("⚠️", err)
		}
		<-ticker.C
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	switch len(cmd.Args) {
	case 0:
		return fmt.Errorf("missing name and url.\nusage: gator add <name> <url>")
	case 1:
		return fmt.Errorf("missing url.\nusage: gator add <name> <url>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feedInfo := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.DB.CreateFeed(context.Background(), feedInfo)
	if err != nil {
		return err
	}
	fmt.Println("✅ Feed added successfully")
	fmt.Printf("%v\n", feed)

	// auto follow added feed
	if err = HandlerFollowFeed(s, Command{Name: "follow", Args: []string{string(feed.Url)}}, user); err != nil {
		fmt.Println("❌ Could not follow feed")
		return err
	}
	fmt.Println("✅ Feed followed successfully")

	return nil
}

func HandlerGetFeeds(s *State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
	} else {
		for _, feed := range feeds {
			fmt.Println(feed)
		}
	}
	return nil
}

func HandlerFollowFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("missing url.\nusage: gator follow <url>")
	}

	feedUrl := cmd.Args[0]
	feed, err := getFeedByUrl(s, feedUrl)
	if err != nil {
		return err
	}

	feedFollowsData := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), feedFollowsData)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			fmt.Printf("ℹ️  %s is already following %s\n", user.Name, feed.Name)
			return nil
		}

		return fmt.Errorf("could not save follow: %w", err)
	}

	fmt.Printf("✅ %s now follows %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	userFeeds, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting user feeds: %w", err)
	}

	if len(userFeeds) == 0 {
		fmt.Println("You are not following any feeds yet")
	} else {
		fmt.Println("You are following:")
		for _, feed := range userFeeds {
			fmt.Printf("- %s\n", feed.FeedName)
		}
	}
	return nil
}

func HandlerUnfollowFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("missing url.\nusage: gator unfollow <url>")
	}

	feed, err := getFeedByUrl(s, cmd.Args[0])
	if err != nil {
		return err
	}

	params := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	if err = s.DB.UnfollowFeed(context.Background(), params); err != nil {
		return fmt.Errorf("could not unfollow feed: %w", err)
	}
	fmt.Printf("you have unfollowed %s\n", feed.Name)

	return nil
}

func getFeedByUrl(s *State, feedUrl string) (database.Feed, error) {
	feed, err := s.DB.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return database.Feed{}, fmt.Errorf("error finding feed: %w", err)
	}

	return feed, nil
}

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		parsed, err := strconv.Atoi(cmd.Args[0])
		if err != nil || parsed < 1 {
			return fmt.Errorf("invalid limit\nusage: gator browse [limit]")
		}

		limit = parsed
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		fmt.Println("no posts yet - try running gator agg 1m")
		return nil
	}

	for _, post := range posts {
		fmt.Printf("\n📌 %s\n🔗 %s\n📰 %s\n", post.FeedName, post.Url, post.Title)
	}

	return nil

}
