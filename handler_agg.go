package main

import (
	"context"
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_requests>", cmd.name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Couldn't parse time between requests: %w", err)
	}
	fmt.Printf("Collecting feeds every %s", timeBetweenRequests.String())
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
	for i, item := range feedData.Channel.Item {
		fmt.Printf("%d: %s\n", i+1, item.Title)
	}
	if err := s.db.MarkFeedFetched(context.Background(), nextFeed.ID); err != nil {
		return fmt.Errorf("Error marking fetched feed: %w", err)
	}
	fmt.Printf("Collected feed %s. Found %d posts", nextFeed.Name, len(feedData.Channel.Item))
	return nil
}