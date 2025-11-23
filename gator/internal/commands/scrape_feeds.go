package commands

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"log"
	"time"

	"github.com/ckm54/go-projects/gator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *State) error {
	ctx := context.Background()

	feed, err := s.DB.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("no feeds to fetch: %w", err)
	}

	fmt.Printf("🪐 fetching %s (%s)\n", feed.Name, feed.Url)
	if err = s.DB.MarkFeedFetched(ctx, feed.ID); err != nil {
		return err
	}

	rss, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("failed fetching feed url: %w", err)
	}

	for _, item := range rss.Channel.Item {
		publishedAt := parsePublishedTime(item.PubDate)
		_, err := s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       html.UnescapeString(item.Title),
			Url:         item.Link,
			Description: sql.NullString{String: html.UnescapeString(item.Description)},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})

		if err != nil {
			if isUniqueValidation(err) {
				continue
			}
			return fmt.Errorf("failed saving post: %w", err)
		}
	}
	log.Printf("✅ fetched and stored posts from %s", feed.Name)

	return nil
}
