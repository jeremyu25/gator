package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jeremyu25/blog_aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}
	username := cmd.Args[0]
	_, err := s.dbQueries.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User does not exist in database. Unable to login.\n")
	}
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("User with username %s has been set\n", username)
	return nil
}

func handlerGetAllUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	users, err := s.dbQueries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get all users")
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}
	username := cmd.Args[0]
	if username == "" {
		return fmt.Errorf("name cannot be empty! usage: %s <name>\n", cmd.Name)
	}
	user, err := s.dbQueries.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		fmt.Printf("Error executing query: %s\n", err)
		os.Exit(1)
	}
	if err = s.cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("User Created:\nID: %s\nCreatedAt: %s\nUpdatedAt: %s\nName: %s\n", user.ID, user.CreatedAt, user.UpdatedAt, user.Name)
	return nil
}
