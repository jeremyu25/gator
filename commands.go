package main

import (
	"fmt"

	"github.com/jeremyu25/blog_aggregator/internal/config"
	"github.com/jeremyu25/blog_aggregator/internal/database"
)

type state struct {
	cfg       *config.Config
	dbQueries *database.Queries
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, found := c.commandMap[cmd.Name]
	if !found {
		fmt.Print(cmd.Args)
		return fmt.Errorf("command could not be found\n")
	}
	return command(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func registerAllCommands(commands *commands) {
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerGetAllUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))
}
