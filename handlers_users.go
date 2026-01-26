package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/morauszkia/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.name)
	}
	name := cmd.args[0]

	if _, err := s.db.GetUser(context.Background(), name); err != nil {
		return fmt.Errorf("Couldn't find user: %w", err)
	}

	if err := s.config.SetUser(name); err != nil {
		return fmt.Errorf("Couldn't set current user: %w", err)
	}
	fmt.Printf("Login successful. Current user is %s\n", s.config.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.name)
	}
	name := cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		Name: 		name,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create user: %w", err)
	}

	if err := s.config.SetUser(user.Name); err != nil {
		return fmt.Errorf("Couldn't set current user: %w", err)
	}
	fmt.Println("User created successfully.")
	printUser(user)
	return nil
}

func printUser(user database.User) {
	fmt.Printf("User id:    %v\n", user.ID)
	fmt.Printf("Username:   %v\n", user.Name)
}

func handlerUsers(s * state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couln't get list of users: %w", err)
	}
	current := s.config.CurrentUserName

	for _, user := range users {
		if user.Name == current {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}