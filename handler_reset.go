package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if err := s.dbQueries.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("Failed to reset DB: %s", err)
	}
	fmt.Println("DB Reset Successful")
	return nil
}
