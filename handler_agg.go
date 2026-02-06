package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/morauszkia/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_requests>", cmd.name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Couldn't parse time between requests: %w", err)
	}
	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests.String())
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't get next feed to fetch: %w", err)
	}

	fmt.Printf("Fetching feed %s\n", nextFeed.Name)
	feedData, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Couldn't fetch feed: %w", err)
	}

	for _, item := range feedData.Channel.Item {
		pubDate, ok := parsePublishedAt(item.PubDate)
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: sql.NullString{
				String: item.Title,
				Valid: item.Title != "",
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid: item.Description != "",
			},
			PublishedAt: sql.NullTime{
				Time: pubDate,
				Valid: ok,
			},
			FeedID: nextFeed.ID,
		})
		if err != nil {
			return err
		}
	}

	if err := s.db.MarkFeedFetched(context.Background(), nextFeed.ID); err != nil {
		return fmt.Errorf("Error marking fetched feed: %w", err)
	}

	fmt.Printf("Collected feed %s. Found %d posts\n", nextFeed.Name, len(feedData.Channel.Item))
	return nil
}