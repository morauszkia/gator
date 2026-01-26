package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}
	if err := s.db.Reset(context.Background()); err != nil {
		return fmt.Errorf("Error resetting database: %w", err)
	}
	fmt.Println("Database reset successfully")
	return nil
}