package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/ckm54/go-projects/gator/internal/commands"
	"github.com/ckm54/go-projects/gator/internal/config"
	"github.com/ckm54/go-projects/gator/internal/database"
)

func main() {
	configuration, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", configuration.DBURL)
	if err != nil {
		log.Fatalf("error opening db connection: %v", err)
	}

	dbQueries := database.New(db)
	state := &commands.State{Config: &configuration, DB: dbQueries}
	cmds := &commands.Commands{}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerGetUsers)
	cmds.Register("agg", commands.HandlerAggregate)
	cmds.Register("addfeed", middlewareLoggedIn(commands.HandlerAddFeed))
	cmds.Register("feeds", commands.HandlerGetFeeds)
	cmds.Register("follow", middlewareLoggedIn(commands.HandlerFollowFeed))
	cmds.Register("following", middlewareLoggedIn(commands.HandlerFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(commands.HandlerUnfollowFeed))
	cmds.Register("browse", middlewareLoggedIn(commands.HandlerBrowse))

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: no command provided")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]
	cmd := commands.Command{Name: cmdName, Args: cmdArgs}

	// run the command
	if err := cmds.Run(state, cmd); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
