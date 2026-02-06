package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/morauszkia/gator/internal/database"
)

const defaultLimit = 2

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error
	if len(cmd.args) == 0 {
		limit = defaultLimit
		fmt.Printf("No limit specified. Using default limit of %d\n", defaultLimit)
		fmt.Printf("You can specify limit by running: %s [limit]\n\n", cmd.name)
	} else if len(cmd.args) == 1 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			fmt.Printf("Could not parse provided limit to integer: %v\n", cmd.args[0])
			fmt.Printf("Using default limit of %d\n\n", defaultLimit)
		}
	} else {
		return fmt.Errorf("Usage: %s [limit]", cmd.name)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Could not fetch latest posts: %w\n", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title.String)
		fmt.Printf("Feed: %s\n", post.FeedName)
		fmt.Printf("URL: %s\n", post.Url)
		var publishedAt string
		if post.PublishedAt.Valid {
			publishedAt = post.PublishedAt.Time.Format(time.RFC1123)
		} else {
			publishedAt = "unknown"
		}	
		fmt.Printf("Published: %s\n", publishedAt)
		if post.Description.Valid {
			fmt.Printf("Description: %s\n", truncateText(strip.StripTags(post.Description.String), 100))
		}
		fmt.Println()
	}


	return nil
}