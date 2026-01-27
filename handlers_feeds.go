package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/morauszkia/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.name)
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Couldn't fetch current user data: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create feed entry: %w", err)
	}

	fmt.Println("Feed added successfully")
	printFeed(feed)
	
	return nil
}


func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}
	feedUrl := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Couldn't fetch feed: %w", err)
	}
	fmt.Printf("%+v", *feed)
	return nil
}