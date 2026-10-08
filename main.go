package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jeremyu25/blog_aggregator/internal/config"
	"github.com/jeremyu25/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfgJson := config.Read()
	dbURL := cfgJson.DbUrl
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening DB")
	}
	dbQueries := database.New(db)
	appState := state{
		cfg:       &cfgJson,
		dbQueries: dbQueries,
	}
	commands := commands{
		commandMap: map[string]func(*state, command) error{},
	}
	registerAllCommands(&commands)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments provided. Proper usage: <command> <args...>")
		os.Exit(1)
	}
	command := command{
		Name: args[1],
		Args: args[2:],
	}
	if err := commands.run(&appState, command); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
