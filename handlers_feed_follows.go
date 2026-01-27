package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/morauszkia/gator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.name)
	}

	userName := s.config.CurrentUserName
	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("Couldn't get current user data: %w", err)
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Couldn't get feed data: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed follow entry created successfully.")
	printFeedFollow(feedFollow)
	
	return nil
}

func handlerFollowing(s* state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Couldn't get current user data: %w", err)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Couldn't get followed feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Printf("No feed follows found for %s\n", user.Name)
	}

	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.FeedName)
	}

	return nil
}
