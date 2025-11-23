package commands

import (
	"context"
	"fmt"
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
		fmt.Println("📰", item.Title)
	}

	return nil
}
