package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/morauszkia/gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title		string		`xml:"title"`
		Link		string		`xml:"link"`
		Description	string		`xml:"description"`
		Item		[]RSSItem	`xml:"item"`
	}	`xml:"channel"`
}

type RSSItem struct {
	Title		string	`xml:"title"`
	Link		string	`xml:"link"`
	Description	string	`xml:"description"`
	PubDate		string	`xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Couln't send request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read body: %w", err)
	}

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("Couldn't decode XML data: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return &feed, nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("Id: 			%s\n", feed.ID)
	fmt.Printf("Created at: 	%v\n", feed.CreatedAt)
	fmt.Printf("Updated at: 	%v\n", feed.UpdatedAt)
	fmt.Printf("Feed title: 	%s\n", feed.Name)
	fmt.Printf("Feed url: 		%s\n", feed.Url)
	fmt.Printf("Added by: 		%s\n", feed.UserID)
}